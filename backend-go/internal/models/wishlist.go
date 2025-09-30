package models
import (
	"time"
)
type WishlistItem struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	ProductID string    `json:"product_id" db:"product_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
type WishlistItemWithProduct struct {
	WishlistItem
	Product *ProductWithRating `json:"product,omitempty"`
}
type WishlistResponse struct {
	Items []WishlistItemWithProduct `json:"items"`
	Count int                       `json:"count"`
}
type WishlistItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
}
type WishlistAddRequest struct {
	ProductID string `json:"product_id" binding:"required"`
}