package app

import (
	"context"
	"fmt"
	"shin/src/app/views"
	"shin/src/config"
	"time"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		c.Set("ctx", ctx)
		c.Next()
	})

	views.Init(router)
	return router
}

func Serve() {
	router := Init()
	router.Run(fmt.Sprintf("0.0.0.0:%d", config.Config.Port))
}
