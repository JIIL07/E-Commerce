package seeds

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/lib/pq"
)

type ProductSeeder struct{}

func (s *ProductSeeder) Name() string {
	return "products"
}

func (s *ProductSeeder) Priority() int {
	return 2
}

func (s *ProductSeeder) Seed(db *sql.DB) error {
	products := []struct {
		name         string
		description  string
		price        float64
		categorySlug string
		images       []string
		stock        int
		featured     bool
	}{
		{"iPhone 15 Pro", "Latest iPhone with titanium design and A17 Pro chip", 999.99, "electronics", []string{"iphone15_pro.jpg", "iphone15_pro_back.jpg"}, 50, true},
		{"Samsung Galaxy S24 Ultra", "Premium Android smartphone with S Pen", 1199.99, "electronics", []string{"galaxy_s24_ultra.jpg"}, 30, true},
		{"MacBook Pro M3", "Professional laptop with M3 chip", 1999.99, "electronics", []string{"macbook_pro_m3.jpg"}, 20, true},
		{"iPad Air", "Powerful tablet for work and creativity", 599.99, "electronics", []string{"ipad_air.jpg"}, 40, false},
		{"AirPods Pro", "Wireless earbuds with noise cancellation", 249.99, "electronics", []string{"airpods_pro.jpg"}, 100, false},
		{"Sony WH-1000XM5", "Premium noise-canceling headphones", 399.99, "electronics", []string{"sony_wh1000xm5.jpg"}, 25, false},

		{"Nike Air Max 270", "Comfortable running shoes with Max Air", 150.99, "clothing", []string{"nike_air_max_270.jpg"}, 100, false},
		{"Adidas Ultraboost 22", "High-performance running shoes", 180.99, "clothing", []string{"adidas_ultraboost_22.jpg"}, 80, false},
		{"Levi's 501 Jeans", "Classic straight-fit jeans", 89.99, "clothing", []string{"levis_501.jpg"}, 200, false},
		{"Uniqlo Heattech T-Shirt", "Thermal base layer t-shirt", 19.99, "clothing", []string{"uniqlo_heattech.jpg"}, 300, false},
		{"Patagonia Fleece Jacket", "Sustainable fleece jacket", 99.99, "clothing", []string{"patagonia_fleece.jpg"}, 60, false},

		{"Clean Code", "A Handbook of Agile Software Craftsmanship", 49.99, "books", []string{"clean_code.jpg"}, 75, true},
		{"JavaScript: The Good Parts", "Essential JavaScript concepts", 39.99, "books", []string{"js_good_parts.jpg"}, 60, false},
		{"Python Crash Course", "Learn Python programming", 44.99, "books", []string{"python_crash_course.jpg"}, 80, false},
		{"Design Patterns", "Elements of Reusable Object-Oriented Software", 59.99, "books", []string{"design_patterns.jpg"}, 45, false},
		{"The Pragmatic Programmer", "Your journey to mastery", 54.99, "books", []string{"pragmatic_programmer.jpg"}, 55, false},

		{"Dyson V15 Detect", "Cordless vacuum with laser dust detection", 749.99, "home-and-garden", []string{"dyson_v15.jpg"}, 15, true},
		{"KitchenAid Stand Mixer", "Professional stand mixer", 399.99, "home-and-garden", []string{"kitchenaid_mixer.jpg"}, 20, false},
		{"Philips Hue Starter Kit", "Smart lighting system", 199.99, "home-and-garden", []string{"philips_hue.jpg"}, 30, false},
		{"Weber Gas Grill", "3-burner gas grill", 449.99, "home-and-garden", []string{"weber_grill.jpg"}, 10, false},

		{"Peloton Bike", "Interactive fitness bike", 1495.00, "sports", []string{"peloton_bike.jpg"}, 5, true},
		{"Bowflex Adjustable Dumbbells", "Space-saving adjustable weights", 549.99, "sports", []string{"bowflex_dumbbells.jpg"}, 25, false},
		{"Yoga Mat Premium", "Non-slip yoga mat", 39.99, "sports", []string{"yoga_mat_premium.jpg"}, 80, false},
		{"Resistance Bands Set", "Complete resistance training set", 29.99, "sports", []string{"resistance_bands.jpg"}, 120, false},

		{"La Mer Moisturizing Cream", "Luxury skincare cream", 195.00, "beauty-and-health", []string{"la_mer_cream.jpg"}, 20, false},
		{"Oral-B Electric Toothbrush", "Advanced electric toothbrush", 89.99, "beauty-and-health", []string{"oral_b_toothbrush.jpg"}, 50, false},
		{"Vitamins Multivitamin", "Daily multivitamin supplement", 24.99, "beauty-and-health", []string{"multivitamin.jpg"}, 200, false},

		{"LEGO Creator Set", "Build and rebuild creative sets", 79.99, "toys-and-games", []string{"lego_creator.jpg"}, 100, false},
		{"Nintendo Switch", "Portable gaming console", 299.99, "toys-and-games", []string{"nintendo_switch.jpg"}, 30, true},
		{"Monopoly Classic", "Classic board game", 24.99, "toys-and-games", []string{"monopoly.jpg"}, 150, false},

		{"Car Phone Mount", "Magnetic phone mount for car", 19.99, "automotive", []string{"car_phone_mount.jpg"}, 200, false},
		{"Dash Cam", "HD dashboard camera", 99.99, "automotive", []string{"dash_cam.jpg"}, 50, false},
		{"Car Air Freshener", "Long-lasting air freshener", 8.99, "automotive", []string{"air_freshener.jpg"}, 500, false},

		{"Organic Coffee Beans", "Premium organic coffee", 24.99, "food-and-beverages", []string{"organic_coffee.jpg"}, 100, false},
		{"Protein Powder", "Whey protein supplement", 49.99, "food-and-beverages", []string{"protein_powder.jpg"}, 80, false},
		{"Green Tea Set", "Premium green tea collection", 34.99, "food-and-beverages", []string{"green_tea_set.jpg"}, 60, false},

		{"Gold Necklace", "14k gold chain necklace", 299.99, "jewelry", []string{"gold_necklace.jpg"}, 25, false},
		{"Silver Ring", "Sterling silver ring", 89.99, "jewelry", []string{"silver_ring.jpg"}, 50, false},
		{"Pearl Earrings", "Classic pearl earrings", 149.99, "jewelry", []string{"pearl_earrings.jpg"}, 30, false},
	}

	categoryMap, err := s.getCategoryMap(db)
	if err != nil {
		return fmt.Errorf("failed to get category map: %w", err)
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 20; i++ {
		randomProduct := s.generateRandomProduct(categoryMap)
		products = append(products, randomProduct)
	}

	for _, product := range products {
		categoryID, exists := categoryMap[product.categorySlug]
		if !exists {
			continue // Пропускаем продукты с несуществующими категориями
		}

		slug := s.generateSlug(product.name)

		_, err := db.Exec(`
			INSERT INTO products (name, slug, description, price, category_id, images, stock, featured, in_stock, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
			ON CONFLICT (slug) DO NOTHING
		`, product.name, slug, product.description, product.price, categoryID, pq.Array(product.images), product.stock, product.featured, product.stock > 0)

		if err != nil {
			return fmt.Errorf("failed to insert product %s: %w", product.name, err)
		}
	}

	return nil
}

func (s *ProductSeeder) getCategoryMap(db *sql.DB) (map[string]string, error) {
	rows, err := db.Query("SELECT id, slug FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categoryMap := make(map[string]string)
	for rows.Next() {
		var id, slug string
		if err := rows.Scan(&id, &slug); err != nil {
			return nil, err
		}
		categoryMap[slug] = id
	}

	return categoryMap, nil
}

func (s *ProductSeeder) generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "&", "and")
	slug = strings.ReplaceAll(slug, "'", "")
	slug = strings.ReplaceAll(slug, ",", "")
	slug = strings.ReplaceAll(slug, ":", "")
	return slug
}

func (s *ProductSeeder) generateRandomProduct(categoryMap map[string]string) struct {
	name         string
	description  string
	price        float64
	categorySlug string
	images       []string
	stock        int
	featured     bool
} {
	names := []string{"Premium Widget", "Deluxe Gadget", "Ultra Tool", "Pro Device", "Smart Component"}
	descriptions := []string{"High-quality product", "Professional grade", "Advanced technology", "Premium materials", "Innovative design"}

	var categorySlug string
	for slug := range categoryMap {
		categorySlug = slug
		break
	}

	price := 10.0 + rand.Float64()*490.0

	stock := rand.Intn(101)

	featured := rand.Float64() < 0.1

	name := names[rand.Intn(len(names))]
	description := descriptions[rand.Intn(len(descriptions))]

	return struct {
		name         string
		description  string
		price        float64
		categorySlug string
		images       []string
		stock        int
		featured     bool
	}{
		name:         name,
		description:  description,
		price:        price,
		categorySlug: categorySlug,
		images:       []string{"default_product.jpg"},
		stock:        stock,
		featured:     featured,
	}
}
