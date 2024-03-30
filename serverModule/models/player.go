package models

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	PlayerId         uuid.UUID       `json:"playerId"`
	PlayerName       string          `json:"playerName"`
	PlayerConnection *websocket.Conn `json:"-"`
}
