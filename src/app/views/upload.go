package views

import (
	"context"
	"net/http"
	"shin/src/app/auth"
	"shin/src/app/models"
	"shin/src/services"

	"github.com/gin-gonic/gin"
)

func uploadGroup(router *gin.Engine) {
	g := router.Group("upload")
	g.PUT("/", auth.LoginRequired(), func(c *gin.Context) {

		file, header, _ := c.Request.FormFile("file")
		mediaUrl, fileName := services.Upload(file, header.Filename)

		u, _ := c.Get("user")

		media := models.Media{
			UserID:   u.(*models.User).ID,
			URL:      mediaUrl,
			Filename: fileName,
		}

		ctx, _ := c.Get("ctx")
		if err := media.Create(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "Couldn't upload media",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"result": media})

		file.Close()
	})
}
