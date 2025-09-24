package services

import (
	"fmt"
	"time"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repositories"
)

type CartService struct {
	cartRepo    *repositories.CartRepository
	productRepo *repositories.ProductRepository
}

func NewCartService(cartRepo *repositories.CartRepository, productRepo *repositories.ProductRepository) *CartService {
	return &CartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *CartService) AddToCart(userID, productID string, quantity int) (*models.CartItem, error) {
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if !product.InStock || product.Stock < quantity {
		return nil, fmt.Errorf("insufficient stock")
	}

	existingItem, err := s.cartRepo.GetByUserAndProduct(userID, productID)
	if err == nil {
		newQuantity := existingItem.Quantity + quantity
		if newQuantity > product.Stock {
			return nil, fmt.Errorf("insufficient stock")
		}

		updates := map[string]interface{}{
			"quantity":   newQuantity,
			"updated_at": time.Now(),
		}

		if err := s.cartRepo.Update(existingItem.ID, updates); err != nil {
			return nil, fmt.Errorf("failed to update cart item: %w", err)
		}

		updatedItem, _ := s.cartRepo.GetByID(existingItem.ID)
		return updatedItem, nil
	}

	cartItem := &models.CartItem{
		ID:        generateID(),
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.cartRepo.Create(cartItem); err != nil {
		return nil, fmt.Errorf("failed to add to cart: %w", err)
	}

	return cartItem, nil
}

func (s *CartService) GetCartItems(userID string) ([]models.CartItemWithProduct, error) {
	return s.cartRepo.GetByUserID(userID)
}

func (s *CartService) UpdateCartItem(userID, itemID string, quantity int) (*models.CartItem, error) {
	item, err := s.cartRepo.GetByID(itemID)
	if err != nil {
		return nil, fmt.Errorf("cart item not found: %w", err)
	}

	if item.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	if quantity <= 0 {
		return s.RemoveFromCart(userID, itemID)
	}

	product, err := s.productRepo.GetByID(item.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if quantity > product.Stock {
		return nil, fmt.Errorf("insufficient stock")
	}

	updates := map[string]interface{}{
		"quantity":   quantity,
		"updated_at": time.Now(),
	}

	if err := s.cartRepo.Update(itemID, updates); err != nil {
		return nil, fmt.Errorf("failed to update cart item: %w", err)
	}

	updatedItem, _ := s.cartRepo.GetByID(itemID)
	return updatedItem, nil
}

func (s *CartService) RemoveFromCart(userID, itemID string) (*models.CartItem, error) {
	item, err := s.cartRepo.GetByID(itemID)
	if err != nil {
		return nil, fmt.Errorf("cart item not found: %w", err)
	}

	if item.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	if err := s.cartRepo.Delete(itemID); err != nil {
		return nil, fmt.Errorf("failed to remove from cart: %w", err)
	}

	return item, nil
}

func (s *CartService) ClearCart(userID string) error {
	return s.cartRepo.DeleteByUserID(userID)
}

func (s *CartService) GetCartTotal(userID string) (float64, error) {
	items, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, item := range items {
		total += item.Product.Price * float64(item.Quantity)
	}

	return total, nil
}
