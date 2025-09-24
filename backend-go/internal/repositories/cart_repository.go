package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"ecommerce-backend/internal/models"

	"github.com/lib/pq"
)

type CartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) Create(item *models.CartItem) error {
	query := `
		INSERT INTO cart_items (id, user_id, product_id, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(query, item.ID, item.UserID, item.ProductID, item.Quantity, item.CreatedAt, item.UpdatedAt)
	return err
}

func (r *CartRepository) GetByID(id string) (*models.CartItem, error) {
	query := `
		SELECT id, user_id, product_id, quantity, created_at, updated_at
		FROM cart_items WHERE id = $1
	`

	item := &models.CartItem{}
	err := r.db.QueryRow(query, id).Scan(
		&item.ID, &item.UserID, &item.ProductID, &item.Quantity, &item.CreatedAt, &item.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("cart item not found")
	}

	return item, err
}

func (r *CartRepository) GetByUserAndProduct(userID, productID string) (*models.CartItem, error) {
	query := `
		SELECT id, user_id, product_id, quantity, created_at, updated_at
		FROM cart_items WHERE user_id = $1 AND product_id = $2
	`

	item := &models.CartItem{}
	err := r.db.QueryRow(query, userID, productID).Scan(
		&item.ID, &item.UserID, &item.ProductID, &item.Quantity, &item.CreatedAt, &item.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("cart item not found")
	}

	return item, err
}

func (r *CartRepository) GetByUserID(userID string) ([]models.CartItemWithProduct, error) {
	query := `
		SELECT ci.id, ci.user_id, ci.product_id, ci.quantity, ci.created_at, ci.updated_at,
		       p.id, p.name, p.slug, p.description, p.price, p.compare_price, p.images, p.in_stock, p.stock, p.featured, p.category_id, p.created_at, p.updated_at
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.user_id = $1
		ORDER BY ci.created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.CartItemWithProduct
	for rows.Next() {
		item := models.CartItem{}
		product := models.Product{}
		var images pq.StringArray

		err := rows.Scan(
			&item.ID, &item.UserID, &item.ProductID, &item.Quantity, &item.CreatedAt, &item.UpdatedAt,
			&product.ID, &product.Name, &product.Slug, &product.Description, &product.Price, &product.ComparePrice,
			&images, &product.InStock, &product.Stock, &product.Featured, &product.CategoryID, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		product.Images = []string(images)
		items = append(items, models.CartItemWithProduct{
			CartItem: item,
			Product:  product,
		})
	}

	return items, nil
}

func (r *CartRepository) Update(id string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	setParts := make([]string, 0, len(updates))
	args := make([]interface{}, 0, len(updates)+1)
	argIndex := 1

	for key, value := range updates {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := fmt.Sprintf("UPDATE cart_items SET %s WHERE id = $%d", strings.Join(setParts, ", "), argIndex)
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *CartRepository) Delete(id string) error {
	query := "DELETE FROM cart_items WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}

func (r *CartRepository) DeleteByUserID(userID string) error {
	query := "DELETE FROM cart_items WHERE user_id = $1"
	_, err := r.db.Exec(query, userID)
	return err
}
