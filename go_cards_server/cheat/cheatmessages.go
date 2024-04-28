package cheat

import (
	"time"

	"example.com/go_cards_server/cards"
	"example.com/go_cards_server/messages"
	"example.com/go_cards_server/user"
	"github.com/google/uuid"
)

type CheatMessageType string

type GameStartedMessage struct {
	MessageInfo            messages.MessageInfo     `json:"messageInfo"`
	SessionId              uuid.UUID                `json:"sessionId"`
	SessionStartTime       time.Time                `json:"sessionStartTime"`
	SessionLastMessageTime time.Time                `json:"sessionLastMessageTime"`
	Players                map[uuid.UUID]*user.User `json:"players"`
	PlayerId               uuid.UUID                `json:"playerId"`
}

type PlayerTurnMessage struct {
	MessageInfo       messages.MessageInfo `json:"messageInfo"`
	PlayerInstruction string               `json:"playerInstruction"`
	PlayerId          string               `json:"playerId"`
	DiscardPileSize   int                  `json:"discardPileSize"`
}

type DeclaredCheatMessage struct {
	PlayerId    string               `json:"playerId"`
	MessageInfo messages.MessageInfo `json:"messageInfo"`
}

type CardsDealtMessage struct {
	MessageInfo messages.MessageInfo `json:"messageInfo"`
	PlayerId    uuid.UUID            `json:"playerId"`
	Cards       []cards.Card         `json:"cards"`
}

type CardsPlayedMessage struct {
	MessageInfo messages.MessageInfo `json:"messageInfo"`
	PlayerId    string               `json:"playerId"`
	Cards       []cards.Card         `json:"cards"`
	TargetId    string               `json:"targetId"`
}

type CheatResultMessage struct {
	MessageInfo messages.MessageInfo `json:"messageInfo"`
	WinnerId    string               `json:"winnerId"`
	LoserId     string               `json:"loserId"`
}
