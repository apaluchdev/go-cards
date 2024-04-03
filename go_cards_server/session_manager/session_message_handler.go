package session_manager

import (
	"example.com/go_cards_server/models"
	"github.com/mitchellh/mapstructure"
)

func handlePlayerReadyMessage(s *models.Session, msg *models.Message, p *models.Player) error {
	messageMap := msg.Message.(map[string]interface{})
	var playerReadyMessage models.PlayerReadyMessage
	if err := mapstructure.Decode(messageMap, &playerReadyMessage); err != nil {
		return err
	}

	// Broadcast player ready message
	playerReadyMessage.PlayerId = p.PlayerId.String()
	p.PlayerReady = playerReadyMessage.PlayerReady
	BroadcastMessage(s.SessionId, models.CreateMessage(playerReadyMessage, models.PlayerReadyMessageType))

	return nil
}
