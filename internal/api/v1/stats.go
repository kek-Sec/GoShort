package v1

import (
	"GoShort/internal/auth"
	"GoShort/internal/db"
	"GoShort/internal/models"
	"GoShort/pkg/logger"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// URLStatistic represents the statistics for a single URL
type URLStatistic struct {
	ID        uint       `json:"id"`
	LongURL   string     `json:"long_url"`
	ShortURL  string     `json:"short_url"`
	CreatedAt time.Time  `json:"created_at"`
	Expiry    *time.Time `json:"expiry,omitempty"`
	Clicks    int        `json:"clicks"`
	UserID    uint       `json:"user_id"`
}

// URLStatsResponse represents the API response for URL statistics
type URLStatsResponse struct {
	URLs  []URLStatistic `json:"urls"`
	Total int            `json:"total"`
}

// GetURLStats returns statistics for URLs based on user role
func GetURLStats(w http.ResponseWriter, r *http.Request) {
	// Get claims from context (set by AuthMiddleware)
	claims, err := auth.GetClaimsFromContext(r.Context())
	if err != nil {
		logger.Error("Failed to get claims from context", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Parse query parameters
	limit := 100 // Default limit
	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0 // Default offset
	if offsetParam := r.URL.Query().Get("offset"); offsetParam != "" {
		if parsedOffset, err := strconv.Atoi(offsetParam); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	var urls []models.URL
	var total int64
	query := db.DB

	// If user is not an admin, only show their URLs
	if claims.Role != "admin" {
		// Get the user ID from the database
		var user models.User
		if err := db.DB.Where("username = ?", claims.Username).First(&user).Error; err != nil {
			logger.Error("Failed to find user", map[string]interface{}{
				"username": claims.Username,
				"error":    err.Error(),
			})
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Only retrieve URLs created by this user
		query = query.Where("user_id = ?", user.ID)
	}

	// Count total records
	if err := query.Model(&models.URL{}).Count(&total).Error; err != nil {
		logger.Error("Failed to count URLs", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Get paginated results
	if err := query.Limit(limit).Offset(offset).Find(&urls).Error; err != nil {
		logger.Error("Failed to retrieve URLs", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Convert to response format
	response := URLStatsResponse{
		Total: int(total),
		URLs:  make([]URLStatistic, len(urls)),
	}

	for i, url := range urls {
		response.URLs[i] = URLStatistic{
			ID:        url.ID,
			LongURL:   url.LongURL,
			ShortURL:  url.ShortURL,
			CreatedAt: url.CreatedAt,
			Expiry:    url.Expiry,
			Clicks:    url.Clicks,
			UserID:    url.UserID,
		}
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Failed to encode response", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
