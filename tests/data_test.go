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

	schemasData = []gin.H{
		{"name": "test", "description": "test desc", "public": false, "attributes": []gin.H{
			{"name": "test_field", "description": "desc of field", "type": "TEXT"},
			{"name": "test_field_2", "description": "desc of field2", "type": "TEXT"},
		}},
		{"name": "test2", "description": "test2 desc", "public": false, "attributes": []gin.H{
			{"name": "test_field", "description": "desc of field", "type": "TEXT"},
			{"name": "test_field_2", "description": "desc of field2", "type": "TEXT"},
		}},
	}
)
