package cheat

import (
	"encoding/json"
	"log"
	"time"

	"example.com/go_cards_server/cards"
	"example.com/go_cards_server/messages"
	"example.com/go_cards_server/user"
	"example.com/go_cards_server/session"
	"github.com/google/uuid"
)

type CheatPlayer struct {
	User *user.User
	Hand   []cards.Card
}

type Cheat struct {
	GameStarted                   bool
	Deck                          *cards.Deck
	cheatSession                  *session.Session
	userTurn                    uuid.UUID
	turnCounter                   int
	isWaitingForCheatDeclarations bool
	CheatPlayers                  []*CheatPlayer
	PlayedCards                   []cards.Card
	ExpectedValue                 string
	DiscardPile                   []cards.Card
}

func CreateNewCheatSession(s *session.Session) *Cheat {
	s.MaxUsers = 2
	var c = &Cheat{GameStarted: false, Deck: GetShuffledCheatDeck(), cheatSession: s, isWaitingForCheatDeclarations: false, userTurn: uuid.Nil, turnCounter: 0, CheatPlayers: make([]*CheatPlayer, 0)}
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
	// Broadcast the card played message to all users
	c.cheatSession.BroadcastMessage(CreateCardsPlayedMessage(typedByteMessage.SentBy.String(), cardsPlayedMessage.Cards, cardsPlayedMessage.TargetId))

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
		if card.Value != c.ExpectedValue {
			caughtCheat = true
			break
		}
	}

	if caughtCheat {
		log.Println("Caught a cheat!")
		// Add all cards including discard to the user who cheated

	} else {
		log.Println("Bad call!")
		// Add all cards to the user who made the bad call
	}

	// Broadcast if declared cheat was successful

	timer := time.NewTimer(3 * time.Second)
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
	c.cheatSession.BroadcastMessage(CreatePlayerTurnMessage(c.userTurn.String()))
}

func (c *Cheat) GetNextUserTurn() {
	c.turnCounter++
	c.turnCounter = c.turnCounter % len(c.CheatPlayers)
	c.userTurn = c.CheatPlayers[c.turnCounter].User.UserId
	c.cheatSession.BroadcastMessage(CreatePlayerTurnMessage(c.userTurn.String()))
}

func (c *Cheat) SetMaxWaitTimeForCheatDeclarations() {
	timer := time.NewTimer(10 * time.Second)
	<-timer.C

	c.EndCheatDeclaration()
}

func (c *Cheat) EndCheatDeclaration() {
	c.isWaitingForCheatDeclarations = false
	log.Println("Cheat declaration over")
	c.GetNextUserTurn()
}
