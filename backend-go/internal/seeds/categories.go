package seeds

import (
	"database/sql"
	"fmt"
	"strings"
)

type CategorySeeder struct{}

func (s *CategorySeeder) Name() string {
	return "categories"
}

func (s *CategorySeeder) Priority() int {
	return 1
}

func (s *CategorySeeder) Seed(db *sql.DB) error {
	categories := []struct {
		name        string
		description string
		image       string
	}{
		{"Electronics", "Electronic devices and gadgets", "electronics.jpg"},
		{"Clothing", "Fashion and apparel", "clothing.jpg"},
		{"Books", "Books and literature", "books.jpg"},
		{"Home & Garden", "Home improvement and garden supplies", "home.jpg"},
		{"Sports", "Sports and fitness equipment", "sports.jpg"},
		{"Beauty & Health", "Beauty products and health supplements", "beauty.jpg"},
		{"Toys & Games", "Toys, games and entertainment", "toys.jpg"},
		{"Automotive", "Car parts and accessories", "automotive.jpg"},
		{"Food & Beverages", "Food items and drinks", "food.jpg"},
		{"Jewelry", "Fine jewelry and accessories", "jewelry.jpg"},
	}

	for _, cat := range categories {
		slug := s.generateSlug(cat.name)

		_, err := db.Exec(`
			INSERT INTO categories (name, slug, description, image, created_at, updated_at)
			VALUES ($1, $2, $3, $4, NOW(), NOW())
			ON CONFLICT (slug) DO NOTHING
		`, cat.name, slug, cat.description, cat.image)

		if err != nil {
			return fmt.Errorf("failed to insert category %s: %w", cat.name, err)
		}
	}

	return nil
}

func (s *CategorySeeder) generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "&", "and")
	slug = strings.ReplaceAll(slug, "'", "")
	slug = strings.ReplaceAll(slug, ",", "")
	return slug
}
