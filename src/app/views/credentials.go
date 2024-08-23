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

func credntialsGroup(router *gin.Engine) {
	g := router.Group("credentials")
	g.Use(auth.LoginRequired())

	g.GET("/schemas", func(c *gin.Context) {

	})

	g.GET("/schemas/:id", func(c *gin.Context) {
		id := c.Param("id")
		// u, _ := c.Get("user")
		s, err := models.GetSchema(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, s)
	})

	g.POST("/schemas", func(c *gin.Context) {
		form := new(SchemaForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		s := new(models.Schema)
		utils.Copy(form, s)
		u, _ := c.Get("user")
		s.CreatedID = &u.(*models.User).ID
		ctx, _ := c.Get("ctx")
		if err := s.Create(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, s)
	})

	g.DELETE("/schemas/:id", func(c *gin.Context) {

	})
}
