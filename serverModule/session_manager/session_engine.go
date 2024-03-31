package session_manager

import (
	"encoding/json"
	"fmt"
	"time"

	"example.com/server/models"
	"example.com/server/models/message_models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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

		if Sessions[s.SessionId] == nil {
			fmt.Println("Session does not exist")
			break
		}

		s.SessionLastMessageTime = time.Now()
		unmarshalClientMessage(s, msg, userId)
	}
}

func unmarshalClientMessage(s *models.Session, msg []byte, userId uuid.UUID) error {
	// Unmarshal message into generic PlayerMessage
	var clientMessage *message_models.Message
	err := json.Unmarshal(msg, &clientMessage)
	if err != nil {
		return err
	}

	handleMessage(s, clientMessage, s.Players[userId])

	return nil
}

func handleMessage(s *models.Session, clientMessage *message_models.Message, p *models.Player) error {
	messageMap := clientMessage.Message.(map[string]interface{})

	switch clientMessage.MessageType {
	case message_models.PlayerReadyMessageType:
		handlePlayerReadyMessage(s, messageMap, p)
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
			sessionInfo := message_models.SessionInfoMessage{
				SessionId:        s.SessionId,
				SessionStartTime: s.SessionStartTime,
				Players:          make(map[uuid.UUID]string),
			}
			if !SendMessage(conn, message_models.CreateMessage(sessionInfo, message_models.SessionInfoMessageType)) {
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
}
