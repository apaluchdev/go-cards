package session

import (
	"fmt"

	"example.com/go_cards_server/messages"
	"example.com/go_cards_server/player"
	"github.com/mitchellh/mapstructure"
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
	messageMap := msg.Message.(map[string]interface{})
	var playerReadyMessage messages.PlayerReadyMessage
	if err := mapstructure.Decode(messageMap, &playerReadyMessage); err != nil {
		return err
	}

	s.BroadcastMessage(CreatePlayerReadyMessage(p.PlayerId, playerReadyMessage.PlayerReady))

	return nil
}
