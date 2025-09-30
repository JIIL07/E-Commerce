package handlers
import (
	"net/http"
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/services"
	"github.com/gin-gonic/gin"
)
type PaymentHandler struct {
	paymentService *services.PaymentService
}
func NewPaymentHandler(paymentService *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}
func (h *PaymentHandler) CreatePaymentIntent(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	var req models.PaymentIntentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	paymentIntent, err := h.paymentService.CreatePaymentIntent(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment intent"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":        "Payment intent created successfully",
		"payment_intent": paymentIntent,
	})
}
func (h *PaymentHandler) ConfirmPayment(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	var req models.PaymentConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payment, err := h.paymentService.ConfirmPayment(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm payment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Payment confirmed successfully",
		"payment": payment,
	})
}
func (h *PaymentHandler) GetPaymentHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	payments, err := h.paymentService.GetUserPayments(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payment history"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "Payment history retrieved successfully",
		"payments": payments,
	})
}
func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	var payload models.StripeWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook payload"})
		return
	}
	err := h.paymentService.HandleWebhook(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process webhook"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}