package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"ecommerce-backend/internal/database"
	"ecommerce-backend/internal/middleware"
	"ecommerce-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var req models.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	var existingUser models.User
	err := database.DB.QueryRow("SELECT id FROM users WHERE email = $1", req.Email).Scan(&existingUser.ID)
	if err == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "User already exists",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to hash password",
		})
		return
	}

	var user models.User
	query := `
		INSERT INTO users (email, password, name, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, name, role, image, created_at, updated_at
	`

	var name *string
	if req.Name != "" {
		name = &req.Name
	}

	err = database.DB.QueryRow(query, req.Email, string(hashedPassword), name, "user").Scan(
		&user.ID, &user.Email, &user.Name, &user.Role, &user.Image, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to create user",
		})
		return
	}

	token, err := generateJWTToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		Message: "User created successfully",
		User:    user.ToResponse(),
		Token:   token,
	})
}

func Login(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	var user models.User
	query := `
		SELECT id, email, name, password, role, image, created_at, updated_at
		FROM users WHERE email = $1
	`
	err := database.DB.QueryRow(query, req.Email).Scan(
		&user.ID, &user.Email, &user.Name, &user.Password, &user.Role, &user.Image, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Message: "Invalid credentials",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Database error",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: "Invalid credentials",
		})
		return
	}

	token, err := generateJWTToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Message: "Login successful",
		User:    user.ToResponse(),
		Token:   token,
	})
}

func Profile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: "User not authenticated",
		})
		return
	}

	var user models.User
	query := `
		SELECT id, email, name, role, image, created_at, updated_at
		FROM users WHERE id = $1
	`
	err := database.DB.QueryRow(query, userID).Scan(
		&user.ID, &user.Email, &user.Name, &user.Role, &user.Image, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Message: "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Profile retrieved successfully",
		Data:    user.ToResponse(),
	})
}

func generateJWTToken(userID, email, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "fallback-secret-key"
	}

	claims := &middleware.Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
