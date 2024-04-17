package messages

import (
	"encoding/json"
	"time"

	"example.com/go_cards_server/user"
	"github.com/google/uuid"
)

type MessageType string

const (
	UserReadyMessageType      MessageType = "userReady"
	SessionStartedMessageType MessageType = "sessionStarted"
	GameStartedMessageType    MessageType = "gameStarted"
	SessionInfoMessageType    MessageType = "sessionInfo"
	SessionEndedMessageType   MessageType = "sessionEnded"
	UserJoinedMessageType     MessageType = "userJoined"
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

type UserReadyMessage struct {
	MessageInfo MessageInfo `json:"messageInfo"`
	UserId      string      `json:"userId"`
	UserReady   bool        `json:"userReady"`
}

type SessionInfoMessage struct {
	MessageInfo      MessageInfo          `json:"messageInfo"`
	SessionId        uuid.UUID            `json:"sessionId"`
	SessionStartTime time.Time            `json:"sessionStartTime"`
	Users            map[uuid.UUID]string `json:"users"`
}

type SessionStartedMessage struct {
	MessageInfo            MessageInfo              `json:"messageInfo"`
	SessionId              uuid.UUID                `json:"sessionId"`
	SessionStartTime       time.Time                `json:"sessionStartTime"`
	SessionLastMessageTime time.Time                `json:"sessionLastMessageTime"`
	Users                  map[uuid.UUID]*user.User `json:"users"`
	UserId                 uuid.UUID                `json:"userId"`
}

type SessionEndedMessage struct {
	MessageInfo MessageInfo `json:"messageInfo"`
	SessionId   uuid.UUID   `json:"sessionId"`
}

type UserJoinedMessage struct {
	MessageInfo MessageInfo `json:"messageInfo"`
	UserId      string      `json:"userId"`
	UserName    string      `json:"userName"`
}
