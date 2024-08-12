package handlers

import (
	"log"
	"net/http"

	"example.com/go_cards_server/gametypes"
	"example.com/go_cards_server/jwthelper"
	"example.com/go_cards_server/session"
	"example.com/go_cards_server/sessionmgr"
	"example.com/go_cards_server/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type AuthRequest struct {
	AuthToken string `json:"authToken" binding:"required"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections | TODO Make specific for web app
	},
}

func ConnectSession(c *gin.Context) {
	var s *session.Session = nil

	tokenString, err := c.Cookie("Authorization")
	if tokenString == "" || err != nil {
		log.Println("No token provided")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		c.Abort()
		return
	}

	claims := &jwthelper.Claims{}
	claims, err = jwthelper.VerifyJWT(tokenString)

	if err != nil {
		log.Println("Invalid token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	userIdUuid, err := uuid.Parse(claims.UserID)
	username := claims.Username

	if err != nil {
		log.Println("Error occurred with parsing user id")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Error occurred with parsing user id"})
	}

	// Check if the session id is provided in the query parameters
	sessionId, err := uuid.Parse(c.Query("id"))
	if err != nil {
		sessionId = uuid.Nil
	}

	// Get or create a session
	if sessionId == uuid.Nil {
		s = sessionmgr.CreateSession(gametypes.Cheat)
	} else {
		s, err = sessionmgr.GetSession(sessionId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
			return
		}
	}

	if (len(s.Users) >= s.MaxUsers) && s.MaxUsers != 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Session is full"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.Error(c.Writer, "Failed to upgrade to websocket connection", http.StatusInternalServerError)
		return
	}

	sessionmgr.HandleUserJoined(s, &user.User{UserId: userIdUuid, UserName: username, UserConnection: conn})
}

func GetSession(c *gin.Context) {
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusOK, sessionmgr.Sessions)
		return
	}

	if _, exists := sessionmgr.Sessions[uuid]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	session := sessionmgr.Sessions[uuid]
	c.JSON(http.StatusOK, session)
}
