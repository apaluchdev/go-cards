package cookieutil

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserIdAndUsernameFromContext(c *gin.Context) (uuid.UUID, string, error) {
	val, exists := c.Get("userId")
	if !exists {
		return uuid.Nil, "", errors.New("user not authenticated")
	}
	userIdUuid := val.(uuid.UUID)

	val, exists = c.Get("username")
	if !exists {
		return uuid.Nil, "", errors.New("user not authenticated")
	}
	username := val.(string)

	return userIdUuid, username, nil
}
