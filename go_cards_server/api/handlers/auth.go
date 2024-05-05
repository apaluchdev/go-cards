package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Test for commit - remove me later
func Login(c *gin.Context) {
	userId := uuid.New().String()
	// Get the value of the "username" query parameter
	username := c.Query("username")

	c.SetCookie("userId", userId, 21600 /* age */, "/" /* valid for all paths */, "localhost", false /* secure */, false /* HTTP only */)
	c.SetCookie("username", username, 21600 /* age */, "/" /* valid for all paths */, "localhost", false /* secure */, false /* HTTP only */)

	fmt.Println("Set cookie!")
	c.JSON(http.StatusOK, gin.H{
		"userId":   userId,
		"username": username,
	})
}
