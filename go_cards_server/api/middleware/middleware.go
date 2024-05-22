package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GuestAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Middleware hit! v1.0.0.0")

		userId, err := c.Cookie("userId")
		if err != nil {
			// userId cookie not set, return an error
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "userId cookie not set.",
			})
			c.Abort() // Abort further processing
			return
		}

		username, err := c.Cookie("username")
		if err != nil {
			// userId cookie not set, return an error
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "username cookie not set",
			})
			c.Abort() // Abort further processing
			return
		}

		userIdUuid, err := uuid.Parse(userId)
		if err != nil {
			// userId cookie is not a valid uuid, return an error
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "userId cookie is not a valid uuid",
			})
			c.Abort() // Abort further processing
			return
		}

		// Set userId in context for downstream handlers
		c.Set("userId", userIdUuid)
		c.Set("username", username)

		// Continue processing
		c.Next()
	}
}
