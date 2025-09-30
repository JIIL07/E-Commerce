package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var message Message
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		c.handleMessage(&message)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(message *Message) {
	switch message.Type {
	case MessageTypePing:
		c.handlePing()
	case MessageTypeNotification:
		c.handleChatMessage(message)
	case MessageTypeUserActivity:
		c.handleJoinRoom(message)
	default:
		log.Printf("Unknown message type: %s", message.Type)
	}
}

func (c *Client) handlePing() {
	pongMsg := CreateMessage(MessageTypePong, map[string]interface{}{
		"timestamp": time.Now().Unix(),
	}, c.UserID)

	data, _ := pongMsg.ToJSON()
	select {
	case c.Send <- data:
	default:
		close(c.Send)
	}
}

func (c *Client) handleChatMessage(message *Message) {
	if c.UserID == "" {
		return
	}

	chatData, ok := message.Data.(map[string]interface{})
	if !ok {
		return
	}

	chatMessage, _ := chatData["message"].(string)
	roomID, _ := chatData["room_id"].(string)

	if roomID == "" {
		return
	}

	if chatMessage == "" {
		return
	}

	chatMsg := CreateMessage(MessageTypeNotification, NotificationData{
		Title:   "Chat Message",
		Message: chatMessage,
		Icon:    "chat",
	}, c.UserID)

	c.Hub.Broadcast(chatMsg)
}

func (c *Client) handleJoinRoom(message *Message) {
	roomData, ok := message.Data.(map[string]interface{})
	if !ok {
		return
	}

	roomID, _ := roomData["room_id"].(string)
	if roomID == "" {
		return
	}

	joinMsg := CreateMessage(MessageTypeUserActivity, UserActivityData{
		UserID:   c.UserID,
		Activity: "joined_room",
		Details:  roomID,
	}, c.UserID)

	c.Hub.Broadcast(joinMsg)
}

func (c *Client) ReadPump() {
	c.readPump()
}

func (c *Client) WritePump() {
	c.writePump()
}
