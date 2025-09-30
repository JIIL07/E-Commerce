package seeds

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

type ReviewSeeder struct{}

func (s *ReviewSeeder) Name() string {
	return "reviews"
}

func (s *ReviewSeeder) Priority() int {
	return 5
}

func (s *ReviewSeeder) Seed(db *sql.DB) error {
	users, err := s.getUsers(db)
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	products, err := s.getProducts(db)
	if err != nil {
		return fmt.Errorf("failed to get products: %w", err)
	}

	if len(users) == 0 || len(products) == 0 {
		return fmt.Errorf("no users or products found for seeding reviews")
	}

	rand.Seed(time.Now().UnixNano())

	for _, product := range products {
		numReviews := rand.Intn(16)

		usedUsers := make(map[string]bool)

		for i := 0; i < numReviews; i++ {
			var user struct {
				ID   string
				Name string
			}

			attempts := 0
			for {
				user = users[rand.Intn(len(users))]
				if !usedUsers[user.ID] || attempts > 20 {
					break
				}
				attempts++
			}

			usedUsers[user.ID] = true

			rating := 1 + rand.Intn(5)

			comment := s.generateComment(rating, product.Name)

			createdAt := time.Now().AddDate(0, 0, -rand.Intn(90))

			_, err = db.Exec(`
				INSERT INTO reviews (user_id, product_id, rating, comment, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6)
				ON CONFLICT (user_id, product_id) DO NOTHING
			`, user.ID, product.ID, rating, comment, createdAt, createdAt)

			if err != nil {
				return fmt.Errorf("failed to create review: %w", err)
			}
		}
	}

	return nil
}

func (s *ReviewSeeder) getUsers(db *sql.DB) ([]struct {
	ID   string
	Name string
}, error) {
	rows, err := db.Query("SELECT id, name FROM users WHERE role = 'user'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []struct {
		ID   string
		Name string
	}

	for rows.Next() {
		var user struct {
			ID   string
			Name string
		}
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *ReviewSeeder) getProducts(db *sql.DB) ([]struct {
	ID   string
	Name string
}, error) {
	rows, err := db.Query("SELECT id, name FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []struct {
		ID   string
		Name string
	}

	for rows.Next() {
		var product struct {
			ID   string
			Name string
		}
		if err := rows.Scan(&product.ID, &product.Name); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (s *ReviewSeeder) generateComment(rating int, productName string) string {
	comments := map[int][]string{
		5: {
			"Excellent product! Highly recommend it.",
			"Amazing quality and fast delivery.",
			"Perfect! Exactly what I was looking for.",
			"Outstanding product, will definitely buy again.",
			"Love it! Great value for money.",
			"Fantastic quality and great customer service.",
			"Best purchase I've made in a while!",
			"Absolutely perfect, exceeded my expectations.",
			"Wonderful product, very satisfied!",
			"Top quality, highly recommended!",
		},
		4: {
			"Very good product, minor issues but overall satisfied.",
			"Good quality, would recommend with minor reservations.",
			"Nice product, works as expected.",
			"Pretty good, meets most of my needs.",
			"Solid product, good value.",
			"Good quality, fast shipping.",
			"Works well, happy with the purchase.",
			"Nice product, minor improvements could be made.",
			"Good overall, would buy again.",
			"Quality product, satisfied with purchase.",
		},
		3: {
			"Average product, nothing special.",
			"Okay quality, could be better.",
			"Decent product, meets basic needs.",
			"Average experience, neither good nor bad.",
			"Fair quality, works but has room for improvement.",
			"Okay for the price, nothing exceptional.",
			"Average product, does the job.",
			"Decent quality, could be improved.",
			"Fair value, meets expectations.",
			"Average product, works as described.",
		},
		2: {
			"Below average quality, not what I expected.",
			"Disappointed with the quality.",
			"Poor build quality, doesn't last long.",
			"Not worth the money, quality issues.",
			"Below expectations, has several problems.",
			"Poor quality control, defective item.",
			"Not satisfied, quality is lacking.",
			"Disappointing purchase, wouldn't recommend.",
			"Poor value for money.",
			"Quality issues, not as described.",
		},
		1: {
			"Terrible product, complete waste of money.",
			"Worst purchase ever, avoid this product.",
			"Completely broken upon arrival.",
			"Poor quality, doesn't work at all.",
			"Waste of money, very disappointed.",
			"Defective product, terrible experience.",
			"Awful quality, would not recommend.",
			"Complete failure, avoid at all costs.",
			"Terrible experience, poor customer service.",
			"Worst product I've ever bought.",
		},
	}

	commentList := comments[rating]
	if len(commentList) == 0 {
		return "No comment provided."
	}

	return commentList[rand.Intn(len(commentList))]
}
