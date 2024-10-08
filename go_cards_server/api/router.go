package router

import (
	"log"
	"net/http"
	"os"

	router "example.com/go_cards_server/api/routes"
	"github.com/gin-gonic/gin"
)

func InitializeRouter() {
	r := gin.Default()

	if os.Getenv("REACT_APP_DOMAIN") == "" {
		os.Setenv("REACT_APP_DOMAIN", "http://localhost:3000")
	}

	// Allow requests from front end domain
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("REACT_APP_DOMAIN"))
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
	r.GET("/", GetVersion)

	if os.Getenv("ENV") == "prod" {
		err := r.RunTLS(":8080", "/certs/fullchain.pem", "/certs/privkey.pem")
		if err != nil {
			log.Fatal("Failed to start server: ", err)
		}
	} else {
		err := r.Run()
		if err != nil {
			log.Fatal("Failed to start server: ", err)
		}
	}
}

func GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"apiVersion":   "1.0.0.7",
		"clientDomain": os.Getenv("REACT_APP_DOMAIN"),
	})
}
