package message_models

import (
	"time"

	"example.com/server/models"
	"github.com/google/uuid"
)

// SessionInfoMessage
type SessionStartedMessage struct {
	SessionId              uuid.UUID                    `json:"sessionId"`
	SessionStartTime       time.Time                    `json:"sessionStartTime"`
	SessionLastMessageTime time.Time                    `json:"sessionLastMessageTime"`
	Players                map[uuid.UUID]*models.Player `json:"players"`
	PlayerId               uuid.UUID                    `json:"playerId"`
	GameType               models.GameType              `json:"gameType"`
}

func CreateSessionStartedMessage(s *models.Session, playerId uuid.UUID) SessionStartedMessage {
	return SessionStartedMessage{
		SessionId:              s.SessionId,
		SessionStartTime:       s.SessionStartTime,
		SessionLastMessageTime: s.SessionLastMessageTime,
		GameType:               s.GameType,
		Players:                s.Players,
		PlayerId:               playerId,
	}
}
