package router

import (
	"log"
	"os"

	router "example.com/go_cards_server/api/routes"
	"github.com/gin-gonic/gin"
)

func InitializeRouter() {
	log.Println("v1.0.0")

	r := gin.Default()
	reactAppDomain := os.Getenv("REACT_APP_DOMAIN")

	if reactAppDomain == "" {
		reactAppDomain = "http://localhost:3000"
	}

	// Allow requests from front end domain
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", reactAppDomain)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Cookie, Set-Cookie, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Handle preflight OPTIONS request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.SetupSessionRoutes(r)
	router.SetupAuthRoutes(r)

	r.Run()
}
