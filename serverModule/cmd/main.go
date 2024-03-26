package main

import (
	"net/http"

	"example.com/server/routes/auth"
	"example.com/server/routes/session"
	"example.com/server/session_manager"
	"github.com/gin-gonic/gin"
)

func main() {
	// file, err := os.Create("logfile.txt")
	// if err != nil {
	// 	log.Fatal("Cannot create log file:", err)
	// }
	// defer file.Close()

	// log.SetOutput(file)

	r := gin.Default()

	reactAppDomain := "http://localhost:3000" //os.Getenv("REACT_APP_DOMAIN")

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

	session_manager.InitSessionEngine()

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", auth.Login)
	}

	sessionGroup := r.Group("/session")
	{
		// Middleware to ensure the client has a userId cookie set
		sessionGroup.Use(func(c *gin.Context) {
			userIdCookie, err := c.Cookie("userId")
			if err != nil || userIdCookie == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				c.Abort()
				return
			}
			c.Next()
		})

		sessionGroup.GET("/", session.GetSession)
		sessionGroup.GET("/connect", session.ConnectSession)
	}

	r.Run()
}
