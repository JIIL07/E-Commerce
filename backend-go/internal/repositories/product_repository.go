package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"ecommerce-backend/internal/models"
	"github.com/lib/pq"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *models.Product) error {
	query := `
		INSERT INTO products (id, name, slug, description, price, compare_price, images, in_stock, stock, featured, category_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	
	_, err := r.db.Exec(query, 
		product.ID, product.Name, product.Slug, product.Description, product.Price, product.ComparePrice, 
		pq.Array(product.Images), product.InStock, product.Stock, product.Featured, product.CategoryID, 
		product.CreatedAt, product.UpdatedAt,
	)
	return err
}

func (r *ProductRepository) GetByID(id string) (*models.Product, error) {
	query := `
		SELECT id, name, slug, description, price, compare_price, images, in_stock, stock, featured, category_id, created_at, updated_at
		FROM products WHERE id = $1
	`
	
	product := &models.Product{}
	var images pq.StringArray
	
	err := r.db.QueryRow(query, id).Scan(
		&product.ID, &product.Name, &product.Slug, &product.Description, &product.Price, &product.ComparePrice,
		&images, &product.InStock, &product.Stock, &product.Featured, &product.CategoryID, &product.CreatedAt, &product.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found")
	}
	
	product.Images = []string(images)
	return product, err
}

func (r *ProductRepository) ListWithFilters(query models.ProductQuery, offset int) ([]models.ProductWithCategory, int, error) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if query.Category != "" {
		whereClause += fmt.Sprintf(" AND p.category_id = $%d", argIndex)
		args = append(args, query.Category)
		argIndex++
	}

	if query.Search != "" {
		whereClause += fmt.Sprintf(" AND (p.name ILIKE $%d OR p.description ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+query.Search+"%")
		argIndex++
	}

	if query.Featured {
		whereClause += fmt.Sprintf(" AND p.featured = $%d", argIndex)
		args = append(args, true)
		argIndex++
	}

	orderClause := "ORDER BY p.created_at DESC"
	if query.SortBy != "" {
		switch query.SortBy {
		case "name":
			orderClause = "ORDER BY p.name"
		case "price":
			orderClause = "ORDER BY p.price"
		case "created_at":
			orderClause = "ORDER BY p.created_at"
		}
		
		if query.SortOrder == "asc" {
			orderClause += " ASC"
		} else {
			orderClause += " DESC"
		}
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM products p %s
	`, whereClause)
	
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	querySQL := fmt.Sprintf(`
		SELECT p.id, p.name, p.slug, p.description, p.price, p.compare_price, p.images, p.in_stock, p.stock, p.featured, p.category_id, p.created_at, p.updated_at,
		       c.id, c.name, c.slug, c.description, c.image, c.created_at, c.updated_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		%s
		%s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderClause, argIndex, argIndex+1)
	
	args = append(args, query.Limit, offset)
	
	rows, err := r.db.Query(querySQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []models.ProductWithCategory
	for rows.Next() {
		product := models.Product{}
		var category models.Category
		var images pq.StringArray
		var categoryID sql.NullString
		var categoryName sql.NullString
		var categorySlug sql.NullString
		var categoryDescription sql.NullString
		var categoryImage sql.NullString
		var categoryCreatedAt sql.NullTime
		var categoryUpdatedAt sql.NullTime

		err := rows.Scan(
			&product.ID, &product.Name, &product.Slug, &product.Description, &product.Price, &product.ComparePrice,
			&images, &product.InStock, &product.Stock, &product.Featured, &categoryID, &product.CreatedAt, &product.UpdatedAt,
			&category.ID, &categoryName, &categorySlug, &categoryDescription, &categoryImage, &categoryCreatedAt, &categoryUpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		product.Images = []string(images)
		product.CategoryID = categoryID.String

		if categoryID.Valid {
			category.ID = category.ID
			category.Name = categoryName.String
			category.Slug = categorySlug.String
			category.Description = &categoryDescription.String
			category.Image = &categoryImage.String
			category.CreatedAt = categoryCreatedAt.Time
			category.UpdatedAt = categoryUpdatedAt.Time
			products = append(products, models.ProductWithCategory{
				Product:  product,
				Category: &category,
			})
		} else {
			products = append(products, models.ProductWithCategory{
				Product:  product,
				Category: nil,
			})
		}
	}

	return products, total, nil
}

func (r *ProductRepository) GetFeatured(limit int) ([]models.ProductWithCategory, error) {
	query := `
		SELECT p.id, p.name, p.slug, p.description, p.price, p.compare_price, p.images, p.in_stock, p.stock, p.featured, p.category_id, p.created_at, p.updated_at,
		       c.id, c.name, c.slug, c.description, c.image, c.created_at, c.updated_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.featured = true AND p.in_stock = true
		ORDER BY p.created_at DESC
		LIMIT $1
	`
	
	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.ProductWithCategory
	for rows.Next() {
		product := models.Product{}
		var category models.Category
		var images pq.StringArray
		var categoryID sql.NullString
		var categoryName sql.NullString
		var categorySlug sql.NullString
		var categoryDescription sql.NullString
		var categoryImage sql.NullString
		var categoryCreatedAt sql.NullTime
		var categoryUpdatedAt sql.NullTime

		err := rows.Scan(
			&product.ID, &product.Name, &product.Slug, &product.Description, &product.Price, &product.ComparePrice,
			&images, &product.InStock, &product.Stock, &product.Featured, &categoryID, &product.CreatedAt, &product.UpdatedAt,
			&category.ID, &categoryName, &categorySlug, &categoryDescription, &categoryImage, &categoryCreatedAt, &categoryUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		product.Images = []string(images)
		product.CategoryID = categoryID.String

		if categoryID.Valid {
			category.ID = category.ID
			category.Name = categoryName.String
			category.Slug = categorySlug.String
			category.Description = &categoryDescription.String
			category.Image = &categoryImage.String
			category.CreatedAt = categoryCreatedAt.Time
			category.UpdatedAt = categoryUpdatedAt.Time
			products = append(products, models.ProductWithCategory{
				Product:  product,
				Category: &category,
			})
		} else {
			products = append(products, models.ProductWithCategory{
				Product:  product,
				Category: nil,
			})
		}
	}

	return products, nil
}

func (r *ProductRepository) Search(query string, limit int) ([]models.ProductWithCategory, error) {
	searchQuery := `
		SELECT p.id, p.name, p.slug, p.description, p.price, p.compare_price, p.images, p.in_stock, p.stock, p.featured, p.category_id, p.created_at, p.updated_at,
		       c.id, c.name, c.slug, c.description, c.image, c.created_at, c.updated_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.name ILIKE $1 OR p.description ILIKE $1
		ORDER BY p.name
		LIMIT $2
	`
	
	rows, err := r.db.Query(searchQuery, "%"+query+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.ProductWithCategory
	for rows.Next() {
		product := models.Product{}
		var category models.Category
		var images pq.StringArray
		var categoryID sql.NullString
		var categoryName sql.NullString
		var categorySlug sql.NullString
		var categoryDescription sql.NullString
		var categoryImage sql.NullString
		var categoryCreatedAt sql.NullTime
		var categoryUpdatedAt sql.NullTime

		err := rows.Scan(
			&product.ID, &product.Name, &product.Slug, &product.Description, &product.Price, &product.ComparePrice,
			&images, &product.InStock, &product.Stock, &product.Featured, &categoryID, &product.CreatedAt, &product.UpdatedAt,
			&category.ID, &categoryName, &categorySlug, &categoryDescription, &categoryImage, &categoryCreatedAt, &categoryUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		product.Images = []string(images)
		product.CategoryID = categoryID.String

		if categoryID.Valid {
			category.ID = category.ID
			category.Name = categoryName.String
			category.Slug = categorySlug.String
			category.Description = &categoryDescription.String
			category.Image = &categoryImage.String
			category.CreatedAt = categoryCreatedAt.Time
			category.UpdatedAt = categoryUpdatedAt.Time
			products = append(products, models.ProductWithCategory{
				Product:  product,
				Category: &category,
			})
		} else {
			products = append(products, models.ProductWithCategory{
				Product:  product,
				Category: nil,
			})
		}
	}

	return products, nil
}

func (r *ProductRepository) Update(id string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	setParts := make([]string, 0, len(updates))
	args := make([]interface{}, 0, len(updates)+1)
	argIndex := 1

	for key, value := range updates {
		if key == "images" {
			setParts = append(setParts, fmt.Sprintf("%s = $%d", key, argIndex))
			args = append(args, pq.Array(value))
		} else {
			setParts = append(setParts, fmt.Sprintf("%s = $%d", key, argIndex))
			args = append(args, value)
		}
		argIndex++
	}

	query := fmt.Sprintf("UPDATE products SET %s WHERE id = $%d", strings.Join(setParts, ", "), argIndex)
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *ProductRepository) Delete(id string) error {
	query := "DELETE FROM products WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}
