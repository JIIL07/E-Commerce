package handlers
import (
	"net/http"
	"strconv"
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/services"
	"github.com/gin-gonic/gin"
)
type WishlistHandler struct {
	wishlistService *services.WishlistService
}
func NewWishlistHandler(wishlistService *services.WishlistService) *WishlistHandler {
	return &WishlistHandler{wishlistService: wishlistService}
}
func (h *WishlistHandler) GetWishlist(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}
	wishlist, total, err := h.wishlistService.GetUserWishlist(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get wishlist"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Wishlist retrieved successfully",
		"wishlist": wishlist,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	})
}
func (h *WishlistHandler) AddToWishlist(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	var req models.WishlistAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.wishlistService.AddToWishlist(userID, req.ProductID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to wishlist"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Product added to wishlist successfully",
	})
}
func (h *WishlistHandler) RemoveFromWishlist(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	productID := c.Param("productId")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product ID is required"})
		return
	}
	err := h.wishlistService.RemoveFromWishlist(userID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from wishlist"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Product removed from wishlist successfully",
	})
}
func (h *WishlistHandler) IsInWishlist(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	productID := c.Param("productId")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product ID is required"})
		return
	}
	isInWishlist, err := h.wishlistService.IsInWishlist(userID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check wishlist status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":        "Wishlist status retrieved successfully",
		"is_in_wishlist": isInWishlist,
	})
}
func (h *WishlistHandler) ClearWishlist(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	err := h.wishlistService.ClearWishlist(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear wishlist"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Wishlist cleared successfully",
	})
}