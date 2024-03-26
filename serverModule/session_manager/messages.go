package session_manager

import "time"

type SessionMessage struct {
	MessageType      string    `json:"messageType"`
	MessageTimestamp time.Time `json:"messageTimestamp"`
}

type StartSessionMessage struct {
	SessionMessage
	Message Session `json:"message"`
}

type SessionInfoMessage struct {
	SessionMessage
	Message Session `json:"message"`
}

type PlayerJoinedMessage struct {
	SessionMessage
	Message Player `json:"message"`
}

func CreateStartSessionMessage(session *Session) StartSessionMessage {
	return StartSessionMessage{
		SessionMessage: SessionMessage{
			MessageType:      "sessionStarted",
			MessageTimestamp: time.Now(),
		},
		Message: *session,
	}
}

func CreatePlayerJoinedMessage(player *Player) PlayerJoinedMessage {
	return PlayerJoinedMessage{
		SessionMessage: SessionMessage{
			MessageType:      "playerJoined",
			MessageTimestamp: time.Now(),
		},
		Message: *player,
	}
}

func CreateSessionInfoMessage(session *Session) SessionInfoMessage {
	return SessionInfoMessage{
		SessionMessage: SessionMessage{
			MessageType:      "sessionInfo",
			MessageTimestamp: time.Now(),
		},
		Message: *session,
	}
}
