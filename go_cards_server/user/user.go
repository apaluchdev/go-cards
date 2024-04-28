package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type User struct {
	UserId         uuid.UUID       `json:"userId"`
	UserName       string          `json:"userName"`
	UserReady      bool            `json:"userReady"`
	UserConnection *websocket.Conn `json:"-"`
}

func (p *User) SendMessage(message any) error {
	log.Println("Sending message: ", message)
	if p.UserConnection == nil {
		return errors.New("User connection is nil")
	}

	// Marshal the message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling message ", err)
		return err
	}

	err = p.UserConnection.WriteMessage(websocket.TextMessage, messageJSON)

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
