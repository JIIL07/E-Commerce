package handlers
import (
	"net/http"
	"strconv"
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/services"
	"github.com/gin-gonic/gin"
)
type CategoryHandler struct {
	categoryService *services.CategoryService
}
func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")
	includeProducts := c.DefaultQuery("include_products", "false")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}
	includeProductsBool := includeProducts == "true"
	categories, total, err := h.categoryService.GetCategories(page, limit, includeProductsBool)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get categories"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":    "Categories retrieved successfully",
		"categories": categories,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	})
}
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category slug is required"})
		return
	}
	includeProducts := c.DefaultQuery("include_products", "true")
	includeProductsBool := includeProducts == "true"
	category, err := h.categoryService.GetCategoryBySlug(slug, includeProductsBool)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "Category retrieved successfully",
		"category": category,
	})
}
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req models.CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category, err := h.categoryService.CreateCategory(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Category created successfully",
		"category": category,
	})
}
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category slug is required"})
		return
	}
	var req models.CategoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category, err := h.categoryService.UpdateCategory(slug, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "Category updated successfully",
		"category": category,
	})
}
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category slug is required"})
		return
	}
	err := h.categoryService.DeleteCategory(slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Category deleted successfully",
	})
}