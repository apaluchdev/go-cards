package session

import (
	"fmt"
	"log"
	"net/http"
	"time"

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

func CheckSessionExists(c *gin.Context) {
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	_, exists := session_manager.Sessions[uuid]
	c.JSON(http.StatusOK, gin.H{"exists": exists})
}

func GetSession(c *gin.Context) {
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusOK, session_manager.Sessions)
		return
	}

	if _, exists := session_manager.Sessions[uuid]; !exists {
		// Session does not exist
		// Handle the case when the session does not exist
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	session := session_manager.Sessions[uuid]
	c.JSON(http.StatusOK, session)
}

func ConnectSession(c *gin.Context) {
	log.Println("User connecting to session...")

	var session *session_manager.Session = nil
	var userId uuid.UUID

	log.Println("Checking for valid userId cookie...")
	// Ensure the client has a userId cookie set
	if cookieVal, err := c.Cookie("userId"); err != nil {
		log.Println("Cookie not found")
		if err == http.ErrNoCookie {
			c.JSON(http.StatusNotFound, gin.H{"error": "No userId cookie set, cannot upgrade connection"})
			return
		}
	} else {
		uuid, err := uuid.Parse(cookieVal)
		if err != nil {
			log.Println("Could not parse cookie into a valid uuid")
			c.JSON(http.StatusNotFound, gin.H{"error": "Could not parse cookie into a valid uuid"})
			return
		} else {
			log.Println("userId Cookie found and parsed successfully")
			userId = uuid
		}
	}

	// Ensure the client has a username cookie set
	username, err := c.Cookie("username")
	if err != nil {
		log.Println("No username cookie set, cannot upgrade connection")
		c.JSON(http.StatusNotFound, gin.H{"error": "No username cookie set, cannot upgrade connection"})
		return
	}

	log.Println("Cookies set properly")
	// Verify a session id was passed in
	id := c.Query("id")
	fmt.Println("Id is: " + id)
	sessionId, err := uuid.Parse(id)
	if err != nil {
		log.Println("No session id; creating new session...")
		sessionId = uuid.New()
		session = &session_manager.Session{SessionId: sessionId, SessionStartTime: time.Now()}
		session.PlayerScores = make(map[uuid.UUID]int16)
		session_manager.Sessions[session.SessionId] = session
		session.Players = make(map[uuid.UUID]*session_manager.Player)

		log.Println("Adding AI player to session...")
		aiUuid := uuid.New()
		session.Players[aiUuid] = &session_manager.Player{PlayerId: aiUuid, PlayerName: "AI"}
	}

	// Loop through session_manager.Sessions and output all the keys
	log.Println("Sessions:")
	for key := range session_manager.Sessions {
		log.Println(key)
	}

	// Verify that the session was successfully created
	log.Printf("Attempting to connect to session %v\n", sessionId)
	session, exists := session_manager.Sessions[sessionId]
	if !exists {
		log.Println("Session does not exist for some reason...")
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Session '%s' does not exist", sessionId)})
		return
	}

	log.Println("Adding player to session...")
	session.Players[userId] = &session_manager.Player{PlayerId: userId, PlayerName: username}

	// Upgrade HTTP connection to WebSocket
	log.Println("Upgrading connection to WebSocket...")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.Error(c.Writer, "Failed to upgrade to websocket connection", http.StatusInternalServerError)
		return
	} else {
		session.Players[userId].PlayerConnection = conn
		log.Println("Sending initial session messages...")
		session_manager.SendInitialSessionMessage(conn, session)
		log.Println("Sending broadcast message...")
		session_manager.BroadcastMessage(session.SessionId, session_manager.CreatePlayerJoinedMessage(session.Players[userId]))
	}

	// Handle WebSocket communication here
	go session_manager.HandleUserSession(conn, session, userId)
}

func CreateSession(c *gin.Context) {
	uuid := uuid.New()
	s := &session_manager.Session{SessionId: uuid, SessionStartTime: time.Now()}

	session_manager.Sessions[uuid] = s
	// Send a response
	c.JSON(http.StatusOK, s)
}
