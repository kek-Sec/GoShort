package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Default configuration values
const (
	DefaultAdminUsername = "admin"
	DefaultAdminPassword = "admin"
	DefaultJWTSecret     = "R9em+CC+xdYN0kiWOIFlSIQvsYn3+9xus7xGC4FsTDX1+tVSHERfCAVHsMSSoi3HPvXA+xM/sQAsax45Qorryg=="
)

var configMap map[string]string

// Load loads environment variables from a `.env` file
func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, falling back to system environment variables")
	}

	// Get environment variables with defaults
	adminUsername := os.Getenv("ADMIN_USERNAME")
	if adminUsername == "" {
		adminUsername = DefaultAdminUsername
		log.Println("ADMIN_USERNAME not set, using default value")
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = DefaultAdminPassword
		log.Println("ADMIN_PASSWORD not set, using default value")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = DefaultJWTSecret
		log.Println("JWT_SECRET not set, using default value")
	}

	configMap = map[string]string{
		"DATABASE_URL":   os.Getenv("DATABASE_URL"),
		"PORT":           os.Getenv("PORT"),
		"ADMIN_USERNAME": adminUsername,
		"ADMIN_PASSWORD": adminPassword,
		"JWT_SECRET":     jwtSecret,
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
