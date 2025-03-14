package db

import (
	"GoShort/internal/models"
	"GoShort/pkg/config"
	"GoShort/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection and applies migrations
func InitDB() {
	// Load the database connection string from environment variables
	dsn := config.Get("DATABASE_URL")

	// Connect to the database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to the database", map[string]interface{}{"error": err.Error()})
		panic("Database connection failed")
	}

	// Apply migrations
	logger.Info("Running migrations...", nil)
	if err := DB.AutoMigrate(&models.URL{}, &models.User{}); err != nil {
		logger.Error("Failed to migrate schemas", map[string]interface{}{"error": err.Error()})
		panic("Migration failed")
	}
	logger.Info("Migrations completed successfully.", nil)

	// Create initial admin user
	if err := createInitialAdminUser(); err != nil {
		logger.Error("Failed to create initial admin user", map[string]interface{}{"error": err.Error()})
	}

	logger.Info("Database connection initialized successfully.", nil)
}

// createInitialAdminUser creates an admin user if no users exist in the database
func createInitialAdminUser() error {
	var count int64
	if err := DB.Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		logger.Info("No users found, creating initial admin user...", nil)
		adminUser := models.User{
			Username: config.Get("ADMIN_USERNAME"),
			Password: config.Get("ADMIN_PASSWORD"),
			Role:     "admin",
		}

		if err := adminUser.HashPassword(); err != nil {
			return err
		}

		if err := DB.Create(&adminUser).Error; err != nil {
			return err
		}
		logger.Info("Initial admin user created successfully", nil)
	}

	return nil
}
