package session

import (
	"time"

	"example.com/go_cards_server/cards"
	"example.com/go_cards_server/messages"
	"example.com/go_cards_server/player"
	"github.com/google/uuid"
)

func CreateSessionInfoMessage(s *Session) messages.SessionInfoMessage {
	players := make(map[uuid.UUID]string)
	for id, player := range s.Players {
		players[id] = player.PlayerName
	}

	return messages.SessionInfoMessage{
		SessionId:        s.SessionId,
		SessionStartTime: s.SessionStartTime,
		Players:          players,
		MessageInfo: messages.Message{
			MessageType:      messages.SessionInfoMessageType,
			MessageTimestamp: time.Now(),
		},
	}
}

func CreateSessionStartedMessage(s *Session, userId uuid.UUID) messages.SessionStartedMessage {
	return messages.SessionStartedMessage{
		SessionId:        s.SessionId,
		SessionStartTime: s.SessionStartTime,
		Players:          s.Players,
		PlayerId:         userId,
		MessageInfo: messages.Message{
			MessageType:      messages.SessionStartedMessageType,
			MessageTimestamp: time.Now(),
		},
	}
}

func CreateSessionEndedMessage(s *Session) messages.SessionEndedMessage {
	return messages.SessionEndedMessage{
		SessionId: s.SessionId,
		MessageInfo: messages.Message{
			MessageType:      messages.SessionInfoMessageType,
			MessageTimestamp: time.Now(),
		},
	}
}

func CreateGameStartedMessage(s *Session) messages.GameStartedMessage {
	return messages.GameStartedMessage{
		SessionId:        s.SessionId,
		SessionStartTime: s.SessionStartTime,
		Players:          s.Players,
		MessageInfo: messages.Message{
			MessageType:      messages.GameStartedMessageType,
			MessageTimestamp: time.Now(),
		},
	}
}

func CreatePlayerReadyMessage(playerId uuid.UUID, playerReady bool) messages.PlayerReadyMessage {
	return messages.PlayerReadyMessage{
		PlayerId:    playerId.String(),
		PlayerReady: playerReady,
		MessageInfo: messages.Message{
			MessageType:      messages.PlayerReadyMessageType,
			MessageTimestamp: time.Now(),
		},
	}
}

func CreatePlayerJoinedMessage(player *player.Player) messages.PlayerJoinedMessage {
	return messages.PlayerJoinedMessage{
		PlayerId:   player.PlayerId.String(),
		PlayerName: player.PlayerName,
		MessageInfo: messages.Message{
			MessageType:      messages.PlayerJoinedMessageType,
			MessageTimestamp: time.Now(),
		},
	}
}

func CreateCardsPlayedMessage(playerId uuid.UUID, cards []cards.Card, targetId uuid.UUID) messages.CardsPlayedMessage {
	return messages.CardsPlayedMessage{
		PlayerId: playerId,
		Cards:    cards,
		TargetId: targetId,
		MessageInfo: messages.Message{
			MessageType:      messages.CardsPlayedMessageType,
			MessageTimestamp: time.Now(),
		},
	}
}

func CreateCardsDealtMessage(playerId uuid.UUID, cards []cards.Card) messages.CardsDealtMessage {
	return messages.CardsDealtMessage{
		PlayerId: playerId,
		Cards:    cards,
		MessageInfo: messages.Message{
			MessageType:      messages.CardsDealtMessageType,
			MessageTimestamp: time.Now(),
		},
	}
}
