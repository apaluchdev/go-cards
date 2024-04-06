package session

import (
	"time"

	"example.com/go_cards_server/messages"
	"example.com/go_cards_server/player"
	"github.com/google/uuid"
)

func CreateSession() *Session {
	// Create a new session id
	sessionId := uuid.New()

	// Create a new session
	session := &Session{SessionId: sessionId, SessionStartTime: time.Now()}
	session.Players = make(map[uuid.UUID]*player.Player)

	return session
}

func (s *Session) AddPlayerToSession(player *player.Player) {
	s.Players[player.PlayerId] = player

	player.SendMessage(CreateSessionStartedMessage(s, player.PlayerId))
	s.BroadcastMessage(messages.CreateMessage(s.Players[player.PlayerId], messages.PlayerJoinedMessageType)) // TODO this should be a CreateXMessage

	// Handle the game communication with this player
	go s.Communicate(player.PlayerId)
}

func (s *Session) ArePlayersReady() bool {
	for _, player := range s.Players {
		if !player.PlayerReady {
			return false
		}
	}
	return true
}
