package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Login(c *gin.Context) {
	userId := uuid.New().String()
	// Get the value of the "username" query parameter
	username := c.Query("username")

	c.SetCookie("userId", userId, 21600 /* age */, "/" /* valid for all paths */, os.Getenv("REACT_APP_DOMAIN"), false /* secure */, false /* HTTP only */)
	c.SetCookie("username", username, 21600 /* age */, "/" /* valid for all paths */, os.Getenv("REACT_APP_DOMAIN"), false /* secure */, false /* HTTP only */)

	fmt.Println("Set cookie!")
	c.JSON(http.StatusOK, gin.H{
		"userId":   userId,
		"username": username,
	})
}
