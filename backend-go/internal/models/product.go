package models
import (
	"time"
)
type Product struct {
	ID           string    `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Slug         string    `json:"slug" db:"slug"`
	Description  *string   `json:"description" db:"description"`
	Price        float64   `json:"price" db:"price"`
	ComparePrice *float64  `json:"compare_price" db:"compare_price"`
	Images       []string  `json:"images" db:"images"`
	InStock      bool      `json:"in_stock" db:"in_stock"`
	Stock        int       `json:"stock" db:"stock"`
	Featured     bool      `json:"featured" db:"featured"`
	CategoryID   string    `json:"category_id" db:"category_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
type ProductWithCategory struct {
	Product
	Category *Category `json:"category,omitempty"`
}
type ProductWithRating struct {
	Product
	Category      *Category `json:"category,omitempty"`
	AverageRating float64   `json:"average_rating"`
	ReviewCount   int       `json:"review_count"`
}
type ProductCreateRequest struct {
	Name         string   `json:"name" binding:"required"`
	Description  string   `json:"description"`
	Price        float64  `json:"price" binding:"required,min=0"`
	ComparePrice *float64 `json:"compare_price"`
	Images       []string `json:"images"`
	Stock        int      `json:"stock" binding:"required,min=0"`
	Featured     bool     `json:"featured"`
	CategoryID   string   `json:"category_id" binding:"required"`
}
type ProductUpdateRequest struct {
	Name         *string  `json:"name"`
	Description  *string  `json:"description"`
	Price        *float64 `json:"price"`
	ComparePrice *float64 `json:"compare_price"`
	Images       []string `json:"images"`
	Stock        *int     `json:"stock"`
	Featured     *bool    `json:"featured"`
	CategoryID   *string  `json:"category_id"`
}
type ProductQuery struct {
	Page      int    `form:"page"`
	Limit     int    `form:"limit"`
	Category  string `form:"category"`
	Search    string `form:"search"`
	Featured  bool   `form:"featured"`
	SortBy    string `form:"sort_by"`
	SortOrder string `form:"sort_order"`
}
type PaginatedProducts struct {
	Data       []ProductWithRating `json:"data"`
	Pagination Pagination          `json:"pagination"`
}
type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	Pages int `json:"pages"`
}