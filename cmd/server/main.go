package main

import (
	"GoShort/pkg/config"
	"GoShort/pkg/logger"
	"GoShort/internal/db"
	"log"
	"net/http"
)

var version = "dev"

func main() {
	// Load configuration
	config.Load()

	// Initialize logger
	logger.Init()

	// Initialize database
	db.InitDB()

	// Set up HTTP server
	router := setupRouter()

	// Start the server
	log.Println("Starting GoShort server on port 8080...")
	log.Printf("Version: %s\n", version)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
