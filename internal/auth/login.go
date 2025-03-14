package auth

import (
	"GoShort/internal/models"
	"GoShort/pkg/logger"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the login response body
type LoginResponse struct {
	Token   string `json:"token"`
	Expires string `json:"expires"`
	User    struct {
		Username string `json:"username"`
		Role     string `json:"role"`
	} `json:"user"`
}

// LoginHandler authenticates users and issues JWT tokens
func LoginHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var loginReq LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Find user in database
		var user models.User
		result := db.Where("username = ?", loginReq.Username).First(&user)
		if result.Error != nil {
			logger.Error("Login failed: "+result.Error.Error(), nil)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Compare passwords
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
		if err != nil {
			logger.Error("Invalid password attempt for user: "+user.Username, map[string]interface{}{
				"username": user.Username,
			})
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Generate JWT token
		token, err := GenerateToken(user)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		// Update last login time
		db.Model(&user).Update("last_login", time.Now())

		// Create response
		resp := LoginResponse{
			Token:   token,
			Expires: time.Now().Add(TokenExpiry).Format(time.RFC3339),
		}
		resp.User.Username = user.Username
		resp.User.Role = user.Role

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		
		// Return the response
		json.NewEncoder(w).Encode(resp)
	}
}
