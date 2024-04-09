package war

import (
	"encoding/json"
	"log"

	"example.com/go_cards_server/cards"
	"example.com/go_cards_server/messages"
	"example.com/go_cards_server/session"
)

type War struct {
	GameStarted bool
	Deck        *cards.Deck
	warSession  *session.Session
}

func CreateNewWarSession(s *session.Session) *War {
	s.MaxPlayers = 2
	var w = &War{GameStarted: false, Deck: GetShuffledWarDeck(), warSession: s}
	go w.Run()
	return w
}

func (w *War) Run() {

	for {
		msg := <-w.warSession.GameChannel

		// Check if message is nil
		if msg == nil {
			if !w.warSession.Active {
				break
			}
			continue
		}

		if !w.warSession.ArePlayersReady() {
			log.Println("Not all players are ready")
			continue
		}

		// TODO - IF SESSION IS OVER THEN BREAK

		if !w.GameStarted {
			w.warSession.BroadcastMessage(session.CreateGameStartedMessage(w.warSession))
			w.DealCards()
		}

		switch msg.MessageType {
		case messages.CardsPlayedMessageType:
			w.handleCardsPlayedMessage(msg)
		default:
			log.Println("war: Unknown message type")
		}
	}
}

func (w *War) DealCards() {
	for _, p := range w.warSession.Players {
		p.SendMessage(session.CreateCardsDealtMessage(p.PlayerId, w.Deck.DrawNCards(5)))
	}
}

func (w *War) handleCardsPlayedMessage(msg *messages.Message) error {
	var cardsPlayedMessage *messages.CardsPlayedMessage

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = json.Unmarshal(msgBytes, &cardsPlayedMessage)
	if err != nil {
		return err
	}

	// Update the war game state here

	// Broadcast the card played message to all players
	w.warSession.BroadcastMessage(session.CreateCardsPlayedMessage(cardsPlayedMessage.PlayerId, cardsPlayedMessage.Cards, cardsPlayedMessage.TargetId))

	return nil
}
