package services

import (
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repositories"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
)

type PaymentService struct {
	paymentRepo *repositories.PaymentRepository
	orderRepo   *repositories.OrderRepository
}

func NewPaymentService(paymentRepo *repositories.PaymentRepository, orderRepo *repositories.OrderRepository) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
	}
}
func (s *PaymentService) CreatePaymentIntent(userID string, req models.PaymentIntentRequest) (*models.PaymentIntentResponse, error) {
	amountInCents := int64(req.Amount * 100)
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountInCents),
		Currency: stripe.String(req.Currency),
		Metadata: map[string]string{
			"user_id": userID,
		},
	}
	if req.OrderID != nil {
		params.Metadata["order_id"] = *req.OrderID
	}
	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, err
	}
	payment := &models.Payment{
		ID:              uuid.New().String(),
		UserID:          userID,
		OrderID:         req.OrderID,
		Amount:          req.Amount,
		Currency:        req.Currency,
		Status:          models.PaymentStatusPending,
		PaymentIntentID: pi.ID,
		ClientSecret:    pi.ClientSecret,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	err = s.paymentRepo.CreatePayment(payment)
	if err != nil {
		return nil, err
	}
	return &models.PaymentIntentResponse{
		ID:           pi.ID,
		ClientSecret: pi.ClientSecret,
		Amount:       amountInCents,
		Currency:     req.Currency,
		Status:       string(pi.Status),
	}, nil
}
func (s *PaymentService) ConfirmPayment(userID string, req models.PaymentConfirmRequest) (*models.Payment, error) {
	pi, err := paymentintent.Get(req.PaymentIntentID, nil)
	if err != nil {
		return nil, err
	}
	payment, err := s.paymentRepo.GetPaymentByIntentID(req.PaymentIntentID)
	if err != nil {
		return nil, err
	}
	if payment.UserID != userID {
		return nil, fmt.Errorf("payment not found")
	}
	switch pi.Status {
	case stripe.PaymentIntentStatusSucceeded:
		payment.Status = models.PaymentStatusSucceeded
	case stripe.PaymentIntentStatusCanceled:
		payment.Status = models.PaymentStatusCancelled
	case stripe.PaymentIntentStatusRequiresPaymentMethod:
		payment.Status = models.PaymentStatusFailed
	default:
		payment.Status = models.PaymentStatusPending
	}
	payment.UpdatedAt = time.Now()
	err = s.paymentRepo.UpdatePayment(payment)
	if err != nil {
		return nil, err
	}
	if payment.Status == models.PaymentStatusSucceeded && payment.OrderID != nil {
		order, err := s.orderRepo.GetOrderByID(*payment.OrderID)
		if err == nil {
			order.Status = models.OrderStatusProcessing
			order.PaymentIntent = &payment.PaymentIntentID
			order.UpdatedAt = time.Now()
			s.orderRepo.UpdateOrder(order)
		}
	}
	return payment, nil
}
func (s *PaymentService) GetUserPayments(userID string) ([]models.Payment, error) {
	return s.paymentRepo.GetUserPayments(userID)
}
func (s *PaymentService) HandleWebhook(payload models.StripeWebhookPayload) error {
	switch payload.Type {
	case "payment_intent.succeeded":
		return s.handlePaymentIntentSucceeded(payload)
	case "payment_intent.payment_failed":
		return s.handlePaymentIntentFailed(payload)
	default:
		return nil
	}
}
func (s *PaymentService) handlePaymentIntentSucceeded(payload models.StripeWebhookPayload) error {
	paymentIntentData, ok := payload.Data["object"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid payment intent data")
	}
	paymentIntentID, ok := paymentIntentData["id"].(string)
	if !ok {
		return fmt.Errorf("invalid payment intent ID")
	}
	payment, err := s.paymentRepo.GetPaymentByIntentID(paymentIntentID)
	if err != nil {
		return err
	}
	payment.Status = models.PaymentStatusSucceeded
	payment.UpdatedAt = time.Now()
	return s.paymentRepo.UpdatePayment(payment)
}
func (s *PaymentService) handlePaymentIntentFailed(payload models.StripeWebhookPayload) error {
	paymentIntentData, ok := payload.Data["object"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid payment intent data")
	}
	paymentIntentID, ok := paymentIntentData["id"].(string)
	if !ok {
		return fmt.Errorf("invalid payment intent ID")
	}
	payment, err := s.paymentRepo.GetPaymentByIntentID(paymentIntentID)
	if err != nil {
		return err
	}
	payment.Status = models.PaymentStatusFailed
	payment.UpdatedAt = time.Now()
	return s.paymentRepo.UpdatePayment(payment)
}
