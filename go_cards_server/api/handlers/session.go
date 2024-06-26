package handlers

import (
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

func GetTicket(c *gin.Context) {
	// Get the authToken from the POST request body
	var authRequest AuthRequest

	// Bind the JSON payload to the struct
	if err := c.ShouldBindJSON(&authRequest); err != nil {
		// If there is an error in binding, return a bad request response
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := jwthelper.VerifyJWT(authRequest.AuthToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Error verifying token"})
	}

	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
	}

	ticket := session.GetTicket()
	session.ClaimsForTickets[ticket] = claims
	c.JSON(http.StatusOK, gin.H{"ticket": ticket})
}

func ConnectSession(c *gin.Context) {
	var s *session.Session = nil
	ticket := c.Query("ticket")
	claims := session.ClaimsForTickets[ticket]

	userIdUuid, err := uuid.Parse(claims.UserID)
	username := claims.Username

	session.DeleteClaimsForTicket(ticket)

	if err != nil {
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
