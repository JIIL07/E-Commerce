package models

import (
	"time"
)

// WishlistItem represents an item in the user's wishlist
type WishlistItem struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	ProductID string    `json:"product_id" db:"product_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// WishlistItemWithProduct represents a wishlist item with product information
type WishlistItemWithProduct struct {
	WishlistItem
	Product *ProductWithRating `json:"product,omitempty"`
}

// WishlistResponse represents the user's wishlist with items
type WishlistResponse struct {
	Items []WishlistItemWithProduct `json:"items"`
	Count int                       `json:"count"`
}

// WishlistItemRequest represents the request to add a wishlist item
type WishlistItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
}
