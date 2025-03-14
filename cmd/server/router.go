package main

import (
	v1 "GoShort/internal/api/v1"

	"github.com/gorilla/mux"
)

// setupRouter initializes the HTTP router and routes
func setupRouter() *mux.Router {
	router := mux.NewRouter()

	// V1 Routes
	apiV1 := router.PathPrefix("/v1").Subrouter()
	apiV1.HandleFunc("/shorten", v1.ShortenURL).Methods("POST")
	apiV1.HandleFunc("/config", v1.GetConfig).Methods("GET") // Add config endpoint

// Auth Routes
	authRouter := apiV1.PathPrefix("/auth").Subrouter()
	v1.RegisterAuthRoutes(authRouter)

	// Redirect Route (catch-all)
	router.HandleFunc("/{shortURL}", v1.RedirectURL).Methods("GET")

	return router
}
