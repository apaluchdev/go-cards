package message_models

// PlayerJoinedMessage
type PlayerReadyMessage struct {
	PlayerId    string `json:"playerId"`
	PlayerReady bool   `json:"playerReady"`
}
