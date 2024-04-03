package war

import (
	"example.com/go_cards_server/models"
	"github.com/mitchellh/mapstructure"
)

type WarGameMessage struct {
	PlayerId    string `json:"playerId"`
	PlayerReady bool   `json:"playerReady"`
}

type WarAction struct {
	PlayerId string `json:"playerId"`
	Action   string `json:"action"`
	Card    models.Card `json:"card"`
}

func GetShuffledWarDeck() *models.Deck {
	deck := GetWarDeck()
	deck.Shuffle()
	return deck
}

func RunWarGame(msgChan chan *models.Message) {
	for {
		msg := <-msgChan

		// Check if message is nil
		if msg == nil {
			continue
		}

		// Type assertion to map[string]interface{}
		messageMap, ok := msg.Message.(map[string]interface{})
		if !ok {
			continue
		}

		var warGameMessage WarGameMessage
		if err := mapstructure.Decode(messageMap, &warGameMessage); err != nil {
			continue
		}

		//ProcessWarGameMessage(warGameMessage)
	}
}

// func ProcessWarGameMessage(msg WarGameMessage) {
// 	// Check if player is ready
// 	if msg.PlayerReady {
// 		// Add player to war game
// 		AddPlayerToWarGame(msg.PlayerId)
// 	}
// }