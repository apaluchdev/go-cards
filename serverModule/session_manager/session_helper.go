package session_manager

import (
	"fmt"
	"log"
	"time"

	"example.com/server/models"
	war "example.com/server/war_manager"
	"github.com/google/uuid"
)

func CreateSession(gameType models.GameType) uuid.UUID {
	// Create a new session id
	sessionId := uuid.New()

	// Create a new session
	session := &models.Session{SessionId: sessionId, SessionStartTime: time.Now()}
	session.Players = make(map[uuid.UUID]*models.Player)
	session.GameType = gameType

	Sessions[session.SessionId] = session
	session.Deck = war.GetWarDeck()

	return sessionId
}

func AddPlayerToSession(sessionId uuid.UUID, player *models.Player) {
	session, exists := Sessions[sessionId]
	if !exists {
		log.Println("Error adding player to session: session does not exist")
		return
	}

	session.Players[player.PlayerId] = player
}

func sessionCleaner() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		for sessionId, session := range Sessions {

			// Ensure each player connection is closed
			for _, player := range session.Players {
				if player.PlayerConnection != nil {
					player.PlayerConnection.WriteMessage(1, []byte("Session ending due to inactivity"))
					player.PlayerConnection.Close()
				}
			}

			if time.Since(session.SessionLastMessageTime) > 60*time.Second {
				fmt.Println("Cleaning session:", sessionId)
				delete(Sessions, sessionId)
			}
		}
	}
}
