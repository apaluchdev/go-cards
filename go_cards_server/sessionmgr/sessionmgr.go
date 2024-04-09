package sessionmgr

import (
	"fmt"
	"log"
	"time"

	"example.com/go_cards_server/gametypes"
	"example.com/go_cards_server/player"
	"example.com/go_cards_server/session"
	"example.com/go_cards_server/war"
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
	case gametypes.War:
		var warGame = war.CreateNewWarSession(session)
		log.Println("Created new War session ", warGame) // TODO store this somewhere so we can track it
	default:
		log.Println("Unknown game type")
	}

	return session
}

func HandlePlayerJoined(s *session.Session, player *player.Player) {
	s.AddPlayerToSession(player)
}

func sessionCleaner() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		for sessionId, session := range Sessions {
			if time.Since(session.SessionLastMessageTime) > 60*time.Second {
				session.EndSession();
				session.Active = false

				// Ensure each player connection is closed
				for _, player := range session.Players {
					if player.PlayerConnection != nil {
						player.PlayerConnection.WriteMessage(1, []byte("Session ending due to inactivity"))
						player.PlayerConnection.Close()
					}
				}

				fmt.Println("Cleaning session:", sessionId)
				delete(Sessions, sessionId)
			}
		}
	}
}
