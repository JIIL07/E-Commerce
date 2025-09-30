package websocket

import (
	"encoding/json"
	"fmt"
	"time"

	"ecommerce-backend/internal/utils"
)

func CreateMessage(msgType MessageType, data interface{}, userID string) *Message {
	return &Message{
		Type:      msgType,
		Data:      data,
		Timestamp: time.Now(),
		ID:        utils.GenerateUUID(),
		UserID:    userID,
	}
}

func CreateNotificationMessage(title, message, icon, priority, category string) *Message {
	return CreateMessage(MessageTypeNotification, NotificationData{
		Title:   title,
		Message: message,
		Icon:    icon,
	}, "")
}

func CreateOrderUpdateMessage(orderID, status, message, userID string) *Message {
	return CreateMessage(MessageTypeOrderUpdate, OrderUpdateData{
		OrderID: orderID,
		Status:  status,
		Message: message,
		UserID:  userID,
	}, userID)
}

func CreateProductUpdateMessage(productID, action string, data interface{}) *Message {
	return CreateMessage(MessageTypeProductUpdate, ProductUpdateData{
		ProductID: productID,
		Action:    action,
		Data:      data,
	}, "")
}

func CreateStockAlertMessage(productID, productName string, currentStock int) *Message {
	return CreateMessage(MessageTypeStockAlert, StockAlertData{
		ProductID:    productID,
		ProductName:  productName,
		CurrentStock: currentStock,
	}, "")
}

func CreatePriceAlertMessage(productID, productName string, oldPrice, newPrice float64) *Message {
	return CreateMessage(MessageTypePriceAlert, PriceAlertData{
		ProductID:   productID,
		ProductName: productName,
		OldPrice:    oldPrice,
		NewPrice:    newPrice,
	}, "")
}

func CreateNewProductAlertMessage(productID, productName string) *Message {
	return CreateMessage(MessageTypeNewProductAlert, NewProductAlertData{
		ProductID:   productID,
		ProductName: productName,
	}, "")
}

func CreatePromotionAlertMessage(title, message, actionURL string) *Message {
	return CreateMessage(MessageTypePromotionAlert, PromotionAlertData{
		Title:     title,
		Message:   message,
		ActionURL: actionURL,
	}, "")
}

func CreateMaintenanceAlertMessage(message string, scheduledTime time.Time) *Message {
	return CreateMessage(MessageTypeMaintenanceAlert, MaintenanceAlertData{
		Message:       message,
		ScheduledTime: scheduledTime,
	}, "")
}

func CreateUserActivityMessage(userID, activity, details string) *Message {
	return CreateMessage(MessageTypeUserActivity, UserActivityData{
		UserID:   userID,
		Activity: activity,
		Details:  details,
	}, userID)
}

func CreateAnalyticsUpdateMessage(metrics map[string]interface{}) *Message {
	return CreateMessage(MessageTypeAnalyticsUpdate, AnalyticsUpdateData{
		Metrics: metrics,
	}, "")
}

func CreateRealTimeStatsMessage(stats map[string]interface{}) *Message {
	return CreateMessage(MessageTypeRealTimeStats, RealTimeStatsData{
		Stats: stats,
	}, "")
}

func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) FromJSON(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *Message) IsValid() bool {
	return m.Type != "" && m.Data != nil
}

func (m *Message) GetAge() time.Duration {
	return time.Since(m.Timestamp)
}

func (m *Message) IsExpired(maxAge time.Duration) bool {
	return m.GetAge() > maxAge
}

func (m *Message) String() string {
	return fmt.Sprintf("Message{Type: %s, ID: %s, UserID: %s, Timestamp: %s}", 
		m.Type, m.ID, m.UserID, m.Timestamp.Format(time.RFC3339))
}

func ValidateMessageType(msgType MessageType) bool {
	validTypes := []MessageType{
		MessageTypeNotification,
		MessageTypeOrderUpdate,
		MessageTypeProductUpdate,
		MessageTypeStockAlert,
		MessageTypePriceAlert,
		MessageTypeNewProductAlert,
		MessageTypePromotionAlert,
		MessageTypeMaintenanceAlert,
		MessageTypeUserActivity,
		MessageTypeAnalyticsUpdate,
		MessageTypeRealTimeStats,
		MessageTypePing,
		MessageTypePong,
	}
	
	for _, validType := range validTypes {
		if msgType == validType {
			return true
		}
	}
	return false
}

func GetMessageTypeFromString(typeStr string) MessageType {
	return MessageType(typeStr)
}

func IsSystemMessage(msgType MessageType) bool {
	return msgType == MessageTypePing || msgType == MessageTypePong
}

func IsUserMessage(msgType MessageType) bool {
	return !IsSystemMessage(msgType)
}

func GetMessagePriority(msgType MessageType) string {
	switch msgType {
	case MessageTypeMaintenanceAlert, MessageTypeStockAlert:
		return "high"
	case MessageTypeOrderUpdate, MessageTypePriceAlert:
		return "medium"
	default:
		return "low"
	}
}

func GetMessageCategory(msgType MessageType) string {
	switch msgType {
	case MessageTypeOrderUpdate:
		return "orders"
	case MessageTypeProductUpdate, MessageTypeStockAlert, MessageTypePriceAlert, MessageTypeNewProductAlert:
		return "products"
	case MessageTypePromotionAlert:
		return "promotions"
	case MessageTypeMaintenanceAlert:
		return "system"
	case MessageTypeUserActivity:
		return "user"
	case MessageTypeAnalyticsUpdate, MessageTypeRealTimeStats:
		return "analytics"
	default:
		return "general"
	}
}
