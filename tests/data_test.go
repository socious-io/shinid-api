package tests_test

import "github.com/gin-gonic/gin"

var (
	authTokens        = []string{}
	authRefreshTokens = []string{}

	usersData = []gin.H{
		{
			"first_name": "TestName",
			"last_name":  "TestLastName",
			"username":   "test",
			"email":      "test@test.com",
			"password":   "test123456",
		},
	}

	organizationsData = []gin.H{
		{"name": "test", "description": "test desc"},
	}
)
