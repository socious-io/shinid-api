package views

import "github.com/gin-gonic/gin"

func orgGroup(router *gin.Engine) {
	g := router.Group("orgs")

	g.POST("/", func(ctx *gin.Context) {})
}
