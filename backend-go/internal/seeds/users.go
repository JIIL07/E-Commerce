package seeds

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"ecommerce-backend/internal/utils"
)

type UserSeeder struct{}

func (s *UserSeeder) Name() string {
	return "users"
}

func (s *UserSeeder) Priority() int {
	return 3
}

func (s *UserSeeder) Seed(db *sql.DB) error {
	users := []struct {
		email    string
		name     string
		password string
		role     string
		image    string
	}{
		{"admin@ecommerce.com", "Admin User", "admin123", "admin", "admin.jpg"},
		{"superadmin@ecommerce.com", "Super Admin", "superadmin123", "admin", "superadmin.jpg"},

		{"john.doe@example.com", "John Doe", "password123", "user", "john.jpg"},
		{"jane.smith@example.com", "Jane Smith", "password123", "user", "jane.jpg"},
		{"mike.johnson@example.com", "Mike Johnson", "password123", "user", "mike.jpg"},
		{"sarah.wilson@example.com", "Sarah Wilson", "password123", "user", "sarah.jpg"},
		{"david.brown@example.com", "David Brown", "password123", "user", "david.jpg"},
		{"lisa.davis@example.com", "Lisa Davis", "password123", "user", "lisa.jpg"},
		{"robert.miller@example.com", "Robert Miller", "password123", "user", "robert.jpg"},
		{"emily.garcia@example.com", "Emily Garcia", "password123", "user", "emily.jpg"},
		{"chris.martinez@example.com", "Chris Martinez", "password123", "user", "chris.jpg"},
		{"amanda.anderson@example.com", "Amanda Anderson", "password123", "user", "amanda.jpg"},
		{"james.taylor@example.com", "James Taylor", "password123", "user", "james.jpg"},
		{"jessica.thomas@example.com", "Jessica Thomas", "password123", "user", "jessica.jpg"},
		{"michael.jackson@example.com", "Michael Jackson", "password123", "user", "michael.jpg"},
		{"jennifer.white@example.com", "Jennifer White", "password123", "user", "jennifer.jpg"},
		{"william.harris@example.com", "William Harris", "password123", "user", "william.jpg"},
		{"linda.martin@example.com", "Linda Martin", "password123", "user", "linda.jpg"},
		{"richard.thompson@example.com", "Richard Thompson", "password123", "user", "richard.jpg"},
		{"patricia.garcia@example.com", "Patricia Garcia", "password123", "user", "patricia.jpg"},
		{"charles.martinez@example.com", "Charles Martinez", "password123", "user", "charles.jpg"},
		{"elizabeth.robinson@example.com", "Elizabeth Robinson", "password123", "user", "elizabeth.jpg"},
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 30; i++ {
		randomUser := s.generateRandomUser()
		users = append(users, randomUser)
	}

	for _, user := range users {
		hashedPassword, err := utils.HashPassword(user.password)
		if err != nil {
			return fmt.Errorf("failed to hash password for user %s: %w", user.email, err)
		}

		_, err = db.Exec(`
			INSERT INTO users (email, name, password, role, image, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
			ON CONFLICT (email) DO NOTHING
		`, user.email, user.name, hashedPassword, user.role, user.image)

		if err != nil {
			return fmt.Errorf("failed to insert user %s: %w", user.email, err)
		}
	}

	return nil
}

func (s *UserSeeder) generateRandomUser() struct {
	email    string
	name     string
	password string
	role     string
	image    string
} {
	firstNames := []string{"Alex", "Jordan", "Taylor", "Casey", "Morgan", "Riley", "Avery", "Quinn", "Sage", "River", "Phoenix", "Skyler", "Cameron", "Drew", "Blake", "Hayden", "Emery", "Finley", "Rowan", "Parker"}
	lastNames := []string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas", "Taylor", "Moore", "Jackson", "Martin"}

	firstName := firstNames[rand.Intn(len(firstNames))]
	lastName := lastNames[rand.Intn(len(lastNames))]

	email := fmt.Sprintf("%s.%s.%d@example.com",
		firstName, lastName, rand.Intn(1000))

	name := fmt.Sprintf("%s %s", firstName, lastName)

	role := "user"
	if rand.Float64() < 0.1 {
		role = "admin"
	}

	return struct {
		email    string
		name     string
		password string
		role     string
		image    string
	}{
		email:    email,
		name:     name,
		password: "password123",
		role:     role,
		image:    fmt.Sprintf("%s_%s.jpg", firstName, lastName),
	}
}
