package views

import (
	"context"
	"net/http"
	"shin/src/app/auth"
	"shin/src/app/models"
	"shin/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func orgGroup(router *gin.Engine) {
	g := router.Group("orgs")
	g.Use(auth.LoginRequired())

	g.GET("/", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		ctx, _ := c.Get("ctx")

		orgs, err := models.GetOrgsByMember(ctx.(context.Context), uuid.MustParse(userID.(string)))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, orgs)
	})

	g.GET("/:id", func(c *gin.Context) {
		orgID := c.Param("id")
		org, err := models.GetOrg(uuid.MustParse(orgID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, org)
	})

	g.POST("/", func(c *gin.Context) {
		form := new(OrganizationForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		o := new(models.Organization)
		utils.Copy(form, o)
		userID, _ := c.Get("user_id")
		ctx, _ := c.Get("ctx")
		if err := o.Create(ctx.(context.Context), uuid.MustParse(userID.(string))); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, o)
	})

	g.PUT("/:id", func(c *gin.Context) {
		orgID := c.Param("id")
		userID, _ := c.Get("user_id")
		// TODO: can be middleware
		o, err := models.GetOrgByMember(uuid.MustParse(orgID), uuid.MustParse(userID.(string)))
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		form := new(OrganizationForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		utils.Copy(form, o)
		ctx, _ := c.Get("ctx")
		if err := o.Update(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, o)
	})
}
