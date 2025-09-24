package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"ecommerce-backend/internal/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(category *models.Category) error {
	query := `
		INSERT INTO categories (id, name, slug, description, image, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	_, err := r.db.Exec(query, category.ID, category.Name, category.Slug, category.Description, category.Image, category.CreatedAt, category.UpdatedAt)
	return err
}

func (r *CategoryRepository) GetByID(id string) (*models.Category, error) {
	query := `
		SELECT id, name, slug, description, image, created_at, updated_at
		FROM categories WHERE id = $1
	`
	
	category := &models.Category{}
	err := r.db.QueryRow(query, id).Scan(
		&category.ID, &category.Name, &category.Slug, &category.Description, &category.Image, &category.CreatedAt, &category.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("category not found")
	}
	
	return category, err
}

func (r *CategoryRepository) GetBySlug(slug string) (*models.Category, error) {
	query := `
		SELECT id, name, slug, description, image, created_at, updated_at
		FROM categories WHERE slug = $1
	`
	
	category := &models.Category{}
	err := r.db.QueryRow(query, slug).Scan(
		&category.ID, &category.Name, &category.Slug, &category.Description, &category.Image, &category.CreatedAt, &category.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("category not found")
	}
	
	return category, err
}

func (r *CategoryRepository) List(limit, offset int) ([]*models.Category, error) {
	query := `
		SELECT id, name, slug, description, image, created_at, updated_at
		FROM categories ORDER BY name LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		category := &models.Category{}
		err := rows.Scan(
			&category.ID, &category.Name, &category.Slug, &category.Description, &category.Image, &category.CreatedAt, &category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *CategoryRepository) Update(id string, updates map[string]interface{}) error {
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

	query := fmt.Sprintf("UPDATE categories SET %s WHERE id = $%d", strings.Join(setParts, ", "), argIndex)
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *CategoryRepository) Delete(id string) error {
	query := "DELETE FROM categories WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}
