package session

import (
	"encoding/json"
	"fmt"

	"example.com/go_cards_server/messages"
	"example.com/go_cards_server/player"
)

func (s *Session) handleMessage(msg *messages.Message, p *player.Player) error {
	switch msg.MessageType {
	case messages.PlayerReadyMessageType:
		s.handlePlayerReadyMessage(msg, p)
	default:
		fmt.Println("Unknown message type")
	}

	return nil
}

func (s *Session) handlePlayerReadyMessage(msg *messages.Message, p *player.Player) error {
	var playerReadyMessage *messages.PlayerReadyMessage

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = json.Unmarshal(msgBytes, &playerReadyMessage)
	if err != nil {
		return err
	}

	p.PlayerReady = playerReadyMessage.PlayerReady

	s.BroadcastMessage(CreatePlayerReadyMessage(p.PlayerId, playerReadyMessage.PlayerReady))

	return nil
}
