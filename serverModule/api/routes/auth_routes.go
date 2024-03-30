package router

import (
	"example.com/server/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.Engine) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", handlers.Login)
	}
}
