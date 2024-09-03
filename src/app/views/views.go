package views

import "github.com/gin-gonic/gin"

func Init(r *gin.Engine) {
	authGroup(r)
	orgGroup(r)
	userGroup(r)
	credntialsGroup(r)
	rootGroup(r)
	verificationsGroup(r)
	credentialsGroup(r)
	recipientsGroup(r)
	uploadGroup(r)
}
