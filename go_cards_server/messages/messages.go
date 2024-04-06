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
	SessionInfoMessageType    MessageType = "sessionInfo"
	PlayerJoinedMessageType   MessageType = "playerJoined"
	GameMessageType           MessageType = "gameMessage"
	CardPlayedMessageType     MessageType = "cardPlayed"
	CardDealedMessageType     MessageType = "cardDealed"
	// Add more message types as needed
)

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

func UnmarshalByteMessage(msg []byte) (*Message, error) {
	var message *Message
	err := json.Unmarshal(msg, &message)

	if err != nil {
		return nil, err
	}

	return message, nil
}

type PlayerReadyMessage struct {
	PlayerId    uuid.UUID `json:"playerId"`
	PlayerReady bool      `json:"playerReady"`
}

type SessionInfoMessage struct {
	SessionId        uuid.UUID            `json:"sessionId"`
	SessionStartTime time.Time            `json:"sessionStartTime"`
	Players          map[uuid.UUID]string `json:"players"`
}

// SessionInfoMessage
type SessionStartedMessage struct {
	SessionId              uuid.UUID                    `json:"sessionId"`
	SessionStartTime       time.Time                    `json:"sessionStartTime"`
	SessionLastMessageTime time.Time                    `json:"sessionLastMessageTime"`
	Players                map[uuid.UUID]*player.Player `json:"players"`
	PlayerId               uuid.UUID                    `json:"playerId"`
}

type CardPlayedMessage struct {
	PlayerId uuid.UUID  `json:"playerId"`
	Card     cards.Card `json:"card"`
	TargetId uuid.UUID  `json:"targetId"`
}

type CardDealedMessage struct {
	PlayerId uuid.UUID    `json:"playerId"`
	Card     []cards.Card `json:"card"`
}
