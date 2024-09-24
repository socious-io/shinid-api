package views

import (
	"context"
	"fmt"
	"net/http"
	"shin/src/app/auth"
	"shin/src/app/models"
	"shin/src/config"
	"shin/src/database"
	"shin/src/lib"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func kybVerificationGroup(router *gin.Engine) {
	g := router.Group("kyb")
	g.Use(auth.LoginRequired())

	g.POST("/:org_id", func(c *gin.Context) {
		orgID := c.Param("org_id")
		u, _ := c.Get("user")
		ctx, _ := c.Get("ctx")

		form := new(KYBVerificationForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		org, err := models.GetOrgByMember(uuid.MustParse(orgID), u.(*models.User).ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   err.Error(),
				"message": "Organization not found",
			})
			return
		}

		kyb := &models.KYBVerification{
			UserID: u.(*models.User).ID,
			OrgID:  org.ID,
		}

		kyb, err = kyb.Create(ctx.(context.Context), form.Documents)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// lib.DiscordSendWithComponents(
		// 	config.Config.Logger.Discord["shin_kyb_channel"],
		// 	fmt.Sprintf(`
		// 		ID: %s
		// 	`, kyb.ID),
		// 	lib.DiscordComponent{
		// 		Type: 1,
		// 		Components: []lib.DiscordButton{
		// 			{
		// 				Type:  lib.DiscordButtonSimple,
		// 				Label: "Approve",
		// 				Style: lib.DiscordStyleSuccess,
		// 				URL:   "URL to Approve",
		// 			},
		// 			{
		// 				Type:  lib.DiscordButtonSimple,
		// 				Label: "Reject",
		// 				Style: lib.DiscordStyleDanger,
		// 				URL:   "URL to Reject",
		// 			},
		// 		},
		// 	},
		// )

		documents := ""
		message := fmt.Sprintf("ID: %s\n", kyb.ID)

		for i, document := range kyb.Documents {
			documents = fmt.Sprintf("%s\n%v. %s/%s", documents, i, config.Config.S3.CDNUrl, document.Url)
		}
		message += fmt.Sprintf("Documents:%s\n\n", documents)
		message += fmt.Sprintf("Approve: https://api.shinid.com/kyb/%s/approve\n", kyb.ID)
		message += fmt.Sprintf("Reject: https://api.shinid.com/kyb/%s/reject\n", kyb.ID)

		lib.DiscordSendTextMessage(
			config.Config.Logger.Discord["shin_kyb_channel"],
			message,
		)

		c.JSON(http.StatusOK, kyb)
	})

	g.GET("/", paginate(), func(c *gin.Context) {
		u, _ := c.Get("user")
		paginate, _ := c.Get("paginate")
		limit, _ := c.Get("limit")
		page, _ := c.Get("page")

		kybVerifications, total, err := models.GetAllByUserId(u.(*models.User).ID, paginate.(database.Paginate))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"results": kybVerifications,
			"page":    page,
			"limit":   limit,
			"total":   total,
		})
	})

	g.GET("/:id", func(c *gin.Context) {
		u, _ := c.Get("user")
		verificationId := c.Param("id")

		verification, err := models.GetById(uuid.MustParse(verificationId), u.(*models.User).ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, verification)
	})

	g.GET("/:id/approve", adminAccessRequired(), func(c *gin.Context) {

		ctx, _ := c.Get("ctx")
		verificationId := c.Param("id")

		kybVerification := models.KYBVerification{
			ID: uuid.MustParse(verificationId),
		}

		if err := kybVerification.ChangeStatus(ctx.(context.Context), models.KYBStatusApproved); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	g.GET("/:id/reject", adminAccessRequired(), func(c *gin.Context) {
		ctx, _ := c.Get("ctx")
		verificationId := c.Param("id")

		kybVerification := models.KYBVerification{
			ID: uuid.MustParse(verificationId),
		}

		if err := kybVerification.ChangeStatus(ctx.(context.Context), models.KYBStatusRejected); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

}
