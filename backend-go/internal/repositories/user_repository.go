package repositories
import (
	"database/sql"
	"fmt"
	"strings"
	"ecommerce-backend/internal/models"
)
type UserRepository struct {
	db *sql.DB
}
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}
func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (id, email, name, password, role, image, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(query, user.ID, user.Email, user.Name, user.Password, user.Role, user.Image, user.CreatedAt, user.UpdatedAt)
	return err
}
func (r *UserRepository) GetByID(id string) (*models.User, error) {
	query := `
		SELECT id, email, name, password, role, image, created_at, updated_at
		FROM users WHERE id = $1
	`
	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.Password, &user.Role, &user.Image, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return user, err
}
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, name, password, role, image, created_at, updated_at
		FROM users WHERE email = $1
	`
	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Name, &user.Password, &user.Role, &user.Image, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return user, err
}
func (r *UserRepository) Update(id string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}
	setParts := make([]string, 0, len(updates))
	args := make([]interface{}, 0, len(updates)+1)
	argIndex := 1
	for key, value := range updates {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(setParts, ", "), argIndex)
	args = append(args, id)
	_, err := r.db.Exec(query, args...)
	return err
}
func (r *UserRepository) Delete(id string) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}
func (r *UserRepository) List(limit, offset int) ([]*models.User, error) {
	query := `
		SELECT id, email, name, password, role, image, created_at, updated_at
		FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.Name, &user.Password, &user.Role, &user.Image, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}