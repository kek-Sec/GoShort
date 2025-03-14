package models

import "time"

// URL represents the structure of a shortened URL
type URL struct {
	ID        uint      `gorm:"primaryKey"`
	LongURL   string    `gorm:"not null"`
	ShortURL  string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Expiry    *time.Time `gorm:"type:timestamp"` // Optional expiry date
	Clicks    int       `gorm:"default:0"`
	UserID    uint      `gorm:"index"` // Add UserID field to track who created the short URL
}
