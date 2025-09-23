package models

import (
	"time"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

// Order represents an order in the system
type Order struct {
	ID              string      `json:"id" db:"id"`
	UserID          string      `json:"user_id" db:"user_id"`
	Status          OrderStatus `json:"status" db:"status"`
	Total           float64     `json:"total" db:"total"`
	Subtotal        float64     `json:"subtotal" db:"subtotal"`
	Tax             float64     `json:"tax" db:"tax"`
	Shipping        float64     `json:"shipping" db:"shipping"`
	ShippingAddress string      `json:"shipping_address" db:"shipping_address"`
	BillingAddress  string      `json:"billing_address" db:"billing_address"`
	PaymentIntent   *string     `json:"payment_intent" db:"payment_intent"`
	CreatedAt       time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at" db:"updated_at"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID        string  `json:"id" db:"id"`
	OrderID   string  `json:"order_id" db:"order_id"`
	ProductID string  `json:"product_id" db:"product_id"`
	Quantity  int     `json:"quantity" db:"quantity"`
	Price     float64 `json:"price" db:"price"`
}

// OrderWithItems represents an order with its items
type OrderWithItems struct {
	Order
	User      *UserResponse      `json:"user,omitempty"`
	OrderItems []OrderItemWithProduct `json:"order_items,omitempty"`
}

// OrderItemWithProduct represents an order item with product information
type OrderItemWithProduct struct {
	OrderItem
	Product *ProductWithRating `json:"product,omitempty"`
}

// OrderCreateRequest represents the request to create a new order
type OrderCreateRequest struct {
	ShippingAddress string `json:"shipping_address" binding:"required"`
	BillingAddress  string `json:"billing_address" binding:"required"`
}

// OrderUpdateRequest represents the request to update an order
type OrderUpdateRequest struct {
	Status *OrderStatus `json:"status"`
}
