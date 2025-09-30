package models
import (
	"time"
)
type Category struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description *string   `json:"description" db:"description"`
	Image       *string   `json:"image" db:"image"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
type CategoryWithProducts struct {
	Category
	Products []ProductWithRating `json:"products,omitempty"`
	Count    int                 `json:"count"`
}
type CategoryCreateRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Image       *string `json:"image"`
}
type CategoryUpdateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Image       *string `json:"image"`
}