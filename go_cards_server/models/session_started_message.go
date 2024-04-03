package models

import (
	"time"

	"github.com/google/uuid"
)

// SessionInfoMessage
type SessionStartedMessage struct {
	SessionId              uuid.UUID             `json:"sessionId"`
	SessionStartTime       time.Time             `json:"sessionStartTime"`
	SessionLastMessageTime time.Time             `json:"sessionLastMessageTime"`
	Players                map[uuid.UUID]*Player `json:"players"`
	PlayerId               uuid.UUID             `json:"playerId"`
}

func CreateSessionStartedMessage(s *Session, playerId uuid.UUID) SessionStartedMessage {
	return SessionStartedMessage{
		SessionId:              s.SessionId,
		SessionStartTime:       s.SessionStartTime,
		SessionLastMessageTime: s.SessionLastMessageTime,
		Players:                s.Players,
		PlayerId:               playerId,
	}
}
