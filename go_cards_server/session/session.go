package session

import (
	"log"
	"time"

	"example.com/go_cards_server/messages"
	"example.com/go_cards_server/user"
	"github.com/google/uuid"
)

type Session struct {
	SessionId              uuid.UUID                       `json:"sessionId"`
	SessionStartTime       time.Time                       `json:"sessionStartTime"`
	SessionLastMessageTime time.Time                       `json:"sessionLastMessageTime"`
	Users                  map[uuid.UUID]*user.User        `json:"users"`
	GameChannel            chan *messages.TypedByteMessage `json:"-"`
	MaxUsers               int                             `json:"maxUsers"`
	Active                 bool                            `json:"active"`
}

func CreateSession() *Session {
	// Create a new session id
	sessionId := uuid.New()

	// Create a new session
	session := &Session{SessionId: sessionId, SessionStartTime: time.Now(), GameChannel: make(chan *messages.TypedByteMessage), Active: true, MaxUsers: 0}
	session.Users = make(map[uuid.UUID]*user.User)

	return session
}

func (s *Session) EndSession() {
	s.Active = false

	// Ensure each user connection is closed
	for _, user := range s.Users {
		if user.UserConnection != nil {
			user.SendMessage(CreateSessionEndedMessage(s))
			user.UserConnection.Close()
		}
	}

	log.Println("Last message was at: ", s.SessionLastMessageTime)
	log.Println("Cleaning session:", s.SessionId)
}

func (s *Session) AddUserToSession(user *user.User) {
	s.Users[user.UserId] = user

	user.SendMessage(CreateSessionStartedMessage(s, user.UserId))
	s.BroadcastMessage(CreateUserJoinedMessage(s.Users[user.UserId]))

	// Handle the game communication with this user
	go s.Communicate(user.UserId)
}

func (s *Session) AreUsersReady() bool {
	if len(s.Users) < 2 {
		return false
	}

	for _, user := range s.Users {
		if !user.UserReady {
			return false
		}
	}
	return true
}
