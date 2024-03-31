package session_manager

import (
	"log"

	"example.com/server/models"
	"example.com/server/models/message_models"
	"github.com/mitchellh/mapstructure"
)

func handlePlayerReadyMessage(s *models.Session, m map[string]interface{}, p *models.Player) error {
	var playerReadyMessage message_models.PlayerReadyMessage
	if err := mapstructure.Decode(m, &playerReadyMessage); err != nil {
		return err
	}
	playerReadyMessage.PlayerId = p.PlayerId.String()
	log.Println("Successfully retrieved PlayerReadyType ", playerReadyMessage.PlayerId, playerReadyMessage.PlayerReady)
	BroadcastMessage(s.SessionId, message_models.CreateMessage(playerReadyMessage, message_models.PlayerReadyMessageType))

	return nil
}
