package session

import (
	"time"

	"example.com/go_cards_server/cards"
	"example.com/go_cards_server/messages"
	"example.com/go_cards_server/player"
	"github.com/google/uuid"
)

type Session struct {
	SessionId              uuid.UUID                    `json:"sessionId"`
	SessionStartTime       time.Time                    `json:"sessionStartTime"`
	SessionLastMessageTime time.Time                    `json:"sessionLastMessageTime"`
	Players                map[uuid.UUID]*player.Player `json:"players"`
	GameChannel            chan *messages.Message       `json:"-"`
}

func CreateSessionInfoMessage(s *Session) messages.Message {
	players := make(map[uuid.UUID]string)
	for id, player := range s.Players {
		players[id] = player.PlayerName
	}

	return messages.CreateMessage(messages.SessionInfoMessage{
		SessionId:        s.SessionId,
		SessionStartTime: s.SessionStartTime,
		Players:          players,
	}, messages.SessionInfoMessageType)
}

func CreateSessionStartedMessage(s *Session, userId uuid.UUID) messages.Message {
	return messages.CreateMessage(messages.SessionStartedMessage{
		SessionId:        s.SessionId,
		SessionStartTime: s.SessionStartTime,
		Players:          s.Players,
		PlayerId:         userId,
	}, messages.SessionStartedMessageType)
}

func CreatePlayerReadyMessage(playerId uuid.UUID, playerReady bool) messages.Message {
	return messages.CreateMessage(messages.PlayerReadyMessage{
		PlayerId:    playerId,
		PlayerReady: playerReady,
	}, messages.PlayerReadyMessageType)
}

func CreateCardPlayedMessage(playerId uuid.UUID, card cards.Card) messages.Message {
	return messages.CreateMessage(messages.CardPlayedMessage{
		PlayerId:    playerId,
		Card: card,
	}, messages.PlayerReadyMessageType)
}

func CreateCardDealedMessage(playerId uuid.UUID, cards []cards.Card) messages.Message {
	return messages.CreateMessage(messages.CardDealedMessage{
		PlayerId:    playerId,
		Card: cards,
	}, messages.PlayerReadyMessageType)
}

