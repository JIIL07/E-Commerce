package repositories
import (
	"database/sql"
	"ecommerce-backend/internal/models"
)
type PaymentRepository struct {
	db *sql.DB
}
func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}
func (r *PaymentRepository) CreatePayment(payment *models.Payment) error {
	query := `
		INSERT INTO payments (id, user_id, order_id, amount, currency, status, 
		                     payment_intent_id, client_secret, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.Exec(query, payment.ID, payment.UserID, payment.OrderID,
		payment.Amount, payment.Currency, payment.Status, payment.PaymentIntentID,
		payment.ClientSecret, payment.CreatedAt, payment.UpdatedAt)
	return err
}
func (r *PaymentRepository) GetPaymentByIntentID(paymentIntentID string) (*models.Payment, error) {
	query := `
		SELECT id, user_id, order_id, amount, currency, status, 
		       payment_intent_id, client_secret, created_at, updated_at
		FROM payments WHERE payment_intent_id = $1`
	payment := &models.Payment{}
	err := r.db.QueryRow(query, paymentIntentID).Scan(
		&payment.ID, &payment.UserID, &payment.OrderID, &payment.Amount,
		&payment.Currency, &payment.Status, &payment.PaymentIntentID,
		&payment.ClientSecret, &payment.CreatedAt, &payment.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return payment, nil
}
func (r *PaymentRepository) GetUserPayments(userID string) ([]models.Payment, error) {
	query := `
		SELECT id, user_id, order_id, amount, currency, status, 
		       payment_intent_id, client_secret, created_at, updated_at
		FROM payments 
		WHERE user_id = $1 
		ORDER BY created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		err := rows.Scan(
			&payment.ID, &payment.UserID, &payment.OrderID, &payment.Amount,
			&payment.Currency, &payment.Status, &payment.PaymentIntentID,
			&payment.ClientSecret, &payment.CreatedAt, &payment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, nil
}
func (r *PaymentRepository) UpdatePayment(payment *models.Payment) error {
	query := `
		UPDATE payments 
		SET status = $2, updated_at = $3
		WHERE id = $1`
	_, err := r.db.Exec(query, payment.ID, payment.Status, payment.UpdatedAt)
	return err
}