package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"ecommerce-backend/internal/models"
)

type ReviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(review *models.Review) error {
	query := `
		INSERT INTO reviews (id, user_id, product_id, rating, comment, helpful, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	_, err := r.db.Exec(query, review.ID, review.UserID, review.ProductID, review.Rating, review.Comment, review.Helpful, review.CreatedAt, review.UpdatedAt)
	return err
}

func (r *ReviewRepository) GetByID(id string) (*models.Review, error) {
	query := `
		SELECT id, user_id, product_id, rating, comment, helpful, created_at, updated_at
		FROM reviews WHERE id = $1
	`
	
	review := &models.Review{}
	err := r.db.QueryRow(query, id).Scan(
		&review.ID, &review.UserID, &review.ProductID, &review.Rating, &review.Comment, &review.Helpful, &review.CreatedAt, &review.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("review not found")
	}
	
	return review, err
}

func (r *ReviewRepository) GetByProductID(productID string, limit, offset int) ([]models.ReviewWithUser, error) {
	query := `
		SELECT r.id, r.user_id, r.product_id, r.rating, r.comment, r.helpful, r.created_at, r.updated_at,
		       u.name, u.image
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		WHERE r.product_id = $1
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.Query(query, productID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.ReviewWithUser
	for rows.Next() {
		review := models.Review{}
		var userName sql.NullString
		var userImage sql.NullString

		err := rows.Scan(
			&review.ID, &review.UserID, &review.ProductID, &review.Rating, &review.Comment, &review.Helpful, &review.CreatedAt, &review.UpdatedAt,
			&userName, &userImage,
		)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, models.ReviewWithUser{
			Review: review,
			UserName: userName.String,
			UserImage: &userImage.String,
		})
	}

	return reviews, nil
}

func (r *ReviewRepository) GetByUserID(userID string, limit, offset int) ([]*models.Review, error) {
	query := `
		SELECT id, user_id, product_id, rating, comment, helpful, created_at, updated_at
		FROM reviews WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*models.Review
	for rows.Next() {
		review := &models.Review{}
		err := rows.Scan(
			&review.ID, &review.UserID, &review.ProductID, &review.Rating, &review.Comment, &review.Helpful, &review.CreatedAt, &review.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

func (r *ReviewRepository) GetUserReviewForProduct(userID, productID string) (*models.Review, error) {
	query := `
		SELECT id, user_id, product_id, rating, comment, helpful, created_at, updated_at
		FROM reviews WHERE user_id = $1 AND product_id = $2
	`
	
	review := &models.Review{}
	err := r.db.QueryRow(query, userID, productID).Scan(
		&review.ID, &review.UserID, &review.ProductID, &review.Rating, &review.Comment, &review.Helpful, &review.CreatedAt, &review.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	return review, err
}

func (r *ReviewRepository) GetProductRating(productID string) (float64, int, error) {
	query := `
		SELECT AVG(rating), COUNT(*)
		FROM reviews WHERE product_id = $1
	`
	
	var avgRating sql.NullFloat64
	var count int
	
	err := r.db.QueryRow(query, productID).Scan(&avgRating, &count)
	if err != nil {
		return 0, 0, err
	}
	
	if avgRating.Valid {
		return avgRating.Float64, count, nil
	}
	
	return 0, 0, nil
}

func (r *ReviewRepository) Update(id string, updates map[string]interface{}) error {
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

	query := fmt.Sprintf("UPDATE reviews SET %s WHERE id = $%d", strings.Join(setParts, ", "), argIndex)
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *ReviewRepository) Delete(id string) error {
	query := "DELETE FROM reviews WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}
