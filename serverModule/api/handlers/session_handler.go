package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"example.com/server/models"
	"example.com/server/session_manager"
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

func GetSession(c *gin.Context) {
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusOK, session_manager.Sessions)
		return
	}

	if _, exists := session_manager.Sessions[uuid]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	session := session_manager.Sessions[uuid]
	c.JSON(http.StatusOK, session)
}

func CreateSession(c *gin.Context) {
	uuid := uuid.New()
	s := &models.Session{SessionId: uuid, SessionStartTime: time.Now()}

	session_manager.Sessions[uuid] = s
	// Send a response
	c.JSON(http.StatusOK, s)
}

func ConnectSession(c *gin.Context) {
	log.Println("User connecting to session...")

	var session *models.Session = nil

	val, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not authenticated"})
		return
	}
	userIdUuid := val.(uuid.UUID)

	val, exists = c.Get("username")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not authenticated"})
		return
	}
	username := val.(string)

	// Check if the session id is provided in the query parameters
	id := c.Query("id")
	sessionId, err := uuid.Parse(id)
	if err != nil {
		sessionId = session_manager.CreateSession()
	}

	// Verify that the session exists or was created successfully
	session, exists = session_manager.Sessions[sessionId]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Session '%s' does not exist", sessionId)})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.Error(c.Writer, "Failed to upgrade to websocket connection", http.StatusInternalServerError)
		return
	}

	// Add the player to the session
	session.Players[userIdUuid] = &models.Player{PlayerId: userIdUuid, PlayerName: username, PlayerConnection: conn}

	// Send the player the session details and notify other players of the new player
	session_manager.SendMessage(conn, models.CreateMessage(session, models.SessionStartedMessageType))
	session_manager.BroadcastMessage(session.SessionId, models.CreateMessage(session.Players[userIdUuid], models.PlayerJoinedMessageType))

	// Handle the websocket communication with this player
	go session_manager.HandleUserSession(conn, session, userIdUuid)
}
