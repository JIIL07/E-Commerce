package websocket

import (
	"log"
	"sync"
	"time"
)

type Hub struct {
	clients          map[*Client]bool
	broadcast        chan []byte
	register         chan *Client
	unregister       chan *Client
	mutex            sync.RWMutex
	startTime        time.Time
	messagesSent     int64
	messagesReceived int64
	lastActivity     time.Time
}

func NewHub() *Hub {
	return &Hub{
		clients:      make(map[*Client]bool),
		broadcast:    make(chan []byte),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		startTime:    time.Now(),
		lastActivity: time.Now(),
	}
}

func (h *Hub) Run() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()

			log.Printf("Client connected. Total clients: %d", len(h.clients))
			h.lastActivity = time.Now()

			welcomeMsg := CreateMessage(MessageTypeNotification, NotificationData{
				Title:   "Welcome",
				Message: "Connected to E-Commerce WebSocket",
				Icon:    "success",
			}, client.UserID)
			h.sendToClient(client, welcomeMsg)

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mutex.Unlock()

			log.Printf("Client disconnected. Total clients: %d", len(h.clients))
			h.lastActivity = time.Now()

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
					h.messagesSent++
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mutex.RUnlock()
			h.lastActivity = time.Now()

		case <-ticker.C:
			h.mutex.RLock()
			for client := range h.clients {
				pingMsg := CreateMessage(MessageTypePing, map[string]interface{}{
					"timestamp": time.Now().Unix(),
				}, client.UserID)
				h.sendToClient(client, pingMsg)
			}
			h.mutex.RUnlock()
		}
	}
}

func (h *Hub) Broadcast(message *Message) {
	data, err := message.ToJSON()
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	select {
	case h.broadcast <- data:
	default:
		log.Println("Broadcast channel is full, dropping message")
	}
}

func (h *Hub) BroadcastToUser(userID string, message *Message) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for client := range h.clients {
		if client.UserID == userID {
			h.sendToClient(client, message)
		}
	}
}

func (h *Hub) BroadcastToRole(role string, message *Message) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for client := range h.clients {
		if client.UserRole == role {
			h.sendToClient(client, message)
		}
	}
}

func (h *Hub) sendToClient(client *Client, message *Message) {
	data, err := message.ToJSON()
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	select {
	case client.Send <- data:
		h.messagesSent++
	default:
		close(client.Send)
		delete(h.clients, client)
	}
}

func (h *Hub) GetClientCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.clients)
}

func (h *Hub) GetConnectedUsers() []ClientInfo {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	users := make([]ClientInfo, 0, len(h.clients))
	userMap := make(map[string]bool)

	for client := range h.clients {
		if client.UserID != "" && !userMap[client.UserID] {
			users = append(users, ClientInfo{
				UserID:   client.UserID,
				UserRole: client.UserRole,
				JoinedAt: client.JoinedAt,
			})
			userMap[client.UserID] = true
		}
	}

	return users
}

func (h *Hub) GetStats() HubStats {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	connectedUsers := make([]ClientInfo, 0, len(h.clients))
	userMap := make(map[string]bool)

	for client := range h.clients {
		if client.UserID != "" && !userMap[client.UserID] {
			connectedUsers = append(connectedUsers, ClientInfo{
				UserID:   client.UserID,
				UserRole: client.UserRole,
				JoinedAt: client.JoinedAt,
			})
			userMap[client.UserID] = true
		}
	}

	return HubStats{
		TotalClients:     len(h.clients),
		ConnectedUsers:   connectedUsers,
		MessagesSent:     h.messagesSent,
		MessagesReceived: h.messagesReceived,
		Uptime:           time.Since(h.startTime),
		LastActivity:     h.lastActivity,
		Metrics: map[string]interface{}{
			"active_connections": len(h.clients),
			"unique_users":       len(connectedUsers),
		},
	}
}

func (h *Hub) SendNotification(title, message, icon, priority, category string) {
	notification := CreateNotificationMessage(title, message, icon, priority, category)
	h.Broadcast(notification)
}

func (h *Hub) SendOrderUpdate(orderID, status, message, userID string) {
	orderUpdate := CreateOrderUpdateMessage(orderID, status, message, userID)

	h.BroadcastToUser(userID, orderUpdate)
	h.BroadcastToRole("admin", orderUpdate)
}

func (h *Hub) SendProductUpdate(productID, action string, data interface{}) {
	productUpdate := CreateProductUpdateMessage(productID, action, data)
	h.Broadcast(productUpdate)
}

func (h *Hub) SendStockAlert(productID, productName string, currentStock int) {
	alert := CreateStockAlertMessage(productID, productName, currentStock)
	h.BroadcastToRole("admin", alert)
}

func (h *Hub) SendPriceAlert(productID, productName string, oldPrice, newPrice float64) {
	alert := CreatePriceAlertMessage(productID, productName, oldPrice, newPrice)
	h.Broadcast(alert)
}

func (h *Hub) SendNewProductAlert(productID, productName string) {
	alert := CreateNewProductAlertMessage(productID, productName)
	h.Broadcast(alert)
}

func (h *Hub) SendPromotionAlert(title, message, actionURL string) {
	alert := CreatePromotionAlertMessage(title, message, actionURL)
	h.Broadcast(alert)
}

func (h *Hub) SendMaintenanceAlert(message string, scheduledTime time.Time) {
	alert := CreateMaintenanceAlertMessage(message, scheduledTime)
	h.Broadcast(alert)
}

func (h *Hub) SendUserActivity(userID, activity, details string) {
	activityMsg := CreateUserActivityMessage(userID, activity, details)
	h.BroadcastToRole("admin", activityMsg)
}

func (h *Hub) SendAnalyticsUpdate(metrics map[string]interface{}) {
	analyticsMsg := CreateAnalyticsUpdateMessage(metrics)
	h.BroadcastToRole("admin", analyticsMsg)
}

func (h *Hub) SendRealTimeStats(stats map[string]interface{}) {
	statsMsg := CreateRealTimeStatsMessage(stats)
	h.BroadcastToRole("admin", statsMsg)
}
