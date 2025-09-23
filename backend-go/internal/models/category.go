package models

import (
	"time"
)

// Category represents a product category
type Category struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description *string   `json:"description" db:"description"`
	Image       *string   `json:"image" db:"image"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CategoryWithProducts represents a category with its products
type CategoryWithProducts struct {
	Category
	Products []ProductWithRating `json:"products,omitempty"`
	Count    int                 `json:"count"`
}

// CategoryCreateRequest represents the request to create a new category
type CategoryCreateRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Image       *string `json:"image"`
}

// CategoryUpdateRequest represents the request to update a category
type CategoryUpdateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Image       *string `json:"image"`
}
