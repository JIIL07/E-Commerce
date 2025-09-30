package websocket

import (
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   string
	UserRole string
	JoinedAt time.Time
}

type MessageType string

const (
	MessageTypeNotification     MessageType = "notification"
	MessageTypeOrderUpdate      MessageType = "order_update"
	MessageTypeProductUpdate    MessageType = "product_update"
	MessageTypeStockAlert       MessageType = "stock_alert"
	MessageTypePriceAlert       MessageType = "price_alert"
	MessageTypeNewProductAlert  MessageType = "new_product_alert"
	MessageTypePromotionAlert   MessageType = "promotion_alert"
	MessageTypeMaintenanceAlert MessageType = "maintenance_alert"
	MessageTypeUserActivity     MessageType = "user_activity"
	MessageTypeAnalyticsUpdate  MessageType = "analytics_update"
	MessageTypeRealTimeStats    MessageType = "real_time_stats"
	MessageTypePing             MessageType = "ping"
	MessageTypePong             MessageType = "pong"
)

type Message struct {
	Type      MessageType `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	ID        string      `json:"id,omitempty"`
	UserID    string      `json:"user_id,omitempty"`
	Priority  string      `json:"priority,omitempty"`
	Category  string      `json:"category,omitempty"`
}

type NotificationData struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Icon    string `json:"icon,omitempty"`
}

type OrderUpdateData struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

type ProductUpdateData struct {
	ProductID string      `json:"product_id"`
	Action    string      `json:"action"`
	Data      interface{} `json:"data"`
}

type StockAlertData struct {
	ProductID    string `json:"product_id"`
	ProductName  string `json:"product_name"`
	CurrentStock int    `json:"current_stock"`
}

type PriceAlertData struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	OldPrice    float64 `json:"old_price"`
	NewPrice    float64 `json:"new_price"`
}

type NewProductAlertData struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
}

type PromotionAlertData struct {
	Title     string `json:"title"`
	Message   string `json:"message"`
	ActionURL string `json:"action_url,omitempty"`
}

type MaintenanceAlertData struct {
	Message       string    `json:"message"`
	ScheduledTime time.Time `json:"scheduled_time"`
}

type UserActivityData struct {
	UserID   string `json:"user_id"`
	Activity string `json:"activity"`
	Details  string `json:"details,omitempty"`
}

type AnalyticsUpdateData struct {
	Metrics map[string]interface{} `json:"metrics"`
}

type RealTimeStatsData struct {
	Stats map[string]interface{} `json:"stats"`
}

type ClientInfo struct {
	UserID   string    `json:"user_id"`
	UserRole string    `json:"user_role"`
	JoinedAt time.Time `json:"joined_at"`
}

type HubStats struct {
	TotalClients     int                    `json:"total_clients"`
	ConnectedUsers   []ClientInfo           `json:"connected_users"`
	MessagesSent     int64                  `json:"messages_sent"`
	MessagesReceived int64                  `json:"messages_received"`
	Uptime           time.Duration          `json:"uptime"`
	LastActivity     time.Time              `json:"last_activity"`
	Metrics          map[string]interface{} `json:"metrics"`
}
