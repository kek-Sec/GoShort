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

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// sendJSONError sends a standardized error response
func sendJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	errorResponse := ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
		Status:  status,
	}
	json.NewEncoder(w).Encode(errorResponse)
}

// LoginHandler authenticates users and issues JWT tokens
func LoginHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var loginReq LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
			sendJSONError(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		// Find user in database
		var user models.User
		result := db.Where("username = ?", loginReq.Username).First(&user)
		if result.Error != nil {
			logger.Error("Login failed: "+result.Error.Error(), nil)
			sendJSONError(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Compare passwords
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
		if err != nil {
			logger.Error("Invalid password attempt for user: "+user.Username, map[string]interface{}{
				"username": user.Username,
			})
			sendJSONError(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Generate JWT token
		token, err := GenerateToken(user)
		if err != nil {
			sendJSONError(w, "Error generating authentication token", http.StatusInternalServerError)
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
