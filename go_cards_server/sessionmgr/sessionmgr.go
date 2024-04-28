package sessionmgr

import (
	"fmt"
	"log"
	"time"

	"example.com/go_cards_server/cheat"
	"example.com/go_cards_server/gametypes"
	"example.com/go_cards_server/user"
	"example.com/go_cards_server/session"
	"github.com/google/uuid"
)

var Sessions map[uuid.UUID]*session.Session

func InitSessionEngine() {
	Sessions = make(map[uuid.UUID]*session.Session)

	// Removes sessions that have not had any messages after a predetermined time
	go sessionCleaner()
}

func GetSession(sessionId uuid.UUID) (*session.Session, error) {
	if session, ok := Sessions[sessionId]; ok {
		return session, nil
	}
	return nil, fmt.Errorf("session not found")
}

func CreateSession(gameType gametypes.GameType) *session.Session {
	session := session.CreateSession()
	Sessions[session.SessionId] = session

	switch gameType {
	case gametypes.Cheat:
		var cheatGame = cheat.CreateNewCheatSession(session)
		log.Println("Created new War session ", cheatGame) // TODO store this somewhere so we can track it
	default:
		log.Println("Unknown game type")
	}

	return session
}

func HandleUserJoined(s *session.Session, user *user.User) {
	s.AddUserToSession(user)
}

func sessionCleaner() {
	ticker := time.NewTicker(20 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		for sessionId, session := range Sessions {
			if time.Since(session.SessionLastMessageTime) > 60*time.Second {
				session.EndSession()
				session.Active = false

				// Ensure each user connection is closed
				for _, user := range session.Users {
					if user.UserConnection != nil {
						user.UserConnection.WriteMessage(1, []byte("Session ending due to inactivity"))
						user.UserConnection.Close()
					}
				}

				fmt.Println("Cleaning session:", sessionId)
				delete(Sessions, sessionId)
			}
		}
	}
}
