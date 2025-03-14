package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	// Test cases
	tests := []struct {
		name        string
		envVars     map[string]string
		expectedKey string
		expectedVal string
	}{
		{
			name: "Default Config Returns Empty When No Env Vars",
			// No env vars set
			expectedKey: "",
			expectedVal: "",
		},
		{
			name: "Returns Title From Environment",
			envVars: map[string]string{
				"BRAND_TITLE": "Custom Title",
			},
			expectedKey: "title",
			expectedVal: "Custom Title",
		},
		{
			name: "Returns Primary Color From Environment",
			envVars: map[string]string{
				"BRAND_PRIMARY_COLOR": "#ff0000",
			},
			expectedKey: "primaryColor",
			expectedVal: "#ff0000",
		},
		{
			name: "Returns Multiple Config Values",
			envVars: map[string]string{
				"BRAND_TITLE":          "Custom Title",
				"BRAND_DESCRIPTION":    "Custom Description",
				"BRAND_PRIMARY_COLOR":  "#ff0000",
				"BRAND_HEADER_TITLE":   "Custom Header",
				"BRAND_SECONDARY_COLOR": "#00ff00",
			},
			expectedKey: "headerTitle", // We'll check just this one
			expectedVal: "Custom Header",
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment variables
			os.Clearenv()

			// Set environment variables for this test
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			// Create request
			req, err := http.NewRequest("GET", "/v1/config", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create response recorder
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetConfig)

			// Serve request
			handler.ServeHTTP(rr, req)

			// Check status code
			assert.Equal(t, http.StatusOK, rr.Code)

			// Check Content-Type
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
			assert.Equal(t, "no-cache, no-store, must-revalidate", rr.Header().Get("Cache-Control"))

			// If we're not expecting any specific value, just verify it's valid JSON
			if tt.expectedKey == "" {
				var result map[string]interface{}
				err = json.Unmarshal(rr.Body.Bytes(), &result)
				assert.NoError(t, err, "Response should be valid JSON")
				return
			}

			// Parse the response
			var result Config
			err = json.Unmarshal(rr.Body.Bytes(), &result)
			assert.NoError(t, err)

			// Check the specific field we're testing for
			switch tt.expectedKey {
			case "title":
				assert.Equal(t, tt.expectedVal, result.Title)
			case "description":
				assert.Equal(t, tt.expectedVal, result.Description)
			case "keywords":
				assert.Equal(t, tt.expectedVal, result.Keywords)
			case "author":
				assert.Equal(t, tt.expectedVal, result.Author)
			case "themeColor":
				assert.Equal(t, tt.expectedVal, result.ThemeColor)
			case "logoText":
				assert.Equal(t, tt.expectedVal, result.LogoText)
			case "primaryColor":
				assert.Equal(t, tt.expectedVal, result.PrimaryColor)
			case "secondaryColor":
				assert.Equal(t, tt.expectedVal, result.SecondaryColor)
			case "headerTitle":
				assert.Equal(t, tt.expectedVal, result.HeaderTitle)
			case "footerText":
				assert.Equal(t, tt.expectedVal, result.FooterText)
			case "footerLink":
				assert.Equal(t, tt.expectedVal, result.FooterLink)
			}
		})
	}
}

func TestGetConfigMultipleValues(t *testing.T) {
	// Clear environment variables
	os.Clearenv()

	// Set multiple environment variables
	os.Setenv("BRAND_TITLE", "Test Title")
	os.Setenv("BRAND_PRIMARY_COLOR", "#ff0000")
	os.Setenv("BRAND_SECONDARY_COLOR", "#00ff00")
	os.Setenv("BRAND_FOOTER_TEXT", "Custom Footer")

	// Create request
	req, err := http.NewRequest("GET", "/v1/config", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetConfig)

	// Serve request
	handler.ServeHTTP(rr, req)

	// Check status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response
	var result Config
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	assert.NoError(t, err)

	// Check all expected values
	assert.Equal(t, "Test Title", result.Title)
	assert.Equal(t, "#ff0000", result.PrimaryColor)
	assert.Equal(t, "#00ff00", result.SecondaryColor)
	assert.Equal(t, "Custom Footer", result.FooterText)

	// Check that others are empty
	assert.Empty(t, result.Description)
	assert.Empty(t, result.Keywords)
}
