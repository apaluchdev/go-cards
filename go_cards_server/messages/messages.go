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
	PlayerTurnMessageType     MessageType = "playerTurn"
	CardsPlayedMessageType    MessageType = "cardsPlayed"
	CardsDealtMessageType     MessageType = "cardsDealt"
	DeclaredCheatMessageType  MessageType = "declaredCheat"
	// Add more message types as needed
)

type MessageInfo struct {
	MessageType      MessageType `json:"messageType"`
	MessageTimestamp time.Time   `json:"messageTimestamp"`
}

type TypedByteMessage struct {
	MessageBytes *[]byte
	MessageType  MessageType
	SentBy       uuid.UUID `json:"-"`
}

type Message struct {
	MessageInfo  MessageInfo `json:"messageInfo"`
	MessageBytes []byte      `json:"-"`
}

func UnmarshalByteMessage(msg []byte) (*Message, error) {
	var message *Message
	err := json.Unmarshal(msg, &message)

	if err != nil {
		return nil, err
	}

	return message, nil
}

type PlayerReadyMessage struct {
	MessageInfo MessageInfo `json:"messageInfo"`
	PlayerId    string      `json:"playerId"`
	PlayerReady bool        `json:"playerReady"`
}

type SessionInfoMessage struct {
	MessageInfo      MessageInfo          `json:"messageInfo"`
	SessionId        uuid.UUID            `json:"sessionId"`
	SessionStartTime time.Time            `json:"sessionStartTime"`
	Players          map[uuid.UUID]string `json:"players"`
}

type SessionStartedMessage struct {
	MessageInfo            MessageInfo                  `json:"messageInfo"`
	SessionId              uuid.UUID                    `json:"sessionId"`
	SessionStartTime       time.Time                    `json:"sessionStartTime"`
	SessionLastMessageTime time.Time                    `json:"sessionLastMessageTime"`
	Players                map[uuid.UUID]*player.Player `json:"players"`
	PlayerId               uuid.UUID                    `json:"playerId"`
}

type SessionEndedMessage struct {
	MessageInfo MessageInfo `json:"messageInfo"`
	SessionId   uuid.UUID   `json:"sessionId"`
}

type GameStartedMessage struct {
	MessageInfo            MessageInfo                  `json:"messageInfo"`
	SessionId              uuid.UUID                    `json:"sessionId"`
	SessionStartTime       time.Time                    `json:"sessionStartTime"`
	SessionLastMessageTime time.Time                    `json:"sessionLastMessageTime"`
	Players                map[uuid.UUID]*player.Player `json:"players"`
	PlayerId               uuid.UUID                    `json:"playerId"`
}

type CardsPlayedMessage struct {
	MessageInfo MessageInfo  `json:"messageInfo"`
	PlayerId    string       `json:"playerId"`
	Cards       []cards.Card `json:"cards"`
	TargetId    string       `json:"targetId"`
}

type PlayerJoinedMessage struct {
	MessageInfo MessageInfo `json:"messageInfo"`
	PlayerId    string      `json:"playerId"`
	PlayerName  string      `json:"playerName"`
}

type PlayerTurnMessage struct {
	MessageInfo MessageInfo `json:"messageInfo"`
	PlayerId    string      `json:"playerId"`
}

type DeclaredCheatMessage struct {
	PlayerId    string      `json:"playerId"`
	MessageInfo MessageInfo `json:"messageInfo"`
}

type CardsDealtMessage struct {
	MessageInfo MessageInfo  `json:"messageInfo"`
	PlayerId    uuid.UUID    `json:"playerId"`
	Cards       []cards.Card `json:"cards"`
}
