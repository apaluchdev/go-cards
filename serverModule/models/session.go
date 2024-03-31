package models

import (
	"time"

	"github.com/google/uuid"
)

type GameType string

const (
	War   GameType = "War"
	Cheat GameType = "Cheat"
	Poker GameType = "Poker"
)

type Session struct {
	SessionId              uuid.UUID             `json:"sessionId"`
	SessionStartTime       time.Time             `json:"sessionStartTime"`
	SessionLastMessageTime time.Time             `json:"sessionLastMessageTime"`
	Players                map[uuid.UUID]*Player `json:"players"`
	GameType               GameType              `json:"gameType"`
	Deck                   *Deck                  `json:"cards"`
}
