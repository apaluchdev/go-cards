package models

import (
	"time"

	"github.com/google/uuid"
)

// SessionInfoMessage
type SessionInfoMessage struct {
	SessionId        uuid.UUID            `json:"sessionId"`
	SessionStartTime time.Time            `json:"sessionStartTime"`
	Players          map[uuid.UUID]string `json:"players"`
}
