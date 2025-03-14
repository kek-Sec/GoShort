package auth

import (
	"GoShort/pkg/config"
	"GoShort/pkg/logger"
	"GoShort/internal/models"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWT claim structure
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Constants for JWT
const (
	TokenExpiry = 24 * time.Hour
)

// GenerateToken creates a new JWT token for a user
func GenerateToken(user models.User) (string, error) {
	// Get JWT secret from config
	secret := config.Get("JWT_SECRET")
	if secret == "" {
		logger.Error("JWT_SECRET not set in configuration", nil)
		return "", errors.New("JWT secret not configured")
	}

	// Create the claims
	claims := &Claims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "goshort",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	// Generate the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		logger.Error("Failed to sign JWT token: "+err.Error(), nil)
		return "", err
	}

	return tokenString, nil
}

// VerifyToken validates a JWT token and returns the claims
func VerifyToken(tokenString string) (*Claims, error) {
	// Get JWT secret from config
	secret := config.Get("JWT_SECRET")
	if secret == "" {
		logger.Error("JWT_SECRET not set in configuration", nil)
		return nil, errors.New("JWT secret not configured")
	}

	// Parse and validate the token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		logger.Error("Failed to parse JWT token: "+err.Error(), nil)
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ExtractTokenFromRequest gets the token from the Authorization header
func ExtractTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	// Check if the header has the Bearer prefix
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("authorization header format must be 'Bearer {token}'")
	}

	return parts[1], nil
}

// GetUserFromToken extracts user information from a token
func GetUserFromToken(tokenString string) (*models.User, error) {
	claims, err := VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Create a user with the information from the token
	user := &models.User{
		Username: claims.Username,
		Role:     claims.Role,
	}

	return user, nil
}
