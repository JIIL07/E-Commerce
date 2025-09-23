package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"ecommerce-backend/internal/database"
	"ecommerce-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	var query models.ProductQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
		return
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 10
	}
	if query.SortBy == "" {
		query.SortBy = "created_at"
	}
	if query.SortOrder == "" {
		query.SortOrder = "desc"
	}

	whereConditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if query.Category != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("c.slug = $%d", argIndex))
		args = append(args, query.Category)
		argIndex++
	}

	if query.Search != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("(p.name ILIKE $%d OR p.description ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+query.Search+"%")
		argIndex++
	}

	if query.Featured {
		whereConditions = append(whereConditions, fmt.Sprintf("p.featured = $%d", argIndex))
		args = append(args, true)
		argIndex++
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	orderBy := fmt.Sprintf("ORDER BY p.%s %s", query.SortBy, strings.ToUpper(query.SortOrder))

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		%s
	`, whereClause)

	var total int
	err := database.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to count products",
		})
		return
	}

	offset := (query.Page - 1) * query.Limit
	pages := (total + query.Limit - 1) / query.Limit

	productsQuery := fmt.Sprintf(`
		SELECT 
			p.id, p.name, p.slug, p.description, p.price, p.compare_price, p.images,
			p.in_stock, p.stock, p.featured, p.category_id, p.created_at, p.updated_at,
			c.id, c.name, c.slug, c.description, c.image,
			COALESCE(AVG(r.rating), 0) as average_rating,
			COUNT(r.id) as review_count
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN reviews r ON p.id = r.product_id
		%s
		GROUP BY p.id, c.id
		%s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderBy, argIndex, argIndex+1)

	args = append(args, query.Limit, offset)

	rows, err := database.DB.Query(productsQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to fetch products",
		})
		return
	}
	defer rows.Close()

	var products []models.ProductWithRating
	for rows.Next() {
		var product models.ProductWithRating
		var category models.Category
		var imagesStr string

		err := rows.Scan(
			&product.ID, &product.Name, &product.Slug, &product.Description, &product.Price,
			&product.ComparePrice, &imagesStr, &product.InStock, &product.Stock,
			&product.Featured, &product.CategoryID, &product.CreatedAt, &product.UpdatedAt,
			&category.ID, &category.Name, &category.Slug, &category.Description, &category.Image,
			&product.AverageRating, &product.ReviewCount,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Message: "Failed to scan product",
			})
			return
		}

		if imagesStr != "" {
			product.Images = strings.Split(strings.Trim(imagesStr, "{}"), ",")
		}

		product.Category = &category
		products = append(products, product)
	}

	c.JSON(http.StatusOK, models.PaginatedProducts{
		Data: products,
		Pagination: models.Pagination{
			Page:  query.Page,
			Limit: query.Limit,
			Total: total,
			Pages: pages,
		},
	})
}

func GetProduct(c *gin.Context) {
	identifier := c.Param("id")

	var query string
	var args []interface{}

	if len(identifier) == 36 {
		query = `
			SELECT 
				p.id, p.name, p.slug, p.description, p.price, p.compare_price, p.images,
				p.in_stock, p.stock, p.featured, p.category_id, p.created_at, p.updated_at,
				c.id, c.name, c.slug, c.description, c.image,
				COALESCE(AVG(r.rating), 0) as average_rating,
				COUNT(r.id) as review_count
			FROM products p
			LEFT JOIN categories c ON p.category_id = c.id
			LEFT JOIN reviews r ON p.id = r.product_id
			WHERE p.id = $1
			GROUP BY p.id, c.id
		`
		args = []interface{}{identifier}
	} else {
		query = `
			SELECT 
				p.id, p.name, p.slug, p.description, p.price, p.compare_price, p.images,
				p.in_stock, p.stock, p.featured, p.category_id, p.created_at, p.updated_at,
				c.id, c.name, c.slug, c.description, c.image,
				COALESCE(AVG(r.rating), 0) as average_rating,
				COUNT(r.id) as review_count
			FROM products p
			LEFT JOIN categories c ON p.category_id = c.id
			LEFT JOIN reviews r ON p.id = r.product_id
			WHERE p.slug = $1
			GROUP BY p.id, c.id
		`
		args = []interface{}{identifier}
	}

	var product models.ProductWithRating
	var category models.Category
	var imagesStr string

	err := database.DB.QueryRow(query, args...).Scan(
		&product.ID, &product.Name, &product.Slug, &product.Description, &product.Price,
		&product.ComparePrice, &imagesStr, &product.InStock, &product.Stock,
		&product.Featured, &product.CategoryID, &product.CreatedAt, &product.UpdatedAt,
		&category.ID, &category.Name, &category.Slug, &category.Description, &category.Image,
		&product.AverageRating, &product.ReviewCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Message: "Product not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to fetch product",
		})
		return
	}

	if imagesStr != "" {
		product.Images = strings.Split(strings.Trim(imagesStr, "{}"), ",")
	}

	product.Category = &category

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Product retrieved successfully",
		Data:    product,
	})
}

func GetFeaturedProducts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "8")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 8
	}

	query := `
		SELECT 
			p.id, p.name, p.slug, p.description, p.price, p.compare_price, p.images,
			p.in_stock, p.stock, p.featured, p.category_id, p.created_at, p.updated_at,
			c.id, c.name, c.slug, c.description, c.image,
			COALESCE(AVG(r.rating), 0) as average_rating,
			COUNT(r.id) as review_count
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN reviews r ON p.id = r.product_id
		WHERE p.featured = true
		GROUP BY p.id, c.id
		ORDER BY p.created_at DESC
		LIMIT $1
	`

	rows, err := database.DB.Query(query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to fetch featured products",
		})
		return
	}
	defer rows.Close()

	var products []models.ProductWithRating
	for rows.Next() {
		var product models.ProductWithRating
		var category models.Category
		var imagesStr string

		err := rows.Scan(
			&product.ID, &product.Name, &product.Slug, &product.Description, &product.Price,
			&product.ComparePrice, &imagesStr, &product.InStock, &product.Stock,
			&product.Featured, &product.CategoryID, &product.CreatedAt, &product.UpdatedAt,
			&category.ID, &category.Name, &category.Slug, &category.Description, &category.Image,
			&product.AverageRating, &product.ReviewCount,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Message: "Failed to scan product",
			})
			return
		}

		if imagesStr != "" {
			product.Images = strings.Split(strings.Trim(imagesStr, "{}"), ",")
		}

		product.Category = &category
		products = append(products, product)
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Featured products retrieved successfully",
		Data:    products,
	})
}