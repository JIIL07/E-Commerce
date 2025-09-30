package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
	RefreshIn time.Duration
	Issuer    string
	Audience  string
}

var jwtConfig *JWTConfig

func InitJWT(secret string, expiresIn, refreshIn time.Duration, issuer, audience string) {
	jwtConfig = &JWTConfig{
		Secret:    secret,
		ExpiresIn: expiresIn,
		RefreshIn: refreshIn,
		Issuer:    issuer,
		Audience:  audience,
	}
}

func GenerateJWT(userID, email, role string) (string, error) {
	if jwtConfig == nil {
		return "", errors.New("JWT not initialized")
	}

	now := time.Now()
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtConfig.Issuer,
			Audience:  []string{jwtConfig.Audience},
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(jwtConfig.ExpiresIn)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtConfig.Secret))
}

func GenerateRefreshToken(userID string) (string, error) {
	if jwtConfig == nil {
		return "", errors.New("JWT not initialized")
	}

	now := time.Now()
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtConfig.Issuer,
			Audience:  []string{jwtConfig.Audience},
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(jwtConfig.RefreshIn)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtConfig.Secret))
}

func ValidateJWT(tokenString string) (*JWTClaims, error) {
	if jwtConfig == nil {
		return nil, errors.New("JWT not initialized")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtConfig.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func RefreshJWT(refreshToken string) (string, error) {
	claims, err := ValidateJWT(refreshToken)
	if err != nil {
		return "", err
	}

	return GenerateJWT(claims.UserID, claims.Email, claims.Role)
}

func ExtractUserID(tokenString string) (string, error) {
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return "", err
	}

	return claims.UserID, nil
}

func ExtractUserRole(tokenString string) (string, error) {
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return "", err
	}

	return claims.Role, nil
}

func IsTokenExpired(tokenString string) bool {
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return true
	}

	return time.Now().After(claims.ExpiresAt.Time)
}

func GetTokenExpiration(tokenString string) (time.Time, error) {
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return time.Time{}, err
	}

	return claims.ExpiresAt.Time, nil
}

func GetTokenIssuedAt(tokenString string) (time.Time, error) {
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return time.Time{}, err
	}

	return claims.IssuedAt.Time, nil
}

func GetTokenTTL(tokenString string) (time.Duration, error) {
	expiresAt, err := GetTokenExpiration(tokenString)
	if err != nil {
		return 0, err
	}

	ttl := time.Until(expiresAt)
	if ttl < 0 {
		return 0, errors.New("token expired")
	}

	return ttl, nil
}

func ValidateJWTWithRole(tokenString, requiredRole string) (*JWTClaims, error) {
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Role != requiredRole {
		return nil, errors.New("insufficient permissions")
	}

	return claims, nil
}

func ValidateJWTWithRoles(tokenString string, requiredRoles []string) (*JWTClaims, error) {
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return nil, err
	}

	for _, role := range requiredRoles {
		if claims.Role == role {
			return claims, nil
		}
	}

	return nil, errors.New("insufficient permissions")
}

func IsAdmin(tokenString string) bool {
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return false
	}

	return claims.Role == "admin"
}

func IsUser(tokenString string) bool {
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return false
	}

	return claims.Role == "user" || claims.Role == "admin"
}

func GetJWTConfig() *JWTConfig {
	return jwtConfig
}

func SetJWTConfig(config *JWTConfig) {
	jwtConfig = config
}
