package websocket

import (
	"net/http"
	"strings"
	"time"

	"ecommerce-backend/internal/utils"

	"github.com/gin-gonic/gin"
	gorilla "github.com/gorilla/websocket"
)

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

func (h *Handler) HandleWebSocket(c *gin.Context) {
	upgrader := gorilla.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.Error("Failed to upgrade connection", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	userID := h.extractUserID(c)
	userRole := h.extractUserRole(c)

	client := &Client{
		Hub:      h.hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		UserID:   userID,
		UserRole: userRole,
		JoinedAt: time.Now(),
	}

	client.Hub.register <- client

	go client.WritePump()
	go client.ReadPump()
}

func (h *Handler) extractUserID(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if claims, err := utils.ValidateJWT(tokenString); err == nil {
			return claims.UserID
		}
	}

	userID := c.Query("user_id")
	if userID != "" {
		return userID
	}

	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}

	return ""
}

func (h *Handler) extractUserRole(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if claims, err := utils.ValidateJWT(tokenString); err == nil {
			return claims.Role
		}
	}

	userRole := c.Query("user_role")
	if userRole != "" {
		return userRole
	}

	if userRole, exists := c.Get("user_role"); exists {
		if role, ok := userRole.(string); ok {
			return role
		}
	}

	return "guest"
}

func (h *Handler) GetConnectedUsers(c *gin.Context) {
	users := h.hub.GetConnectedUsers()
	c.JSON(http.StatusOK, gin.H{
		"connected_users": users,
		"total_count":     len(users),
	})
}

func (h *Handler) GetClientCount(c *gin.Context) {
	count := h.hub.GetClientCount()
	c.JSON(http.StatusOK, gin.H{
		"client_count": count,
	})
}

func (h *Handler) SendNotification(c *gin.Context) {
	var req struct {
		Title    string `json:"title" binding:"required"`
		Message  string `json:"message" binding:"required"`
		Icon     string `json:"icon,omitempty"`
		Priority string `json:"priority,omitempty"`
		Category string `json:"category,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.SendNotification(req.Title, req.Message, req.Icon, req.Priority, req.Category)
	c.JSON(http.StatusOK, gin.H{"message": "Notification sent"})
}

func (h *Handler) SendOrderUpdate(c *gin.Context) {
	var req struct {
		OrderID string `json:"order_id" binding:"required"`
		Status  string `json:"status" binding:"required"`
		Message string `json:"message" binding:"required"`
		UserID  string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.SendOrderUpdate(req.OrderID, req.Status, req.Message, req.UserID)
	c.JSON(http.StatusOK, gin.H{"message": "Order update sent"})
}

func (h *Handler) SendProductUpdate(c *gin.Context) {
	var req struct {
		ProductID string      `json:"product_id" binding:"required"`
		Action    string      `json:"action" binding:"required"`
		Data      interface{} `json:"data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.SendProductUpdate(req.ProductID, req.Action, req.Data)
	c.JSON(http.StatusOK, gin.H{"message": "Product update sent"})
}

func (h *Handler) SendStockAlert(c *gin.Context) {
	var req struct {
		ProductID    string `json:"product_id" binding:"required"`
		ProductName  string `json:"product_name" binding:"required"`
		CurrentStock int    `json:"current_stock" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.SendStockAlert(req.ProductID, req.ProductName, req.CurrentStock)
	c.JSON(http.StatusOK, gin.H{"message": "Stock alert sent"})
}

func (h *Handler) SendPriceAlert(c *gin.Context) {
	var req struct {
		ProductID   string  `json:"product_id" binding:"required"`
		ProductName string  `json:"product_name" binding:"required"`
		OldPrice    float64 `json:"old_price" binding:"required"`
		NewPrice    float64 `json:"new_price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.SendPriceAlert(req.ProductID, req.ProductName, req.OldPrice, req.NewPrice)
	c.JSON(http.StatusOK, gin.H{"message": "Price alert sent"})
}

func (h *Handler) SendNewProductAlert(c *gin.Context) {
	var req struct {
		ProductID   string `json:"product_id" binding:"required"`
		ProductName string `json:"product_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.SendNewProductAlert(req.ProductID, req.ProductName)
	c.JSON(http.StatusOK, gin.H{"message": "New product alert sent"})
}

func (h *Handler) SendPromotionAlert(c *gin.Context) {
	var req struct {
		Title     string `json:"title" binding:"required"`
		Message   string `json:"message" binding:"required"`
		ActionURL string `json:"action_url,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.SendPromotionAlert(req.Title, req.Message, req.ActionURL)
	c.JSON(http.StatusOK, gin.H{"message": "Promotion alert sent"})
}

func (h *Handler) SendMaintenanceAlert(c *gin.Context) {
	var req struct {
		Message       string `json:"message" binding:"required"`
		ScheduledTime string `json:"scheduled_time" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	scheduledTime, err := time.Parse("2006-01-02 15:04", req.ScheduledTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scheduled_time format"})
		return
	}

	h.hub.SendMaintenanceAlert(req.Message, scheduledTime)
	c.JSON(http.StatusOK, gin.H{"message": "Maintenance alert sent"})
}

func (h *Handler) SendUserActivity(c *gin.Context) {
	var req struct {
		UserID   string `json:"user_id" binding:"required"`
		Activity string `json:"activity" binding:"required"`
		Details  string `json:"details,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.SendUserActivity(req.UserID, req.Activity, req.Details)
	c.JSON(http.StatusOK, gin.H{"message": "User activity sent"})
}

func (h *Handler) SendAnalyticsUpdate(c *gin.Context) {
	var req struct {
		Metrics map[string]interface{} `json:"metrics" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.SendAnalyticsUpdate(req.Metrics)
	c.JSON(http.StatusOK, gin.H{"message": "Analytics update sent"})
}

func (h *Handler) SendRealTimeStats(c *gin.Context) {
	var req struct {
		Stats map[string]interface{} `json:"stats" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.SendRealTimeStats(req.Stats)
	c.JSON(http.StatusOK, gin.H{"message": "Real-time stats sent"})
}
