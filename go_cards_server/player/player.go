package player

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	PlayerId         uuid.UUID       `json:"playerId"`
	PlayerName       string          `json:"playerName"`
	PlayerReady      bool            `json:"playerReady"`
	PlayerConnection *websocket.Conn `json:"-"`
}

func (p *Player) SendMessage(message any) error {
	if p.PlayerConnection == nil {
		return errors.New("Player connection is nil")
	}

	// Marshal the message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling message ", err)
		return err
	}

	err = p.PlayerConnection.WriteMessage(websocket.TextMessage, messageJSON)

	// If there was an error writing to the client, return false
	if err != nil {
		fmt.Println("Write message error:", err)
		if err == websocket.ErrCloseSent {
			fmt.Println("Breaking connection")
			return err
		}
	}

	return nil
}
