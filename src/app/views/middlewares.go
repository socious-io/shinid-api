package views

import (
	"shin/src/database"
	"strconv"

	"github.com/gin-gonic/gin"
)

func paginate() gin.HandlerFunc {
	return func(c *gin.Context) {

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 1
		}

		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			limit = 10
		}
		if page < 1 {
			page = 1
		}
		if limit > 100 || limit < 1 {
			limit = 10
		}

		c.Set("paginate", database.Paginate{
			Limit: limit,
			Offet: (page - 1) * limit,
		})
		c.Set("limit", limit)
		c.Set("page", page)
		c.Next()

	}
}
