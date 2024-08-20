package views

import (
	"context"
	"net/http"
	"shin/src/app/auth"
	"shin/src/app/models"
	"shin/src/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func userGroup(router *gin.Engine) {
	g := router.Group("users")
	g.Use(auth.LoginRequired())

	g.GET("/", func(c *gin.Context) {
		u, _ := c.Get("user")
		c.JSON(http.StatusOK, u)
	})

	g.POST("/profile/update", func(c *gin.Context) {
		form := new(ProfileUpdateForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		u, _ := c.Get("user")
		ctx, _ := c.Get("ctx")

		newProfile := new(models.User)
		utils.Copy(u, newProfile)

		//TODO: Scan the struct and iterate through it to Replace the form.* in newProfile.*

		newProfile.ID = u.(*models.User).ID
		newProfile.Username = strings.ToLower(newProfile.Username)

		if err := newProfile.UpdateProfile(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, "")
	})
}
