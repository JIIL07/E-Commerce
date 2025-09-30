package services

import (
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repositories"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type OrderService struct {
	orderRepo   *repositories.OrderRepository
	cartRepo    *repositories.CartRepository
	productRepo *repositories.ProductRepository
}

func NewOrderService(orderRepo *repositories.OrderRepository, cartRepo *repositories.CartRepository, productRepo *repositories.ProductRepository) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}
func (s *OrderService) GetUserOrders(userID string, page, limit int) ([]models.OrderWithItems, int, error) {
	offset := (page - 1) * limit
	orders, err := s.orderRepo.GetUserOrders(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.orderRepo.CountUserOrders(userID)
	if err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}
func (s *OrderService) GetOrderByID(orderID, userID string) (*models.OrderWithItems, error) {
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return nil, err
	}
	if order.UserID != userID {
		return nil, fmt.Errorf("order not found")
	}
	orderItems, err := s.orderRepo.GetOrderItems(orderID)
	if err != nil {
		return nil, err
	}
	orderWithItems := &models.OrderWithItems{
		Order:      *order,
		OrderItems: orderItems,
	}
	return orderWithItems, nil
}
func (s *OrderService) CreateOrder(userID string, req models.OrderCreateRequest) (*models.OrderWithItems, error) {
	cartItems, err := s.cartRepo.GetUserCartItems(userID)
	if err != nil {
		return nil, err
	}
	if len(cartItems) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}
	var subtotal float64
	var orderItems []models.OrderItem
	for _, item := range cartItems {
		product, err := s.productRepo.GetProductByID(item.ProductID)
		if err != nil {
			return nil, err
		}
		itemTotal := product.Price * float64(item.Quantity)
		subtotal += itemTotal
		orderItems = append(orderItems, models.OrderItem{
			ID:        uuid.New().String(),
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
	}
	tax := subtotal * 0.1 // 10% tax
	shipping := 10.0      // Fixed shipping cost
	total := subtotal + tax + shipping
	order := &models.Order{
		ID:              uuid.New().String(),
		UserID:          userID,
		Status:          models.OrderStatusPending,
		Total:           total,
		Subtotal:        subtotal,
		Tax:             tax,
		Shipping:        shipping,
		ShippingAddress: req.ShippingAddress,
		BillingAddress:  req.BillingAddress,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	err = s.orderRepo.CreateOrder(order)
	if err != nil {
		return nil, err
	}
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
		err = s.orderRepo.CreateOrderItem(&orderItems[i])
		if err != nil {
			return nil, err
		}
	}
	err = s.cartRepo.ClearUserCart(userID)
	if err != nil {
		return nil, err
	}
	orderItemsWithProduct := make([]models.OrderItemWithProduct, len(orderItems))
	for i, item := range orderItems {
		orderItemsWithProduct[i] = models.OrderItemWithProduct{
			OrderItem: item,
		}
	}
	orderWithItems := &models.OrderWithItems{
		Order:      *order,
		OrderItems: orderItemsWithProduct,
	}
	return orderWithItems, nil
}
func (s *OrderService) UpdateOrderStatus(orderID string, status models.OrderStatus) (*models.Order, error) {
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return nil, err
	}
	order.Status = status
	order.UpdatedAt = time.Now()
	err = s.orderRepo.UpdateOrder(order)
	if err != nil {
		return nil, err
	}
	return order, nil
}
func (s *OrderService) CancelOrder(orderID, userID string) error {
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return err
	}
	if order.UserID != userID {
		return fmt.Errorf("order not found")
	}
	if order.Status != models.OrderStatusPending {
		return fmt.Errorf("order cannot be cancelled")
	}
	order.Status = models.OrderStatusCancelled
	order.UpdatedAt = time.Now()
	return s.orderRepo.UpdateOrder(order)
}
