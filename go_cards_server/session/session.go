package session

import (
	"log"
	"time"

	"example.com/go_cards_server/messages"
	"example.com/go_cards_server/player"
	"github.com/google/uuid"
)

type Session struct {
	SessionId              uuid.UUID                       `json:"sessionId"`
	SessionStartTime       time.Time                       `json:"sessionStartTime"`
	SessionLastMessageTime time.Time                       `json:"sessionLastMessageTime"`
	Players                map[uuid.UUID]*player.Player    `json:"players"`
	GameChannel            chan *messages.TypedByteMessage `json:"-"`
	MaxPlayers             int                             `json:"maxPlayers"`
	Active                 bool                            `json:"active"`
}

func CreateSession() *Session {
	// Create a new session id
	sessionId := uuid.New()

	// Create a new session
	session := &Session{SessionId: sessionId, SessionStartTime: time.Now(), GameChannel: make(chan *messages.TypedByteMessage), Active: true, MaxPlayers: 0}
	session.Players = make(map[uuid.UUID]*player.Player)

	return session
}

func (s *Session) EndSession() {
	s.Active = false

	// Ensure each player connection is closed
	for _, player := range s.Players {
		if player.PlayerConnection != nil {
			player.SendMessage(CreateSessionEndedMessage(s))
			player.PlayerConnection.Close()
		}
	}

	log.Println("Cleaning session:", s.SessionId)
}

func (s *Session) AddPlayerToSession(player *player.Player) {
	s.Players[player.PlayerId] = player

	player.SendMessage(CreateSessionStartedMessage(s, player.PlayerId))
	s.BroadcastMessage(CreatePlayerJoinedMessage(s.Players[player.PlayerId]))

	// Handle the game communication with this player
	go s.Communicate(player.PlayerId)
}

func (s *Session) ArePlayersReady() bool {
	if len(s.Players) < 2 {
		return false
	}

	for _, player := range s.Players {
		if !player.PlayerReady {
			return false
		}
	}
	return true
}
