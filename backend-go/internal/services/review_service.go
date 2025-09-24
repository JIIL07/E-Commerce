package services

import (
	"fmt"
	"time"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repositories"
)

type ReviewService struct {
	reviewRepo *repositories.ReviewRepository
}

func NewReviewService(reviewRepo *repositories.ReviewRepository) *ReviewService {
	return &ReviewService{reviewRepo: reviewRepo}
}

func (s *ReviewService) CreateReview(userID string, req models.ReviewCreateRequest) (*models.Review, error) {
	existingReview, err := s.reviewRepo.GetUserReviewForProduct(userID, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing review: %w", err)
	}

	if existingReview != nil {
		return nil, fmt.Errorf("review already exists for this product")
	}

	review := &models.Review{
		ID:        generateID(),
		UserID:    userID,
		ProductID: req.ProductID,
		Rating:    req.Rating,
		Comment:   &req.Comment,
		Helpful:   req.Helpful,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.reviewRepo.Create(review); err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	return review, nil
}

func (s *ReviewService) GetProductReviews(productID string, page, limit int) ([]models.ReviewWithUser, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	offset := (page - 1) * limit
	return s.reviewRepo.GetByProductID(productID, limit, offset)
}

func (s *ReviewService) GetUserReviews(userID string, page, limit int) ([]*models.Review, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	offset := (page - 1) * limit
	return s.reviewRepo.GetByUserID(userID, limit, offset)
}

func (s *ReviewService) UpdateReview(userID, reviewID string, req models.ReviewUpdateRequest) (*models.Review, error) {
	review, err := s.reviewRepo.GetByID(reviewID)
	if err != nil {
		return nil, fmt.Errorf("review not found: %w", err)
	}

	if review.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	updates := make(map[string]interface{})

	if req.Rating != nil {
		updates["rating"] = *req.Rating
	}
	if req.Comment != nil {
		updates["comment"] = *req.Comment
	}
	if req.Helpful != nil {
		updates["helpful"] = *req.Helpful
	}

	if len(updates) > 0 {
		updates["updated_at"] = time.Now()

		if err := s.reviewRepo.Update(reviewID, updates); err != nil {
			return nil, fmt.Errorf("failed to update review: %w", err)
		}
	}

	updatedReview, err := s.reviewRepo.GetByID(reviewID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated review: %w", err)
	}

	return updatedReview, nil
}

func (s *ReviewService) DeleteReview(userID, reviewID string) error {
	review, err := s.reviewRepo.GetByID(reviewID)
	if err != nil {
		return fmt.Errorf("review not found: %w", err)
	}

	if review.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	return s.reviewRepo.Delete(reviewID)
}

func (s *ReviewService) GetUserReviewForProduct(userID, productID string) (*models.Review, error) {
	return s.reviewRepo.GetUserReviewForProduct(userID, productID)
}
