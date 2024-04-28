package cheat

import (
	"encoding/json"
	"log"
	"time"

	"example.com/go_cards_server/cards"
	"example.com/go_cards_server/messages"
	"example.com/go_cards_server/session"
	"example.com/go_cards_server/user"
	"github.com/google/uuid"
)

type CheatPlayer struct {
	User *user.User
	Hand []cards.Card
}

type Cheat struct {
	GameStarted                   bool
	Deck                          *cards.Deck
	cheatSession                  *session.Session
	userTurn                      uuid.UUID
	turnCounter                   int
	isWaitingForCheatDeclarations bool
	CheatPlayers                  []*CheatPlayer
	PlayedCards                   []cards.Card
	CurrentCardValueIndex         int
	DiscardPile                   []cards.Card
	CurrentValue                  string
}

func CreateNewCheatSession(s *session.Session) *Cheat {
	s.MaxUsers = 2
	var c = &Cheat{GameStarted: false, Deck: GetShuffledCheatDeck(), cheatSession: s, isWaitingForCheatDeclarations: false, userTurn: uuid.Nil, turnCounter: 0, CurrentCardValueIndex: 0, CheatPlayers: make([]*CheatPlayer, 0)}
	go c.Run()
	return c
}

func (w *Cheat) Run() {

	for {
		// TODO - End game session once session is cleaned up
		msg := <-w.cheatSession.GameChannel

		// Check if message is nil
		if msg == nil {
			if !w.cheatSession.Active {
				break
			}
			continue
		}

		if !w.GameStarted {
			w.HandleGameStart()
			continue
		} else {
			if w.userTurn == uuid.Nil {
				w.GetNextUserTurn()
			}
		}

		for _, p := range w.CheatPlayers {
			log.Println("Player ", p.User.UserName, " should have ", string(len(p.Hand)), " cards.")
		}

		switch msg.MessageType {
		case messages.CardsPlayedMessageType:
			w.handleCardsPlayedMessage(msg)
		case messages.DeclaredCheatMessageType:
			w.handleDeclaredCheatMessage(msg)
		default:
			log.Println("war: Skipping message type: ", msg.MessageType)
		}
	}
}

func (w *Cheat) DealCards() {
	for _, p := range w.CheatPlayers {
		cards := w.Deck.DrawNCards(5)
		p.Hand = append(p.Hand, cards...)
		p.User.SendMessage(CreateCardsDealtMessage(p.User.UserId, cards))
	}
}

func (c *Cheat) handleCardsPlayedMessage(typedByteMessage *messages.TypedByteMessage) error {
	var cardsPlayedMessage *CardsPlayedMessage

	err := json.Unmarshal(*typedByteMessage.MessageBytes, &cardsPlayedMessage)
	if err != nil {
		return err
	}

	// Verify this message can be processed given the current state of the game
	if typedByteMessage.SentBy != c.userTurn || c.isWaitingForCheatDeclarations {
		log.Println("Cannot accept card played message currently")
		return nil
	}

	c.PlayedCards = cardsPlayedMessage.Cards

	// Validate the played cards
	if len(c.PlayedCards) == 0 {
		log.Println("No cards played")
		return nil
	}

	// Move all the played cards from the player's hand to the discard pile
	updatedHand := make([]cards.Card, 0)
	for _, card := range c.PlayedCards {
		for _, handCard := range c.CheatPlayers[c.turnCounter].Hand {
			if handCard.Suit == card.Suit && handCard.Value == card.Value {
				c.DiscardPile = append(c.DiscardPile, card)
			} else {
				updatedHand = append(updatedHand, handCard)
			}
		}
	}
	c.CheatPlayers[c.turnCounter].Hand = updatedHand

	// Create the "hidden cards" to show the other players
	cheatCards := make([]cards.Card, len(c.PlayedCards))
	for i, _ := range c.PlayedCards {
		cheatCards[i] = cards.Card{Suit: "Maybe", Value: "Maybe"}
	}

	// Broadcast the card played message to all users
	c.cheatSession.BroadcastMessage(CreateCardsPlayedMessage(typedByteMessage.SentBy.String(), cheatCards, cardsPlayedMessage.TargetId))

	c.isWaitingForCheatDeclarations = true
	go c.SetMaxWaitTimeForCheatDeclarations()

	return nil
}

