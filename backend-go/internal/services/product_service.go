package services

import (
	"fmt"
	"math"
	"strings"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repositories"
)

type ProductService struct {
	productRepo *repositories.ProductRepository
	categoryRepo *repositories.CategoryRepository
	reviewRepo   *repositories.ReviewRepository
}

func NewProductService(productRepo *repositories.ProductRepository, categoryRepo *repositories.CategoryRepository, reviewRepo *repositories.ReviewRepository) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		reviewRepo:   reviewRepo,
	}
}

func (s *ProductService) CreateProduct(req models.ProductCreateRequest) (*models.ProductWithCategory, error) {
	product := &models.Product{
		ID:          generateID(),
		Name:        req.Name,
		Slug:        generateSlug(req.Name),
		Description: &req.Description,
		Price:       req.Price,
		ComparePrice: req.ComparePrice,
		Images:      req.Images,
		InStock:     req.Stock > 0,
		Stock:       req.Stock,
		Featured:    req.Featured,
		CategoryID:  req.CategoryID,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return s.GetProductWithCategory(product.ID)
}

func (s *ProductService) GetProduct(id string) (*models.ProductWithCategory, error) {
	return s.GetProductWithCategory(id)
}

func (s *ProductService) GetProductWithCategory(id string) (*models.ProductWithCategory, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	var category *models.Category
	if product.CategoryID != "" {
		category, _ = s.categoryRepo.GetByID(product.CategoryID)
	}

	return &models.ProductWithCategory{
		Product:  *product,
		Category: category,
	}, nil
}

func (s *ProductService) GetProducts(query models.ProductQuery) (*models.PaginatedProducts, error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 20
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	offset := (query.Page - 1) * query.Limit

	products, total, err := s.productRepo.ListWithFilters(query, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	productsWithRating := make([]models.ProductWithRating, len(products))
	for i, product := range products {
		rating, reviewCount := s.getProductRating(product.ID)
		productsWithRating[i] = models.ProductWithRating{
			Product:       product.Product,
			Category:      product.Category,
			AverageRating: rating,
			ReviewCount:   reviewCount,
		}
	}

	pages := int(math.Ceil(float64(total) / float64(query.Limit)))

	return &models.PaginatedProducts{
		Data: productsWithRating,
		Pagination: models.Pagination{
			Page:  query.Page,
			Limit: query.Limit,
			Total: total,
			Pages: pages,
		},
	}, nil
}

func (s *ProductService) GetFeaturedProducts(limit int) ([]models.ProductWithRating, error) {
	if limit <= 0 {
		limit = 10
	}

	products, err := s.productRepo.GetFeatured(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get featured products: %w", err)
	}

	productsWithRating := make([]models.ProductWithRating, len(products))
	for i, product := range products {
		rating, reviewCount := s.getProductRating(product.ID)
		productsWithRating[i] = models.ProductWithRating{
			Product:       product.Product,
			Category:      product.Category,
			AverageRating: rating,
			ReviewCount:   reviewCount,
		}
	}

	return productsWithRating, nil
}

func (s *ProductService) UpdateProduct(id string, req models.ProductUpdateRequest) (*models.ProductWithCategory, error) {
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
		updates["slug"] = generateSlug(*req.Name)
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.ComparePrice != nil {
		updates["compare_price"] = *req.ComparePrice
	}
	if req.Images != nil {
		updates["images"] = req.Images
	}
	if req.Stock != nil {
		updates["stock"] = *req.Stock
		updates["in_stock"] = *req.Stock > 0
	}
	if req.Featured != nil {
		updates["featured"] = *req.Featured
	}
	if req.CategoryID != nil {
		updates["category_id"] = *req.CategoryID
	}

	if len(updates) > 0 {
		if err := s.productRepo.Update(id, updates); err != nil {
			return nil, fmt.Errorf("failed to update product: %w", err)
		}
	}

	return s.GetProductWithCategory(id)
}

func (s *ProductService) DeleteProduct(id string) error {
	return s.productRepo.Delete(id)
}

func (s *ProductService) SearchProducts(query string, limit int) ([]models.ProductWithRating, error) {
	if limit <= 0 {
		limit = 20
	}

	products, err := s.productRepo.Search(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}

	productsWithRating := make([]models.ProductWithRating, len(products))
	for i, product := range products {
		rating, reviewCount := s.getProductRating(product.ID)
		productsWithRating[i] = models.ProductWithRating{
			Product:       product.Product,
			Category:      product.Category,
			AverageRating: rating,
			ReviewCount:   reviewCount,
		}
	}

	return productsWithRating, nil
}

func (s *ProductService) getProductRating(productID string) (float64, int) {
	rating, count, err := s.reviewRepo.GetProductRating(productID)
	if err != nil {
		return 0, 0
	}
	return rating, count
}

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	return slug
}
