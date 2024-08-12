package handlers

import (
	"net/http"

	"example.com/go_cards_server/jwthelper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Login(c *gin.Context) {
	userId := uuid.New().String()
	username := c.Query("username")

	tokenString, err := jwthelper.GenerateJWT(userId, username)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
