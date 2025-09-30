package models
import (
	"time"
)
type User struct {
	ID        string    `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Name      *string   `json:"name" db:"name"`
	Password  string    `json:"-" db:"password"`
	Role      string    `json:"role" db:"role"`
	Image     *string   `json:"image" db:"image"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
type UserCreateRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name"`
}
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      *string   `json:"name"`
	Role      string    `json:"role"`
	Image     *string   `json:"image"`
	CreatedAt time.Time `json:"created_at"`
}
type AuthResponse struct {
	Message string       `json:"message"`
	User    UserResponse `json:"user"`
	Token   string       `json:"token"`
}
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Role:      u.Role,
		Image:     u.Image,
		CreatedAt: u.CreatedAt,
	}
}