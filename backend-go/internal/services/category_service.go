package services

import (
	"database/sql"
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repositories"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
	productRepo  *repositories.ProductRepository
}

func NewCategoryService(categoryRepo *repositories.CategoryRepository, productRepo *repositories.ProductRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		productRepo:  productRepo,
	}
}
func (s *CategoryService) GetCategories(page, limit int, includeProducts bool) ([]models.CategoryWithProducts, int, error) {
	offset := (page - 1) * limit
	categories, err := s.categoryRepo.GetCategories(limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.categoryRepo.CountCategories()
	if err != nil {
		return nil, 0, err
	}
	var categoriesWithProducts []models.CategoryWithProducts
	for _, category := range categories {
		categoryWithProducts := models.CategoryWithProducts{
			Category: *category,
		}
		if includeProducts {
			products, err := s.productRepo.GetProductsByCategory(category.ID, 10, 0)
			if err == nil {
				productsWithRating := make([]models.ProductWithRating, len(products))
				for i, product := range products {
					productsWithRating[i] = models.ProductWithRating{
						Product: *product,
					}
				}
				categoryWithProducts.Products = productsWithRating
				categoryWithProducts.Count = len(products)
			}
		}
		categoriesWithProducts = append(categoriesWithProducts, categoryWithProducts)
	}
	return categoriesWithProducts, total, nil
}
func (s *CategoryService) GetCategoryBySlug(slug string, includeProducts bool) (*models.CategoryWithProducts, error) {
	category, err := s.categoryRepo.GetCategoryBySlug(slug)
	if err != nil {
		return nil, err
	}
	categoryWithProducts := &models.CategoryWithProducts{
		Category: *category,
	}
	if includeProducts {
		products, err := s.productRepo.GetProductsByCategory(category.ID, 50, 0)
		if err == nil {
			productsWithRating := make([]models.ProductWithRating, len(products))
			for i, product := range products {
				productsWithRating[i] = models.ProductWithRating{
					Product: *product,
				}
			}
			categoryWithProducts.Products = productsWithRating
			categoryWithProducts.Count = len(products)
		}
	}
	return categoryWithProducts, nil
}
func (s *CategoryService) CreateCategory(req models.CategoryCreateRequest) (*models.Category, error) {
	slug := s.generateSlug(req.Name)
	existing, err := s.categoryRepo.GetCategoryBySlug(slug)
	if err != sql.ErrNoRows && err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("category with this name already exists")
	}
	category := &models.Category{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Slug:        slug,
		Description: &req.Description,
		Image:       req.Image,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = s.categoryRepo.CreateCategory(category)
	if err != nil {
		return nil, err
	}
	return category, nil
}
func (s *CategoryService) UpdateCategory(slug string, req models.CategoryUpdateRequest) (*models.Category, error) {
	category, err := s.categoryRepo.GetCategoryBySlug(slug)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		category.Name = *req.Name
		category.Slug = s.generateSlug(*req.Name)
	}
	if req.Description != nil {
		category.Description = req.Description
	}
	if req.Image != nil {
		category.Image = req.Image
	}
	category.UpdatedAt = time.Now()
	err = s.categoryRepo.UpdateCategory(category.ID, map[string]interface{}{
		"name":        category.Name,
		"slug":        category.Slug,
		"description": category.Description,
		"image":       category.Image,
		"updated_at":  category.UpdatedAt,
	})
	if err != nil {
		return nil, err
	}
	return category, nil
}
func (s *CategoryService) DeleteCategory(slug string) error {
	category, err := s.categoryRepo.GetCategoryBySlug(slug)
	if err != nil {
		return err
	}
	products, err := s.productRepo.GetProductsByCategory(category.ID, 1, 0)
	if err != nil {
		return err
	}
	if len(products) > 0 {
		return fmt.Errorf("cannot delete category with products")
	}
	return s.categoryRepo.DeleteCategory(category.ID)
}
func (s *CategoryService) generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	var result strings.Builder
	for _, char := range slug {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			result.WriteRune(char)
		}
	}
	return result.String()
}
