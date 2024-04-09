package messages

import (
	"encoding/json"
	"time"

	"example.com/go_cards_server/cards"
	"example.com/go_cards_server/player"
	"github.com/google/uuid"
)

type MessageType string

const (
	PlayerReadyMessageType    MessageType = "playerReady"
	SessionStartedMessageType MessageType = "sessionStarted"
	GameStartedMessageType    MessageType = "gameStarted"
	SessionInfoMessageType    MessageType = "sessionInfo"
	SessionEndedMessageType   MessageType = "sessionEnded"
	PlayerJoinedMessageType   MessageType = "playerJoined"
	GameMessageType           MessageType = "gameMessage"
	CardsPlayedMessageType    MessageType = "cardsPlayed"
	CardsDealtMessageType     MessageType = "cardsDealt"
	// Add more message types as needed
)

type Message struct {
	MessageType      MessageType `json:"messageType"`
	MessageTimestamp time.Time   `json:"messageTimestamp"`
}

// func CreateMessage(message any, messageType MessageType) Message {
// 	return Message{
// 		MessageType:      messageType,
// 		MessageTimestamp: time.Now(),
// 		Message:          message,
// 	}
// }

func UnmarshalByteMessage(msg []byte) (*Message, error) {
	var message *Message
	err := json.Unmarshal(msg, &message)

	if err != nil {
		return nil, err
	}

	return message, nil
}

type PlayerReadyMessage struct {
	MessageInfo Message `json:"messageInfo"`
	PlayerId    string  `json:"playerId"`
	PlayerReady bool    `json:"playerReady"`
}

type SessionInfoMessage struct {
	MessageInfo      Message              `json:"messageInfo"`
	SessionId        uuid.UUID            `json:"sessionId"`
	SessionStartTime time.Time            `json:"sessionStartTime"`
	Players          map[uuid.UUID]string `json:"players"`
}

type SessionStartedMessage struct {
	MessageInfo            Message                      `json:"messageInfo"`
	SessionId              uuid.UUID                    `json:"sessionId"`
	SessionStartTime       time.Time                    `json:"sessionStartTime"`
	SessionLastMessageTime time.Time                    `json:"sessionLastMessageTime"`
	Players                map[uuid.UUID]*player.Player `json:"players"`
	PlayerId               uuid.UUID                    `json:"playerId"`
}

type SessionEndedMessage struct {
	MessageInfo Message   `json:"messageInfo"`
	SessionId   uuid.UUID `json:"sessionId"`
}

type GameStartedMessage struct {
	MessageInfo            Message                      `json:"messageInfo"`
	SessionId              uuid.UUID                    `json:"sessionId"`
	SessionStartTime       time.Time                    `json:"sessionStartTime"`
	SessionLastMessageTime time.Time                    `json:"sessionLastMessageTime"`
	Players                map[uuid.UUID]*player.Player `json:"players"`
	PlayerId               uuid.UUID                    `json:"playerId"`
}

type CardsPlayedMessage struct {
	MessageInfo Message      `json:"messageInfo"`
	PlayerId    uuid.UUID    `json:"playerId"`
	Cards       []cards.Card `json:"cards"`
	TargetId    uuid.UUID    `json:"targetId"`
}

type PlayerJoinedMessage struct {
	MessageInfo Message `json:"messageInfo"`
	PlayerId    string  `json:"playerId"`
	PlayerName  string  `json:"playerName"`
}

type CardsDealtMessage struct {
	MessageInfo Message      `json:"messageInfo"`
	PlayerId    uuid.UUID    `json:"playerId"`
	Cards       []cards.Card `json:"cards"`
}
