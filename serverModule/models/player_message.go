package models

// type PlayerMessageType string

// type PlayerMessage struct {
// 	MessageType      PlayerMessageType `json:"type"`
// 	MessageTimestamp time.Time         `json:"messageTimestamp"`
// 	Message          any               `json:"message"`
// }

// const (
// 	PlayerReadyType PlayerMessageType = "playerReady"
// 	TypeB           PlayerMessageType = "TypeB"
// 	TypeC           PlayerMessageType = "TypeC"
// 	// Add more message types as needed
// )

// type PlayerReadyMessage struct {
// 	// Fields specific to TypeA
// 	PlayerId    string `json:"playerId"`
// 	PlayerReady bool   `json:"playerReady"`
// }

// func CreatePlayerReadyMessage(message any, messageType string) PlayerMessage {
// 	return PlayerMessage{
// 		MessageType:      PlayerReadyType,
// 		MessageTimestamp: time.Now(),
// 		Message:          message,
// 	}
// }

type TypeBData struct {
	// Fields specific to TypeB
	Age int `json:"age"`
}

type TypeCData struct {
	// Fields specific to TypeC
	Enabled bool `json:"enabled"`
}
