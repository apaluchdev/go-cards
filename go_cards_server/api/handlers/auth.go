package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Login(c *gin.Context) {
	userId := uuid.New().String()
	// Get the value of the "username" query parameter
	username := c.Query("username")

	setCookie(c, "userId", userId, 21600 /* age */, "/" /* valid for all paths */, os.Getenv("REACT_APP_DOMAIN"), !strings.Contains(os.Getenv("REACT_APP_DOMAIN"), "localhost"), false /* HTTP only */)
	setCookie(c, "username", username, 21600 /* age */, "/" /* valid for all paths */, os.Getenv("REACT_APP_DOMAIN"), !strings.Contains(os.Getenv("REACT_APP_DOMAIN"), "localhost"), false /* HTTP only */)

	fmt.Println("Set cookie!")
	c.JSON(http.StatusOK, gin.H{
		"userId":   userId,
		"username": username,
	})
}

func setCookie(c *gin.Context, name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		Domain:   domain,
		MaxAge:   maxAge,
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: http.SameSiteNoneMode, // explicitly set SameSite=None
	})
}
