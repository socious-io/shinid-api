package views

import (
	"context"
	"fmt"
	"net/http"
	"shin/src/app/auth"
	"shin/src/app/models"
	"shin/src/database"
	"shin/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func verificationsGroup(router *gin.Engine) {
	g := router.Group("verifications")

	g.GET("", paginate(), auth.LoginRequired(), func(c *gin.Context) {
		u, _ := c.Get("user")
		page, _ := c.Get("paginate")
		verifications, total, err := models.GetVerifications(u.(*models.User).ID, page.(database.Paginate))
		if err != nil {
			fmt.Println(err, "----------")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"results": verifications,
			"total":   total,
		})
	})

	g.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		v, err := models.GetVerification(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, v)
	})

	g.POST("", auth.LoginRequired(), func(c *gin.Context) {
		form := new(VerificationForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		v := new(models.Verification)
		utils.Copy(form, v)
		u, _ := c.Get("user")
		v.UserID = u.(*models.User).ID
		ctx, _ := c.Get("ctx")
		if err := v.Create(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, v)
	})

	g.PUT("/:id", auth.LoginRequired(), func(c *gin.Context) {
		form := new(VerificationForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id := c.Param("id")
		v, err := models.GetVerification(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u, _ := c.Get("user")
		if v.UserID.String() != u.(*models.User).ID.String() {
			c.JSON(http.StatusForbidden, gin.H{"error": "not allow"})
			return
		}

		if v.VerifiedAt != nil {
			form.SchemaID = v.SchemaID
		}
		utils.Copy(form, v)

		ctx, _ := c.Get("ctx")
		if err := v.Update(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, v)
	})

	g.DELETE("/:id", auth.LoginRequired(), func(c *gin.Context) {
		id := c.Param("id")
		v, err := models.GetVerification(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u, _ := c.Get("user")
		if v.UserID.String() != u.(*models.User).ID.String() {
			c.JSON(http.StatusForbidden, gin.H{"error": "not allow"})
			return
		}
		ctx, _ := c.Get("ctx")
		if err := v.Delete(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
}
