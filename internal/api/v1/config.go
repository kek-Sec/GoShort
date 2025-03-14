package v1

import (
	"encoding/json"
	"net/http"
	"os"
)

// Config represents the frontend configuration
type Config struct {
	Title          string `json:"title,omitempty"`
	Description    string `json:"description,omitempty"`
	Keywords       string `json:"keywords,omitempty"`
	Author         string `json:"author,omitempty"`
	ThemeColor     string `json:"themeColor,omitempty"`
	LogoText       string `json:"logoText,omitempty"`
	PrimaryColor   string `json:"primaryColor,omitempty"`
	SecondaryColor string `json:"secondaryColor,omitempty"`
	HeaderTitle    string `json:"headerTitle,omitempty"`
	FooterText     string `json:"footerText,omitempty"`
	FooterLink     string `json:"footerLink,omitempty"`
}

// GetConfig returns the frontend configuration
func GetConfig(w http.ResponseWriter, r *http.Request) {
	config := Config{}

	// Load config from environment variables
	envVars := map[string]*string{
		"BRAND_TITLE":          &config.Title,
		"BRAND_DESCRIPTION":    &config.Description,
		"BRAND_KEYWORDS":       &config.Keywords,
		"BRAND_AUTHOR":         &config.Author,
		"BRAND_THEME_COLOR":    &config.ThemeColor,
		"BRAND_LOGO_TEXT":      &config.LogoText,
		"BRAND_PRIMARY_COLOR":  &config.PrimaryColor,
		"BRAND_SECONDARY_COLOR": &config.SecondaryColor,
		"BRAND_HEADER_TITLE":   &config.HeaderTitle,
		"BRAND_FOOTER_TEXT":    &config.FooterText,
		"BRAND_FOOTER_LINK":    &config.FooterLink,
	}

	for envVar, field := range envVars {
		if value := os.Getenv(envVar); value != "" {
			*field = value
		}
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // Prevent caching

	// Encode config as JSON and send response
	json.NewEncoder(w).Encode(config)
}
