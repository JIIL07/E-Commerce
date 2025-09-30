package seeds

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

type OrderSeeder struct{}

func (s *OrderSeeder) Name() string {
	return "orders"
}

func (s *OrderSeeder) Priority() int {
	return 4
}

func (s *OrderSeeder) Seed(db *sql.DB) error {
	users, err := s.getUsers(db)
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	products, err := s.getProducts(db)
	if err != nil {
		return fmt.Errorf("failed to get products: %w", err)
	}

	if len(users) == 0 || len(products) == 0 {
		return fmt.Errorf("no users or products found for seeding orders")
	}

	rand.Seed(time.Now().UnixNano())
	orderStatuses := []string{"pending", "processing", "shipped", "delivered", "cancelled"}

	numOrders := 50 + rand.Intn(51)

	for i := 0; i < numOrders; i++ {
		user := users[rand.Intn(len(users))]

		numItems := 1 + rand.Intn(5)

		var orderID string
		total := 0.0

		status := orderStatuses[rand.Intn(len(orderStatuses))]

		createdAt := time.Now().AddDate(0, 0, -rand.Intn(180))

		err := db.QueryRow(`
			INSERT INTO orders (user_id, total, subtotal, status, shipping_address, billing_address, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`, user.ID, 0.0, 0.0, status, s.generateAddress(), s.generateAddress(), createdAt, createdAt).Scan(&orderID)

		if err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		usedProducts := make(map[string]bool)
		for j := 0; j < numItems; j++ {
			var product struct {
				ID    string
				Price float64
			}

			attempts := 0
			for {
				product = products[rand.Intn(len(products))]
				if !usedProducts[product.ID] || attempts > 10 {
					break
				}
				attempts++
			}

			usedProducts[product.ID] = true

			quantity := 1 + rand.Intn(3)
			itemTotal := product.Price * float64(quantity)
			total += itemTotal

			_, err = db.Exec(`
				INSERT INTO order_items (order_id, product_id, quantity, price, created_at)
				VALUES ($1, $2, $3, $4, $5)
			`, orderID, product.ID, quantity, product.Price, createdAt)

			if err != nil {
				return fmt.Errorf("failed to create order item: %w", err)
			}
		}

		_, err = db.Exec(`
			UPDATE orders SET total = $1, subtotal = $2 WHERE id = $3
		`, total, total, orderID)

		if err != nil {
			return fmt.Errorf("failed to update order total: %w", err)
		}
	}

	return nil
}

func (s *OrderSeeder) getUsers(db *sql.DB) ([]struct {
	ID string
}, error) {
	rows, err := db.Query("SELECT id FROM users WHERE role = 'user'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []struct {
		ID string
	}

	for rows.Next() {
		var user struct {
			ID string
		}
		if err := rows.Scan(&user.ID); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *OrderSeeder) getProducts(db *sql.DB) ([]struct {
	ID    string
	Price float64
}, error) {
	rows, err := db.Query("SELECT id, price FROM products WHERE in_stock = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []struct {
		ID    string
		Price float64
	}

	for rows.Next() {
		var product struct {
			ID    string
			Price float64
		}
		if err := rows.Scan(&product.ID, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (s *OrderSeeder) generateAddress() string {
	addresses := []string{
		`{"street": "123 Main St", "city": "New York", "state": "NY", "zip": "10001", "country": "USA"}`,
		`{"street": "456 Oak Ave", "city": "Los Angeles", "state": "CA", "zip": "90210", "country": "USA"}`,
		`{"street": "789 Pine Rd", "city": "Chicago", "state": "IL", "zip": "60601", "country": "USA"}`,
		`{"street": "321 Elm St", "city": "Houston", "state": "TX", "zip": "77001", "country": "USA"}`,
		`{"street": "654 Maple Dr", "city": "Phoenix", "state": "AZ", "zip": "85001", "country": "USA"}`,
		`{"street": "987 Cedar Ln", "city": "Philadelphia", "state": "PA", "zip": "19101", "country": "USA"}`,
		`{"street": "147 Birch Way", "city": "San Antonio", "state": "TX", "zip": "78201", "country": "USA"}`,
		`{"street": "258 Spruce St", "city": "San Diego", "state": "CA", "zip": "92101", "country": "USA"}`,
		`{"street": "369 Willow Ave", "city": "Dallas", "state": "TX", "zip": "75201", "country": "USA"}`,
		`{"street": "741 Poplar Rd", "city": "San Jose", "state": "CA", "zip": "95101", "country": "USA"}`,
	}

	return addresses[rand.Intn(len(addresses))]
}
