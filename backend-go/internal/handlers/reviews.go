package handlers

import (
	"net/http"
	"strconv"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewService *services.ReviewService
}

func NewReviewHandler(reviewService *services.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID := c.GetString("user_id")
	
	var req models.ReviewCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	review, err := h.reviewService.CreateReview(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Review created successfully",
		"review":  review,
	})
}

func (h *ReviewHandler) GetProductReviews(c *gin.Context) {
	productID := c.Param("productId")
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	reviews, err := h.reviewService.GetProductReviews(productID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reviews"})
		return
	}
	
	c.JSON(http.StatusOK, reviews)
}

func (h *ReviewHandler) GetUserReviews(c *gin.Context) {
	userID := c.GetString("user_id")
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	reviews, err := h.reviewService.GetUserReviews(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reviews"})
		return
	}
	
	c.JSON(http.StatusOK, reviews)
}

func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	userID := c.GetString("user_id")
	reviewID := c.Param("id")
	
	var req models.ReviewUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	review, err := h.reviewService.UpdateReview(userID, reviewID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Review updated successfully",
		"review":  review,
	})
}

func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	userID := c.GetString("user_id")
	reviewID := c.Param("id")
	
	if err := h.reviewService.DeleteReview(userID, reviewID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}

func (h *ReviewHandler) GetUserReviewForProduct(c *gin.Context) {
	userID := c.GetString("user_id")
	productID := c.Param("productId")
	
	review, err := h.reviewService.GetUserReviewForProduct(userID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get review"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"review": review})
}
