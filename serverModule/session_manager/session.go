package session_manager

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var Sessions map[uuid.UUID]*Session

type PlayerScores struct {
	Player1Score uint16 `json:"player1Score"`
	Player2Score uint16 `json:"player2Score"`
}

type SessionInfo struct {
	SessionId              uuid.UUID            `json:"sessionId"`
	SessionStartTime       time.Time            `json:"sessionStartTime"`
	SessionLastMessageTime time.Time            `json:"sessionLastMessageTime"`
	Players                map[uuid.UUID]string `json:"players"`
}

type Player struct {
	PlayerId         uuid.UUID `json:"playerId"`
	PlayerName       string    `json:"playerName"`
	PlayerConnection *websocket.Conn
}

// This should not be converted to JSON
type Session struct {
	SessionId              uuid.UUID             `json:"sessionId"`
	PlayerScores           map[uuid.UUID]int16   `json:"playerScores"`
	SessionStartTime       time.Time             `json:"sessionStartTime"`
	SessionLastMessageTime time.Time             `json:"sessionLastMessageTime"`
	Players                map[uuid.UUID]*Player `json:"players"`
}

type PlayerUpdate struct {
	Score int16
	X     int16
	Y     int16
}

func (s *Session) GetPlayerScores() map[uuid.UUID]int16 {
	return s.PlayerScores
}

func InitSessionEngine() {
	Sessions = make(map[uuid.UUID]*Session)

	go sessionCleaner()
}

func SendInitialSessionMessage(conn *websocket.Conn, s *Session) {
	sessionJSON, err := json.Marshal(CreateStartSessionMessage(s))

	if err != nil {
		fmt.Println("Error marshalling initial session message:", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, sessionJSON)
	if err != nil {
		fmt.Println("Error writing initial session message:", err)
	}
}

// make each message type have a method that makes a byte[] of json, then this method can accept any one of those message types
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

// TODO have a type that all SESSION MESSAGE types implement
// Then have individual SESSION MESSAGE types that implement the interface

func HandleUserSession(conn *websocket.Conn, s *Session, userId uuid.UUID) {
	defer conn.Close()

	go updateClient(conn, s)

	for {
		msg, err := readClientMessage(conn, s, userId)

		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				fmt.Println("WebSocket connection closed")
			}
			break
		}
		msgStr := string(msg)
		BroadcastMessage(s.SessionId, msgStr)
		fmt.Println("Received message:", msgStr)
	}
}

func updateClient(conn *websocket.Conn, s *Session) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			if !sendMessageToClient(conn, s) {
				done <- true
			}
		}
	}
}

func readClientMessage(conn *websocket.Conn, s *Session, userId uuid.UUID) (clientMessage []byte, err error) {
	// Read message from WebSocket
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return msg, nil
	//var message *PlayerUpdate
	// err = json.Unmarshal(msg, &message)
	// if err != nil {
	// 	fmt.Println("Error parsing JSON:", err)
	// 	return nil, err
	// }

	// fmt.Printf("Received new score: %v\nx: %v\ny: %v\n", message.Score, message.X, message.Y)
	// s.PlayerScores[userId] = message.Score

}

func sendMessageToClient(conn *websocket.Conn, s *Session) bool {
	// Marshal the client's session into JSON
	sessionInfoJSON, err := json.Marshal(CreateSessionInfoMessage(s))
	if err != nil {
		fmt.Println("JSON marshal error:", err)
	}

	// Write the session JSON to the client
	err = conn.WriteMessage(websocket.TextMessage, sessionInfoJSON)

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

func sessionCleaner() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		for sessionId, session := range Sessions {
			if time.Since(session.SessionLastMessageTime) > 60*time.Second {
				fmt.Println("Cleaning session:", sessionId)
				delete(Sessions, sessionId)
			}
		}
	}
}