func (c *Cheat) handleDeclaredCheatMessage(typedByteMessage *messages.TypedByteMessage) {
	if !c.isWaitingForCheatDeclarations {
		log.Println("Cannot accept cheat declaration message currently")
		return
	}
	c.isWaitingForCheatDeclarations = false
	c.userTurn = uuid.Nil
	c.cheatSession.BroadcastMessage(CreateDeclaredCheatMessage(typedByteMessage.SentBy))

	// Cheat just has to check if the played cards are of the correct value, REDO THE BELOW
	caughtCheat := false
	for _, card := range c.PlayedCards {
		curValue := cards.CardValues[c.CurrentCardValueIndex]
		if card.Value != curValue {
			caughtCheat = true
			break
		}
	}

	// Give the discard pile to the appropriate player
	if caughtCheat {
		c.CheatPlayers[c.turnCounter].Hand = append(c.CheatPlayers[c.turnCounter].Hand, c.DiscardPile...)

		// Send a message to the cheater to pickup the discard pile
		c.CheatPlayers[c.turnCounter].User.SendMessage(CreateCardsDealtMessage(c.CheatPlayers[c.turnCounter].User.UserId, c.DiscardPile))

		// Broadcast the cheat result message
		c.cheatSession.BroadcastMessage(CreateCheatResultMessage(typedByteMessage.SentBy.String(), c.CheatPlayers[c.turnCounter].User.UserId.String()))

	} else {
		for _, cheatPlayer := range c.CheatPlayers {
			if cheatPlayer.User.UserId == typedByteMessage.SentBy {
				cheatPlayer.Hand = append(cheatPlayer.Hand, c.DiscardPile...)

				// Send a message to the accuser to pickup the discard pile
				cheatPlayer.User.SendMessage(CreateCardsDealtMessage(cheatPlayer.User.UserId, c.DiscardPile))

				// Broadcast the cheat result message
				c.cheatSession.BroadcastMessage(CreateCheatResultMessage(c.CheatPlayers[c.turnCounter].User.UserId.String(), typedByteMessage.SentBy.String()))
				break
			}
		}
	}

	// Empty the discard pile, should always be empty after a cheat declaration
	c.DiscardPile = make([]cards.Card, 0)

	timer := time.NewTimer(5 * time.Second)
	<-timer.C

	c.EndCheatDeclaration()
}

func (c *Cheat) HandleGameStart() {
	if !c.cheatSession.AreUsersReady() {
		log.Println("Not all users are ready")
		return
	}

	c.GameStarted = true

	// Create a CheatPlayer for each user in the session
	for _, user := range c.cheatSession.Users {
		c.CheatPlayers = append(c.CheatPlayers, &CheatPlayer{User: user, Hand: make([]cards.Card, 0)})
	}

	// Broadcast the game started message
	c.cheatSession.BroadcastMessage(CreateGameStartedMessage(c.cheatSession))

	c.DealCards()

	c.userTurn = c.CheatPlayers[0].User.UserId
	c.cheatSession.BroadcastMessage(CreatePlayerTurnMessage(c.userTurn.String(), "Play one or more "+cards.CardValues[c.CurrentCardValueIndex]+"'s", len(c.DiscardPile)))
}

func (c *Cheat) GetNextUserTurn() {
	c.CurrentCardValueIndex = (c.CurrentCardValueIndex + 1) % len(cards.CardValues)
	c.turnCounter++
	c.turnCounter = c.turnCounter % len(c.CheatPlayers)
	c.userTurn = c.CheatPlayers[c.turnCounter].User.UserId
	c.cheatSession.BroadcastMessage(CreatePlayerTurnMessage(c.userTurn.String(), "Play one or more "+cards.CardValues[c.CurrentCardValueIndex]+"'s", len(c.DiscardPile)))
}

func (c *Cheat) SetMaxWaitTimeForCheatDeclarations() {
	timer := time.NewTimer(9 * time.Second)
	<-timer.C

	if c.isWaitingForCheatDeclarations {
		c.EndCheatDeclaration()
	}
}

func (c *Cheat) EndCheatDeclaration() {
	c.isWaitingForCheatDeclarations = false
	log.Println("Cheat declaration over")
	c.GetNextUserTurn()
}
