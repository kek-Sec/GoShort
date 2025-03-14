package v1

import (
	"GoShort/internal/auth"
	"GoShort/internal/db"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterAuthRoutes registers all authentication related routes
func RegisterAuthRoutes(router *mux.Router) {
	router.HandleFunc("/login", Login).Methods("POST")
}

// Login handles user authentication and returns JWT token
func Login(w http.ResponseWriter, r *http.Request) {
	auth.LoginHandler(db.DB)(w, r)
}