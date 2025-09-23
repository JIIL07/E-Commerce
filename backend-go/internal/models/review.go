package models

import (
	"time"
)

// Review represents a product review
type Review struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	ProductID string    `json:"product_id" db:"product_id"`
	Rating    int       `json:"rating" db:"rating"`
	Comment   *string   `json:"comment" db:"comment"`
	Helpful   *bool     `json:"helpful" db:"helpful"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ReviewWithUser represents a review with user information
type ReviewWithUser struct {
	Review
	User *UserResponse `json:"user,omitempty"`
}

// ReviewCreateRequest represents the request to create a new review
type ReviewCreateRequest struct {
	ProductID string  `json:"product_id" binding:"required"`
	Rating    int     `json:"rating" binding:"required,min=1,max=5"`
	Comment   string  `json:"comment" binding:"required"`
	Helpful   *bool   `json:"helpful"`
}

// ReviewUpdateRequest represents the request to update a review
type ReviewUpdateRequest struct {
	Rating  *int    `json:"rating"`
	Comment *string `json:"comment"`
	Helpful *bool   `json:"helpful"`
}
