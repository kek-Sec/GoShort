package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetSparklineStats(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedPoints int
	}{
		{
			name:           "Default range (day)",
			queryParams:    "",
			expectedPoints: 24,
		},
		{
			name:           "Day range",
			queryParams:    "range=day",
			expectedPoints: 24,
		},
		{
			name:           "Week range",
			queryParams:    "range=week",
			expectedPoints: 7,
		},
		{
			name:           "Month range",
			queryParams:    "range=month",
			expectedPoints: 30,
		},
		{
			name:           "Invalid range (should default to day)",
			queryParams:    "range=invalid",
			expectedPoints: 24,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request
			req, err := http.NewRequest("GET", "/api/v1/stats/sparkline?"+tt.queryParams, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetSparklineStats)

			// Call the handler
			handler.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}

			// Check response body
			var response SparklineResponse
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Errorf("failed to unmarshal response: %v", err)
			}

			// Verify response structure
			if !response.Success {
				t.Error("expected success to be true")
			}

			// Check number of data points
			if len(response.Data) != tt.expectedPoints {
				t.Errorf("expected %d data points, got %d", tt.expectedPoints, len(response.Data))
			}
		})
	}
}

func TestGenerateMockSparklineData(t *testing.T) {
	tests := []struct {
		name           string
		timeRange      string
		expectedPoints int
		checkInterval  bool
	}{
		{
			name:           "Day range",
			timeRange:      "day",
			expectedPoints: 24,
			checkInterval:  true,
		},
		{
			name:           "Week range",
			timeRange:      "week",
			expectedPoints: 7,
			checkInterval:  true,
		},
		{
			name:           "Month range",
			timeRange:      "month",
			expectedPoints: 30,
			checkInterval:  true,
		},
		{
			name:           "Invalid range defaults to day",
			timeRange:      "invalid",
			expectedPoints: 24,
			checkInterval:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := generateMockSparklineData(tt.timeRange)

			// Check number of data points
			if len(data) != tt.expectedPoints {
				t.Errorf("expected %d data points, got %d", tt.expectedPoints, len(data))
			}

			// Check that timestamps are in order
			for i := 1; i < len(data); i++ {
				if !data[i].Timestamp.After(data[i-1].Timestamp) {
					t.Errorf("data points not in chronological order at index %d", i)
				}
			}

			// Check intervals between timestamps if required
			if tt.checkInterval && len(data) > 1 {
				var expectedInterval time.Duration
				switch tt.timeRange {
				case "day":
					expectedInterval = time.Hour
				case "week", "month":
					expectedInterval = 24 * time.Hour
				}

				for i := 1; i < len(data); i++ {
					interval := data[i].Timestamp.Sub(data[i-1].Timestamp)
					if interval != expectedInterval {
						t.Errorf("unexpected interval at index %d: expected %v, got %v",
							i, expectedInterval, interval)
					}
				}
			}

			// Check that counts are within expected range (5-14 in the mock data)
			for i, point := range data {
				if point.Count < 5 || point.Count > 14 {
					t.Errorf("count at index %d out of expected range: got %d, expected between 5-14",
						i, point.Count)
				}
			}
		})
	}
}
