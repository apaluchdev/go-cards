package models

import (
	"time"

	"github.com/google/uuid"
)

type MessageType string

type Message struct {
	MessageType      MessageType `json:"messageType"`
	MessageTimestamp time.Time   `json:"messageTimestamp"`
	Message          any         `json:"message"`
}

const (
	PlayerReadyMessageType    MessageType = "playerReady"
	SessionStartedMessageType MessageType = "sessionStarted"
	SessionInfoMessageType    MessageType = "sessionInfo"
	PlayerJoinedMessageType   MessageType = "playerJoined"
	// Add more message types as needed
)

// SessionInfoMessage
type SessionInfoMessage struct {
	SessionId        uuid.UUID            `json:"sessionId"`
	SessionStartTime time.Time            `json:"sessionStartTime"`
	Players          map[uuid.UUID]string `json:"players"`
}

func CreateMessage(message any, messageType MessageType) Message {
	return Message{
		MessageType:      messageType,
		MessageTimestamp: time.Now(),
		Message:          message,
	}
}

// SessionStartedMessage
// func CreateSessionStartedMessage(message any) Message {
// 	return Message{
// 		MessageType:      SessionStarted,
// 		MessageTimestamp: time.Now(),
// 		Message:          message,
// 	}
// }

// PlayerJoinedMessage
type PlayerReadyMessage struct {
	PlayerId    string `json:"playerId"`
	PlayerReady bool   `json:"playerReady"`
}

// func CreatePlayerReadyMessage(message PlayerReadyMessage) Message {
// 	return Message{
// 		MessageType:      PlayerReady,
// 		MessageTimestamp: time.Now(),
// 		Message:          message,
// 	}
// }
