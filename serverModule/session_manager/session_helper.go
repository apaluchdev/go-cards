package session_manager

import (
	"fmt"
	"time"

	"example.com/server/models"
	"github.com/google/uuid"
)

func CreateSession() uuid.UUID {
	// Create a new session id
	sessionId := uuid.New()

	// Create a new session
	session := &models.Session{SessionId: sessionId, SessionStartTime: time.Now()}
	session.Players = make(map[uuid.UUID]*models.Player)

	Sessions[session.SessionId] = session

	return sessionId
}

func sessionCleaner() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		for sessionId, session := range Sessions {
			if time.Since(session.SessionLastMessageTime) > 60*time.Second {
				fmt.Println("Cleaning session:", sessionId)
				delete(Sessions, sessionId)
			}
		}
	}
}
