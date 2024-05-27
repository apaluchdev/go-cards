package middleware

import (
	"github.com/gin-gonic/gin"
)

func GuestAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// TODO - Lock routes except the authorization route

		// Continue processing
		c.Next()
	}
}
