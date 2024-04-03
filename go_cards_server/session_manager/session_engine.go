package session_manager

import (
	"encoding/json"
	"fmt"
	"time"

	"example.com/go_cards_server/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var Sessions map[uuid.UUID]*models.Session

func InitSessionEngine() {
	Sessions = make(map[uuid.UUID]*models.Session)

	// Removes sessions that have not had any messages after a predetermined time
	go sessionCleaner()
}

func HandleMessagesFromPlayer(conn *websocket.Conn, s *models.Session, userId uuid.UUID) {
	defer conn.Close()

	//go updateClient(conn, s)

	for {
		msg, err := getClientMessage(conn)
		s.SessionLastMessageTime = time.Now()

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

		clientMessage, err := unmarshalClientMessage(msg)
		if err != nil {
			fmt.Println("Error unmarshalling client message: ", err)
			continue
		}
		handleMessage(s, clientMessage, s.Players[userId])
	}
}

func unmarshalClientMessage(msg []byte) (*models.Message, error) {
	// Unmarshal message into generic PlayerMessage
	var clientMessage *models.Message
	err := json.Unmarshal(msg, &clientMessage)
	if err != nil {
		return nil, err
	}

	return clientMessage, nil
}

func handleMessage(s *models.Session, msg *models.Message, p *models.Player) error {
	

	switch msg.MessageType {
	case models.GameMessageType:
		s.GameChannel <- msg
	case models.PlayerReadyMessageType:
		handlePlayerReadyMessage(s, msg, p)
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

func getClientMessage(conn *websocket.Conn) (clientMessage []byte, err error) {
	// Read message from WebSocket
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return msg, nil
}
