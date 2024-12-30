package v1

import (
	"GoShort/internal/db"
	"GoShort/internal/models"
	"GoShort/internal/utils"
	"encoding/json"
	"net/http"
	"regexp"
	"time"
)

// ShortenRequest represents the request payload for URL shortening
type ShortenRequest struct {
	LongURL   string `json:"long_url"`
	CustomURL string `json:"custom_url,omitempty"`
	Expiry    string `json:"expiry,omitempty"` // Optional expiry date
}

// ShortenResponse represents the response payload for URL shortening
type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

// validateCustomURL ensures the custom URL does not contain spaces or illegal characters
func validateCustomURL(customURL string) bool {
	// Define a regex for allowed characters (alphanumeric and dashes/underscores only)
	validURLPattern := `^[a-zA-Z0-9_-]+$`
	re := regexp.MustCompile(validURLPattern)
	return re.MatchString(customURL)
}

// ShortenURL handles the URL shortening request
func ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the long URL
	if !utils.ValidateURL(req.LongURL) {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	// Use custom URL if provided, otherwise generate a new one
	shortURL := req.CustomURL
	if shortURL != "" {
		if !validateCustomURL(shortURL) {
			http.Error(w, "Custom URL contains invalid characters", http.StatusBadRequest)
			return
		}
		// Check if the custom URL already exists
		var existingURL models.URL
		if err := db.DB.Where("short_url = ?", shortURL).First(&existingURL).Error; err == nil {
			http.Error(w, "Custom URL is already taken", http.StatusConflict)
			return
		}
	} else {
		shortURL = utils.GenerateShortURL()
	}

	// Parse expiry if provided
	var expiry *time.Time
	if req.Expiry != "" {
		parsedExpiry, err := time.Parse(time.RFC3339, req.Expiry)
		if err != nil {
			http.Error(w, "Invalid expiry format", http.StatusBadRequest)
			return
		}
		expiry = &parsedExpiry
	}

	// Save the URL to the database
	url := models.URL{
		LongURL:  req.LongURL,
		ShortURL: shortURL,
		Expiry:   expiry,
	}
	if err := db.DB.Create(&url).Error; err != nil {
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	// Return the shortened URL
	resp := ShortenResponse{
		ShortURL: shortURL,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
