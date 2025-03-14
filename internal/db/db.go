// db/db.go
package db

import (
    "GoShort/internal/models"
    "GoShort/pkg/config"
    "GoShort/pkg/logger"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
    dsn := config.Get("DATABASE_URL")
    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        logger.Error("Failed to connect to DB", map[string]interface{}{"error": err.Error()})
        panic("Database connection failed")
    }

    logger.Info("Database connection established successfully.", nil)
}

// CreateInitialAdminUser after tables exist
func CreateInitialAdminUser() error {
    var count int64
    if err := DB.Model(&models.User{}).Count(&count).Error; err != nil {
        return err
    }
    if count == 0 {
        logger.Info("No users found, creating admin...", nil)
        admin := models.User{
            Username: config.Get("ADMIN_USERNAME"),
            Password: config.Get("ADMIN_PASSWORD"),
            Role:     "admin",
        }
        if err := admin.HashPassword(); err != nil {
            return err
        }
        return DB.Create(&admin).Error
    }
	logger.Info("Admin user already exists", nil)
    return nil
}
