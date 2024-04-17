package session

import (
	"encoding/json"
	"fmt"

	"example.com/go_cards_server/messages"
	"example.com/go_cards_server/user"
)

func (s *Session) handleMessage(msg *messages.Message, p *user.User) error {
	switch msg.MessageInfo.MessageType {
	case messages.UserReadyMessageType:
		s.handleUserReadyMessage(msg.MessageBytes, p)
	default:
		fmt.Println("Unknown message type")
	}

	return nil
}

func (s *Session) handleUserReadyMessage(msg []byte, p *user.User) error {
	userReadyMessage := &messages.UserReadyMessage{}

	if err := json.Unmarshal(msg, userReadyMessage); err != nil {
		return err
	}

	p.UserReady = userReadyMessage.UserReady

	s.BroadcastMessage(CreateUserReadyMessage(p.UserId, userReadyMessage.UserReady))

	return nil
}
