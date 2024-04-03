package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	SessionId              uuid.UUID             `json:"sessionId"`
	SessionStartTime       time.Time             `json:"sessionStartTime"`
	SessionLastMessageTime time.Time             `json:"sessionLastMessageTime"`
	Players                map[uuid.UUID]*Player `json:"players"`
	GameChannel            chan *Message         `json:"-"`
}
