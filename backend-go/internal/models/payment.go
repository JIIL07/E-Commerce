package models
import (
	"time"
)
type PaymentStatus string
const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSucceeded PaymentStatus = "succeeded"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)
type Payment struct {
	ID              string        `json:"id" db:"id"`
	UserID          string        `json:"user_id" db:"user_id"`
	OrderID         *string       `json:"order_id" db:"order_id"`
	Amount          float64       `json:"amount" db:"amount"`
	Currency        string        `json:"currency" db:"currency"`
	Status          PaymentStatus `json:"status" db:"status"`
	PaymentIntentID string        `json:"payment_intent_id" db:"payment_intent_id"`
	ClientSecret    string        `json:"client_secret" db:"client_secret"`
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" db:"updated_at"`
}
type PaymentIntentRequest struct {
	Amount   float64 `json:"amount" binding:"required,min=0.01"`
	Currency string  `json:"currency" binding:"required"`
	OrderID  *string `json:"order_id"`
}
type PaymentConfirmRequest struct {
	PaymentIntentID string `json:"payment_intent_id" binding:"required"`
}
type PaymentIntentResponse struct {
	ID           string `json:"id"`
	ClientSecret string `json:"client_secret"`
	Amount       int64  `json:"amount"`
	Currency     string `json:"currency"`
	Status       string `json:"status"`
}
type StripeWebhookPayload struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Type    string                 `json:"type"`
	Data    map[string]interface{} `json:"data"`
	Created int64                  `json:"created"`
}