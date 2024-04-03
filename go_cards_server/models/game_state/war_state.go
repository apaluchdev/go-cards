package models

import "example.com/go_cards_server/models"

// should derive from a game state model (composition)
type WarState struct {
	Deck *models.Deck `json:"cards"`
}
