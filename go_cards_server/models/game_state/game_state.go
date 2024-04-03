package models

import "example.com/go_cards_server/models"

type GameState struct {
	Deck *models.Deck `json:"cards"`
}
