package views

import "github.com/gin-gonic/gin"

func userGroup(router *gin.Engine) {
	g := router.Group("users")

	g.POST("/", func(ctx *gin.Context) {})
}
