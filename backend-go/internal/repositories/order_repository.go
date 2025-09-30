package repositories
import (
	"database/sql"
	"ecommerce-backend/internal/models"
)
type OrderRepository struct {
	db *sql.DB
}
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}
func (r *OrderRepository) CreateOrder(order *models.Order) error {
	query := `
		INSERT INTO orders (id, user_id, status, total, subtotal, tax, shipping, 
		                   shipping_address, billing_address, payment_intent, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.db.Exec(query, order.ID, order.UserID, order.Status, order.Total,
		order.Subtotal, order.Tax, order.Shipping, order.ShippingAddress,
		order.BillingAddress, order.PaymentIntent, order.CreatedAt, order.UpdatedAt)
	return err
}
func (r *OrderRepository) CreateOrderItem(item *models.OrderItem) error {
	query := `
		INSERT INTO order_items (id, order_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query, item.ID, item.OrderID, item.ProductID, item.Quantity, item.Price)
	return err
}
func (r *OrderRepository) GetOrderByID(orderID string) (*models.Order, error) {
	query := `
		SELECT id, user_id, status, total, subtotal, tax, shipping, 
		       shipping_address, billing_address, payment_intent, created_at, updated_at
		FROM orders WHERE id = $1`
	order := &models.Order{}
	err := r.db.QueryRow(query, orderID).Scan(
		&order.ID, &order.UserID, &order.Status, &order.Total,
		&order.Subtotal, &order.Tax, &order.Shipping, &order.ShippingAddress,
		&order.BillingAddress, &order.PaymentIntent, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return order, nil
}
func (r *OrderRepository) GetOrderItems(orderID string) ([]models.OrderItemWithProduct, error) {
	query := `
		SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price,
		       p.id, p.name, p.description, p.price, p.image, p.category_id,
		       p.stock_quantity, p.is_featured, p.created_at, p.updated_at
		FROM order_items oi
		JOIN products p ON oi.product_id = p.id
		WHERE oi.order_id = $1`
	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.OrderItemWithProduct
	for rows.Next() {
		var item models.OrderItemWithProduct
		var product models.Product
		err := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price,
			&product.ID, &product.Name, &product.Description, &product.Price,
			&product.Images, &product.CategoryID, &product.Stock,
			&product.Featured, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		item.Product = &models.ProductWithRating{
			Product: product,
		}
		items = append(items, item)
	}
	return items, nil
}
func (r *OrderRepository) GetUserOrders(userID string, limit, offset int) ([]models.OrderWithItems, error) {
	query := `
		SELECT id, user_id, status, total, subtotal, tax, shipping, 
		       shipping_address, billing_address, payment_intent, created_at, updated_at
		FROM orders 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []models.OrderWithItems
	for rows.Next() {
		var order models.OrderWithItems
		err := rows.Scan(
			&order.ID, &order.UserID, &order.Status, &order.Total,
			&order.Subtotal, &order.Tax, &order.Shipping, &order.ShippingAddress,
			&order.BillingAddress, &order.PaymentIntent, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items, err := r.GetOrderItems(order.ID)
		if err == nil {
			order.OrderItems = items
		}
		orders = append(orders, order)
	}
	return orders, nil
}
func (r *OrderRepository) CountUserOrders(userID string) (int, error) {
	query := `SELECT COUNT(*) FROM orders WHERE user_id = $1`
	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	return count, err
}
func (r *OrderRepository) UpdateOrder(order *models.Order) error {
	query := `
		UPDATE orders 
		SET status = $2, total = $3, subtotal = $4, tax = $5, shipping = $6,
		    shipping_address = $7, billing_address = $8, payment_intent = $9, updated_at = $10
		WHERE id = $1`
	_, err := r.db.Exec(query, order.ID, order.Status, order.Total,
		order.Subtotal, order.Tax, order.Shipping, order.ShippingAddress,
		order.BillingAddress, order.PaymentIntent, order.UpdatedAt)
	return err
}
