package v1

import (
	"GoShort/internal/auth"
	"GoShort/internal/db"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// TokenValidationResponse represents the response from the token validation endpoint
type TokenValidationResponse struct {
	Valid    bool   `json:"valid"`
	Username string `json:"username,omitempty"`
	Role     string `json:"role,omitempty"`
}

// RegisterAuthRoutes registers all authentication related routes
func RegisterAuthRoutes(router *mux.Router) {
	router.HandleFunc("/login", Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/validate", ValidateToken).Methods("GET", "POST", "OPTIONS")
	// Add CORS preflight handling
	router.HandleFunc("/login", PreflightHandler).Methods("OPTIONS")
	router.HandleFunc("/validate", PreflightHandler).Methods("OPTIONS")
}

// PreflightHandler handles OPTIONS requests for CORS preflight
func PreflightHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
	w.WriteHeader(http.StatusOK)
}

// Login handles user authentication and returns JWT token
func Login(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for API responses
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
	
	// Handle the actual login
	auth.LoginHandler(db.DB)(w, r)
}

// ValidateToken checks if a JWT token is valid and returns the claims
func ValidateToken(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
	w.Header().Set("Content-Type", "application/json")

	// Get token either from Authorization header or request body
	var tokenString string
	var err error

	// First try to get token from Authorization header
	tokenString, err = auth.ExtractTokenFromRequest(r)
	
	// If not found in header, look for it in request body
	if err != nil && r.Method == "POST" {
		var requestBody struct {
			Token string `json:"token"`
		}
		
		if decodeErr := json.NewDecoder(r.Body).Decode(&requestBody); decodeErr == nil && requestBody.Token != "" {
			tokenString = requestBody.Token
			err = nil // Reset error if we found a token in the body
		}
	}

	// If we still don't have a token, return error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "No token provided. Please include a token in the Authorization header or request body.",
		})
		return
	}

	// Validate token
	claims, err := auth.VerifyToken(tokenString)
	
	response := TokenValidationResponse{}
	
	if err != nil {
		response.Valid = false
	} else {
		response.Valid = true
		response.Username = claims.Username
		response.Role = claims.Role
	}

	// Return the validation result
	json.NewEncoder(w).Encode(response)
}