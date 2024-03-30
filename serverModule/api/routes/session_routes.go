package router

import (
	"example.com/server/api/handlers"
	"example.com/server/api/middleware"
	"github.com/gin-gonic/gin"
)

func SetupSessionRoutes(r *gin.Engine) {
	sessionGroup := r.Group("/session")
	{
		sessionGroup.Use(middleware.GuestAuthMiddleware())

		sessionGroup.GET("/", handlers.GetSession)
		sessionGroup.GET("/connect", handlers.ConnectSession)
	}
}
