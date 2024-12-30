package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var configMap map[string]string

// Load loads environment variables from a `.env` file
func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, falling back to system environment variables")
	}

	configMap = map[string]string{
		"DATABASE_URL": os.Getenv("DATABASE_URL"),
		"PORT":         os.Getenv("PORT"),
	}
}

// Get retrieves a configuration value by key
func Get(key string) string {
	value, exists := configMap[key]
	if !exists {
		log.Fatalf("Configuration key %s not found", key)
	}
	return value
}
