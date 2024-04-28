package cheat

import (
	"time"

	"example.com/go_cards_server/cards"
	"example.com/go_cards_server/messages"
	"example.com/go_cards_server/session"
	"github.com/google/uuid"
)

func GetShuffledCheatDeck() *cards.Deck {
	deck := GetCheatDeck()
	deck.Shuffle()
	return deck
}

func GetCheatDeck() *cards.Deck {
	var cheatCards []cards.Card = []cards.Card{
		{Suit: "Hearts", Value: "2"},
		{Suit: "Hearts", Value: "3"},
		{Suit: "Hearts", Value: "4"},
		{Suit: "Hearts", Value: "5"},
		{Suit: "Hearts", Value: "6"},
		{Suit: "Hearts", Value: "7"},
		{Suit: "Hearts", Value: "8"},
		{Suit: "Hearts", Value: "9"},
		{Suit: "Hearts", Value: "10"},
		{Suit: "Hearts", Value: "J"},
		{Suit: "Hearts", Value: "Q"},
		{Suit: "Hearts", Value: "K"},
		{Suit: "Hearts", Value: "A"},
		{Suit: "Diamonds", Value: "2"},
		{Suit: "Diamonds", Value: "3"},
		{Suit: "Diamonds", Value: "4"},
		{Suit: "Diamonds", Value: "5"},
		{Suit: "Diamonds", Value: "6"},
		{Suit: "Diamonds", Value: "7"},
		{Suit: "Diamonds", Value: "8"},
		{Suit: "Diamonds", Value: "9"},
		{Suit: "Diamonds", Value: "10"},
		{Suit: "Diamonds", Value: "J"},
		{Suit: "Diamonds", Value: "Q"},
		{Suit: "Diamonds", Value: "K"},
		{Suit: "Diamonds", Value: "A"},
		{Suit: "Clubs", Value: "2"},
		{Suit: "Clubs", Value: "3"},
		{Suit: "Clubs", Value: "4"},
		{Suit: "Clubs", Value: "5"},
		{Suit: "Clubs", Value: "6"},
		{Suit: "Clubs", Value: "7"},
		{Suit: "Clubs", Value: "8"},
		{Suit: "Clubs", Value: "9"},
		{Suit: "Clubs", Value: "10"},
		{Suit: "Clubs", Value: "J"},
		{Suit: "Clubs", Value: "Q"},
		{Suit: "Clubs", Value: "K"},
		{Suit: "Clubs", Value: "A"},
		{Suit: "Spades", Value: "2"},
		{Suit: "Spades", Value: "3"},
		{Suit: "Spades", Value: "4"},
		{Suit: "Spades", Value: "5"},
		{Suit: "Spades", Value: "6"},
		{Suit: "Spades", Value: "7"},
		{Suit: "Spades", Value: "8"},
		{Suit: "Spades", Value: "9"},
		{Suit: "Spades", Value: "10"},
		{Suit: "Spades", Value: "J"},
		{Suit: "Spades", Value: "Q"},
		{Suit: "Spades", Value: "K"},
		{Suit: "Spades", Value: "A"},
	}
	return &cards.Deck{Cards: cheatCards}
}

func CreateDeclaredCheatMessage(playerId uuid.UUID) DeclaredCheatMessage {
	return DeclaredCheatMessage{
		PlayerId: playerId.String(),
		MessageInfo: messages.MessageInfo{
			MessageType:      messages.DeclaredCheatMessageType,
			MessageTimestamp: time.Now(),
		},
	}
}

func CreatePlayerTurnMessage(playerId string, playerInstruction string, discardPileSize int) PlayerTurnMessage {
	return PlayerTurnMessage{
		PlayerId:          playerId,
		PlayerInstruction: playerInstruction,
		DiscardPileSize:   discardPileSize,
		MessageInfo: messages.MessageInfo{
			MessageType:      messages.PlayerTurnMessageType,
			MessageTimestamp: time.Now(),
		},
	}
}

func CreateCardsPlayedMessage(playerId string, cards []cards.Card, targetId string) CardsPlayedMessage {
	return CardsPlayedMessage{
		PlayerId: playerId,
		Cards:    cards,
		TargetId: targetId,
		MessageInfo: messages.MessageInfo{
			MessageType:      messages.CardsPlayedMessageType,
			MessageTimestamp: time.Now(),
		},
	}
}

func CreateCardsDealtMessage(playerId uuid.UUID, cards []cards.Card) CardsDealtMessage {
	return CardsDealtMessage{
		PlayerId: playerId,
		Cards:    cards,
		MessageInfo: messages.MessageInfo{
			MessageType:      messages.CardsDealtMessageType,
			MessageTimestamp: time.Now(),
		},
	}
}

func CreateGameStartedMessage(s *session.Session) GameStartedMessage {
	return GameStartedMessage{
		SessionId:        s.SessionId,
		SessionStartTime: s.SessionStartTime,
		Players:          s.Users,
		MessageInfo: messages.MessageInfo{
			MessageType:      messages.GameStartedMessageType,
			MessageTimestamp: time.Now(),
		},
	}
}

func CreateCheatResultMessage(winnerId string, loserId string) CheatResultMessage {
	return CheatResultMessage{
		WinnerId: winnerId,
		LoserId:  loserId,
		MessageInfo: messages.MessageInfo{
			MessageType:      messages.CheatResultMessageType,
			MessageTimestamp: time.Now(),
		},
	}
}
