package handlers

import (
	"net/http"

	"example.com/go_cards_server/api/cookieutil"
	"example.com/go_cards_server/gametypes"
	"example.com/go_cards_server/player"
	"example.com/go_cards_server/session"
	"example.com/go_cards_server/sessionmgr"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections | TODO Make specific for web app
	},
}

func ConnectSession(c *gin.Context) {
	var s *session.Session = nil

	userIdUuid, username, err := cookieutil.GetUserIdAndUsernameFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.Error(c.Writer, "Failed to upgrade to websocket connection", http.StatusInternalServerError)
		return
	}

	// Check if the session id is provided in the query parameters
	sessionId, err := uuid.Parse(c.Query("id"))
	if err != nil {
		sessionId = uuid.Nil
	}

	// Get or create a session
	if sessionId == uuid.Nil {
		s = sessionmgr.CreateSession(gametypes.War)
	} else {
		s, err = sessionmgr.GetSession(sessionId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
			return
		}
	}

	if (len(s.Players) >= s.MaxPlayers) && s.MaxPlayers != 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Session is full"})
		return
	}

	sessionmgr.HandlePlayerJoined(s, &player.Player{PlayerId: userIdUuid, PlayerName: username, PlayerConnection: conn})
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
