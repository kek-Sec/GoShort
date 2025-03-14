package v1

import (
	"GoShort/internal/auth"
	"GoShort/internal/db"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterAuthRoutes registers all authentication related routes
func RegisterAuthRoutes(router *mux.Router) {
	router.HandleFunc("/login", Login).Methods("POST", "OPTIONS")
	// Add CORS preflight handling
	router.HandleFunc("/login", PreflightHandler).Methods("OPTIONS")
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