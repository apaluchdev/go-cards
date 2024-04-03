package router

import (
	"example.com/go_cards_server/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.Engine) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", handlers.Login)
	}
}
