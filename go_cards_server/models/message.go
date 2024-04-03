package models

import (
	"time"
)

const (
	PlayerReadyMessageType    MessageType = "playerReady"
	SessionStartedMessageType MessageType = "sessionStarted"
	SessionInfoMessageType    MessageType = "sessionInfo"
	PlayerJoinedMessageType   MessageType = "playerJoined"
	GameMessageType           MessageType = "gameMessage"
	// Add more message types as needed
)

type MessageType string

type Message struct {
	MessageType      MessageType `json:"messageType"`
	MessageTimestamp time.Time   `json:"messageTimestamp"`
	Message          any         `json:"message"`
}

func CreateMessage(message any, messageType MessageType) Message {
	return Message{
		MessageType:      messageType,
		MessageTimestamp: time.Now(),
		Message:          message,
	}
}
