package models

import (
	"time"
)

// CartItem represents an item in the user's cart
type CartItem struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	ProductID string    `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CartItemWithProduct represents a cart item with product information
type CartItemWithProduct struct {
	CartItem
	Product *ProductWithRating `json:"product,omitempty"`
}

// CartResponse represents the user's cart with items
type CartResponse struct {
	Items []CartItemWithProduct `json:"items"`
	Total float64               `json:"total"`
}

// CartItemRequest represents the request to add/update a cart item
type CartItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}
