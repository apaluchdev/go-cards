package session

import (
	"encoding/json"
	"fmt"

	"example.com/go_cards_server/messages"
	"example.com/go_cards_server/player"
)

func (s *Session) handleMessage(msg *messages.Message, p *player.Player) error {
	switch msg.MessageInfo.MessageType {
	case messages.PlayerReadyMessageType:
		s.handlePlayerReadyMessage(msg.MessageBytes, p)
	case messages.CardsPlayedMessageType:
		// s.handleCardsPlayedMessage(msg)
	default:
		fmt.Println("Unknown message type")
	}

	return nil
}

func (s *Session) handlePlayerReadyMessage(msg []byte, p *player.Player) error {
	playerReadyMessage := &messages.PlayerReadyMessage{}

	if err := json.Unmarshal(msg, playerReadyMessage); err != nil {
		return err
	}

	p.PlayerReady = playerReadyMessage.PlayerReady

	s.BroadcastMessage(CreatePlayerReadyMessage(p.PlayerId, playerReadyMessage.PlayerReady))

	return nil
}
