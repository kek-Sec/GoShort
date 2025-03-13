package v1

import (
	"encoding/json"
	"net/http"
	"time"
)

// SparklineResponse represents the response structure for sparkline data
type SparklineResponse struct {
	Success bool             `json:"success"`
	Data    []SparklinePoint `json:"data"`
}

// SparklinePoint represents a single data point in the sparkline
type SparklinePoint struct {
	Timestamp time.Time `json:"timestamp"`
	Count     int       `json:"count"`
}

// GetSparklineStats handles requests for URL shortening statistics in sparkline format
func GetSparklineStats(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	timeRange := r.URL.Query().Get("range") // e.g., "day", "week", "month"
	if timeRange == "" {
		timeRange = "day" // Default to daily stats
	}
	
	// TODO: Implement actual data retrieval from database
	// This is mock data for demonstration
	data := generateMockSparklineData(timeRange)
	
	// Prepare response
	response := SparklineResponse{
		Success: true,
		Data:    data,
	}
	
	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// generateMockSparklineData creates sample data for demonstration
func generateMockSparklineData(timeRange string) []SparklinePoint {
	now := time.Now()
	data := []SparklinePoint{}
	
	// Generate different time-based data points depending on the requested range
	var points int
	var interval time.Duration
	
	switch timeRange {
	case "day":
		points = 24
		interval = time.Hour
	case "week":
		points = 7
		interval = 24 * time.Hour
	case "month":
		points = 30
		interval = 24 * time.Hour
	default:
		points = 24
		interval = time.Hour
	}
	
	// Generate sample data points
	for i := points - 1; i >= 0; i-- {
		timestamp := now.Add(-time.Duration(i) * interval)
		// Mock count - in a real implementation, this would come from the database
		count := 5 + (i % 10) // Just a simple pattern for demonstration
		
		data = append(data, SparklinePoint{
			Timestamp: timestamp,
			Count:     count,
		})
	}
	
	return data
}
