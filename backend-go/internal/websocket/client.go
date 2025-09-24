package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	UserID  string      `json:"user_id,omitempty"`
	Role    string      `json:"role,omitempty"`
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func HandleWebSocket(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}

		userID := c.GetString("user_id")
		userRole := c.GetString("user_role")

		client := &Client{
			hub:      hub,
			conn:     conn,
			send:     make(chan []byte, 256),
			userID:   userID,
			userRole: userRole,
		}

		client.hub.register <- client

		go client.writePump()
		go client.readPump()

		welcomeMessage := Message{
			Type: "welcome",
			Data: map[string]interface{}{
				"message": "Connected to E-Commerce WebSocket",
				"user_id": userID,
			},
		}

		messageBytes, _ := json.Marshal(welcomeMessage)
		client.send <- messageBytes
	}
}

func BroadcastOrderUpdate(hub *Hub, userID string, orderData interface{}) {
	message := Message{
		Type:   "order_update",
		Data:   orderData,
		UserID: userID,
	}

	messageBytes, _ := json.Marshal(message)
	hub.BroadcastToUser(userID, messageBytes)
}

func BroadcastProductUpdate(hub *Hub, productData interface{}) {
	message := Message{
		Type: "product_update",
		Data: productData,
	}

	messageBytes, _ := json.Marshal(message)
	hub.broadcast <- messageBytes
}

func BroadcastAdminNotification(hub *Hub, notification interface{}) {
	message := Message{
		Type: "admin_notification",
		Data: notification,
		Role: "admin",
	}

	messageBytes, _ := json.Marshal(message)
	hub.BroadcastToRole("admin", messageBytes)
}
