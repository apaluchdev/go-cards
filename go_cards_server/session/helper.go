package session

import (
	"time"

	"example.com/go_cards_server/messages"
	"example.com/go_cards_server/user"
	"github.com/google/uuid"
)

func CreateSessionInfoMessage(s *Session) messages.SessionInfoMessage {
	users := make(map[uuid.UUID]string)
	for id, user := range s.Users {
		users[id] = user.UserName
	}

	return messages.SessionInfoMessage{
		SessionId:        s.SessionId,
		SessionStartTime: s.SessionStartTime,
		Users:          users,
		MessageInfo: messages.MessageInfo{
			MessageType:      messages.SessionInfoMessageType,
			MessageTimestamp: time.Now(),
		},
	}
}

func CreateSessionStartedMessage(s *Session, userId uuid.UUID) messages.SessionStartedMessage {
	return messages.SessionStartedMessage{
		SessionId:        s.SessionId,
		SessionStartTime: s.SessionStartTime,
		Users:          s.Users,
		UserId:         userId,
		MessageInfo: messages.MessageInfo{
			MessageType:      messages.SessionStartedMessageType,
			MessageTimestamp: time.Now(),
		},
	}
}

func CreateSessionEndedMessage(s *Session) messages.SessionEndedMessage {
	return messages.SessionEndedMessage{
		SessionId: s.SessionId,
		MessageInfo: messages.MessageInfo{
			MessageType:      messages.SessionInfoMessageType,
			MessageTimestamp: time.Now(),
		},
	}
}

func CreateUserReadyMessage(userId uuid.UUID, userReady bool) messages.UserReadyMessage {
	return messages.UserReadyMessage{
		UserId:    userId.String(),
		UserReady: userReady,
		MessageInfo: messages.MessageInfo{
			MessageType:      messages.UserReadyMessageType,
			MessageTimestamp: time.Now(),
		},
	}
}

func CreateUserJoinedMessage(user *user.User) messages.UserJoinedMessage {
	return messages.UserJoinedMessage{
		UserId:   user.UserId.String(),
		UserName: user.UserName,
		MessageInfo: messages.MessageInfo{
			MessageType:      messages.UserJoinedMessageType,
			MessageTimestamp: time.Now(),
		},
	}
}