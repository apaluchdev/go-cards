package session

import (
	"fmt"
	"time"

	"example.com/go_cards_server/messages"
	"github.com/google/uuid"
)

// BroadcastMessage sends a message to all players in a session
func (s *Session) BroadcastMessage(message messages.Message) {
	for _, player := range s.Players {
		err := player.SendMessage(message)
		if err != nil {
			fmt.Println("Error marshalling initial session message:", err)
		}
	}
}

func (s *Session) Communicate(playerId uuid.UUID) {
	playerConnection := s.Players[playerId].PlayerConnection

	defer playerConnection.Close()

	for {
		_, msg, err := playerConnection.ReadMessage()
		fmt.Println("Received message: ", string(msg))
		if err != nil {
			fmt.Println("Error reading message: ", err)
			break
		}

		s.SessionLastMessageTime = time.Now()

		playerMessage, err := messages.UnmarshalByteMessage(msg)
		if err != nil {
			fmt.Println("Error unmarshalling client message: ", err)
			continue
		}

		// Allow the session to do any processing of the message first
		s.handleMessage(playerMessage, s.Players[playerId])

		// Send the message off to the game channel to be handled by whichever game is being played
		s.GameChannel <- playerMessage
	}
}
