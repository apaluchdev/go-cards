package session_manager

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// BroadcastMessage sends a message to all players in a session
func BroadcastMessage(sessionId uuid.UUID, message any) {

	messageJSON, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling broadcast message ", err)
	}

	for _, player := range Sessions[sessionId].Players {
		if player.PlayerConnection == nil {
			continue
		}

		err := player.PlayerConnection.WriteMessage(websocket.TextMessage, messageJSON)
		if err != nil {
			fmt.Println("Error marshalling initial session message:", err)
		}
	}
}

func SendMessage(conn *websocket.Conn, message any) bool {
	// Marshal the message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling message ", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, messageJSON)

	// If there was an error writing to the client, return false
	if err != nil {
		fmt.Println("Write message error:", err)
		if err == websocket.ErrCloseSent {
			fmt.Println("Breaking connection")
			return false
		}
	}

	return true
}
