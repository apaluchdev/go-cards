package router

import (
	"example.com/go_cards_server/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupHomeRoutes(r *gin.Engine) {
	authGroup := r.Group("/")
	{
		authGroup.GET("/login", handlers.Login)
	}
}
