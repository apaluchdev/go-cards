package session_manager

import (
	"encoding/json"
	"fmt"
	"time"

	"example.com/server/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

var Sessions map[uuid.UUID]*models.Session

func InitSessionEngine() {
	Sessions = make(map[uuid.UUID]*models.Session)

	// Removes sessions that have not had any messages after a predetermined time
	go sessionCleaner()
}

func HandleUserSession(conn *websocket.Conn, s *models.Session, userId uuid.UUID) {
	defer conn.Close()

	go updateClient(conn, s)

	for {
		msg, err := getClientMessage(conn, s, userId)

		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				fmt.Println("WebSocket connection closed")
			}
			break
		}

		unmashalClientMessage(s, msg, userId)

		// msgStr := string(msg)
		// BroadcastMessage(s.SessionId, msgStr)
		// fmt.Println("Received message:", msgStr)
	}
}

func unmashalClientMessage(s *models.Session, msg []byte, userId uuid.UUID) error {
	// Unmarshal message into generic PlayerMessage
	var clientMessage *models.Message
	err := json.Unmarshal(msg, &clientMessage)
	if err != nil {
		return err
	}
	messageMap := clientMessage.Message.(map[string]interface{})

	// Determine the type and unmarshal accordingly
	switch clientMessage.MessageType {
	case models.PlayerReadyMessageType:
		var playerReadyMessage models.PlayerReadyMessage
		if err := mapstructure.Decode(messageMap, &playerReadyMessage); err != nil {
			return err
		}
		playerReadyMessage.PlayerId = userId.String()
		fmt.Println("Successfully retrieved PlayerReadyType ", playerReadyMessage.PlayerId, playerReadyMessage.PlayerReady)
		BroadcastMessage(s.SessionId, models.CreateMessage(playerReadyMessage, models.PlayerReadyMessageType))

	default:
		fmt.Println("Unknown message type")
	}

	return nil
}

func updateClient(conn *websocket.Conn, s *models.Session) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			sessionInfo := models.SessionInfoMessage{
				SessionId:        s.SessionId,
				SessionStartTime: s.SessionStartTime,
				Players:          make(map[uuid.UUID]string),
			}
			if !SendMessage(conn, models.CreateMessage(sessionInfo, models.SessionInfoMessageType)) {
				done <- true
			}
		}
	}
}

func getClientMessage(conn *websocket.Conn, s *models.Session, userId uuid.UUID) (clientMessage []byte, err error) {
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
