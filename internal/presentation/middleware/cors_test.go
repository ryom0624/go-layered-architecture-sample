package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"layered-architecture-template/pkg/config"
)

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		corsConfig         *CORSConfig
		requestOrigin      string
		requestMethod      string
		expectedOrigin     string
		expectedMethods    string
		expectedHeaders    string
		expectedCredentials string
		expectedStatus     int
	}{
		{
			name: "Allowed origin with credentials",
			corsConfig: &CORSConfig{
				AllowedOrigins:   []string{"http://localhost:3000"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Content-Type", "Authorization"},
				AllowCredentials: true,
			},
			requestOrigin:       "http://localhost:3000",
			requestMethod:       "GET",
			expectedOrigin:      "http://localhost:3000",
			expectedMethods:     "GET,POST,PUT,DELETE,OPTIONS",
			expectedHeaders:     "Content-Type,Authorization",
			expectedCredentials: "true",
			expectedStatus:      200,
		},
		{
			name: "Wildcard origin",
			corsConfig: &CORSConfig{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET", "POST"},
				AllowedHeaders:   []string{"Content-Type"},
				AllowCredentials: false,
			},
			requestOrigin:       "http://example.com",
			requestMethod:       "GET",
			expectedOrigin:      "http://example.com",
			expectedMethods:     "GET,POST",
			expectedHeaders:     "Content-Type",
			expectedCredentials: "",
			expectedStatus:      200,
		},
		{
			name: "Disallowed origin",
			corsConfig: &CORSConfig{
				AllowedOrigins:   []string{"http://localhost:3000"},
				AllowedMethods:   []string{"GET", "POST"},
				AllowedHeaders:   []string{"Content-Type"},
				AllowCredentials: false,
			},
			requestOrigin:       "http://malicious.com",
			requestMethod:       "GET",
			expectedOrigin:      "",
			expectedMethods:     "GET,POST",
			expectedHeaders:     "Content-Type",
			expectedCredentials: "",
			expectedStatus:      200,
		},
		{
			name: "OPTIONS preflight request",
			corsConfig: &CORSConfig{
				AllowedOrigins:   []string{"http://localhost:3000"},
				AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
				AllowedHeaders:   []string{"Content-Type", "Authorization"},
				AllowCredentials: true,
			},
			requestOrigin:       "http://localhost:3000",
			requestMethod:       "OPTIONS",
			expectedOrigin:      "http://localhost:3000",
			expectedMethods:     "GET,POST,OPTIONS",
			expectedHeaders:     "Content-Type,Authorization",
			expectedCredentials: "true",
			expectedStatus:      204,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(CORS(tt.corsConfig))
			r.GET("/test", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "test"})
			})

			req := httptest.NewRequest(tt.requestMethod, "/test", nil)
			if tt.requestOrigin != "" {
				req.Header.Set("Origin", tt.requestOrigin)
			}
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedOrigin != "" {
				if origin := w.Header().Get("Access-Control-Allow-Origin"); origin != tt.expectedOrigin {
					t.Errorf("Expected origin %s, got %s", tt.expectedOrigin, origin)
				}
			} else {
				if origin := w.Header().Get("Access-Control-Allow-Origin"); origin != "" {
					t.Errorf("Expected no origin header, got %s", origin)
				}
			}

			if methods := w.Header().Get("Access-Control-Allow-Methods"); methods != tt.expectedMethods {
				t.Errorf("Expected methods %s, got %s", tt.expectedMethods, methods)
			}

			if headers := w.Header().Get("Access-Control-Allow-Headers"); headers != tt.expectedHeaders {
				t.Errorf("Expected headers %s, got %s", tt.expectedHeaders, headers)
			}

			if tt.expectedCredentials != "" {
				if credentials := w.Header().Get("Access-Control-Allow-Credentials"); credentials != tt.expectedCredentials {
					t.Errorf("Expected credentials %s, got %s", tt.expectedCredentials, credentials)
				}
			} else {
				if credentials := w.Header().Get("Access-Control-Allow-Credentials"); credentials != "" {
					t.Errorf("Expected no credentials header, got %s", credentials)
				}
			}
		})
	}
}

func TestNewCORSConfig(t *testing.T) {
	cfg := &config.Config{
		CORS: config.CORSConfig{
			AllowedOrigins:   "http://localhost:3000,http://localhost:3001",
			AllowedMethods:   "GET,POST,PUT,DELETE,OPTIONS",
			AllowedHeaders:   "Content-Type,Authorization,X-Requested-With",
			AllowCredentials: true,
		},
	}

	corsConfig := NewCORSConfig(cfg)

	expectedOrigins := []string{"http://localhost:3000", "http://localhost:3001"}
	expectedMethods := []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	expectedHeaders := []string{"Content-Type", "Authorization", "X-Requested-With"}

	if len(corsConfig.AllowedOrigins) != len(expectedOrigins) {
		t.Errorf("Expected %d origins, got %d", len(expectedOrigins), len(corsConfig.AllowedOrigins))
	}

	for i, origin := range corsConfig.AllowedOrigins {
		if origin != expectedOrigins[i] {
			t.Errorf("Expected origin %s, got %s", expectedOrigins[i], origin)
		}
	}

	if len(corsConfig.AllowedMethods) != len(expectedMethods) {
		t.Errorf("Expected %d methods, got %d", len(expectedMethods), len(corsConfig.AllowedMethods))
	}

	if len(corsConfig.AllowedHeaders) != len(expectedHeaders) {
		t.Errorf("Expected %d headers, got %d", len(expectedHeaders), len(corsConfig.AllowedHeaders))
	}

	if !corsConfig.AllowCredentials {
		t.Error("Expected AllowCredentials to be true")
	}
}

func TestParseCommaSeparated(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Normal comma-separated values",
			input:    "a,b,c",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "Values with spaces",
			input:    "a, b , c",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "Empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "Single value",
			input:    "single",
			expected: []string{"single"},
		},
		{
			name:     "Empty values",
			input:    "a,,b",
			expected: []string{"a", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseCommaSeparated(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d items, got %d", len(tt.expected), len(result))
			}

			for i, item := range result {
				if item != tt.expected[i] {
					t.Errorf("Expected item %s, got %s", tt.expected[i], item)
				}
			}
		})
	}
}

func TestIsOriginAllowed(t *testing.T) {
	tests := []struct {
		name           string
		origin         string
		allowedOrigins []string
		expected       bool
	}{
		{
			name:           "Allowed origin",
			origin:         "http://localhost:3000",
			allowedOrigins: []string{"http://localhost:3000", "http://localhost:3001"},
			expected:       true,
		},
		{
			name:           "Wildcard origin",
			origin:         "http://example.com",
			allowedOrigins: []string{"*"},
			expected:       true,
		},
		{
			name:           "Not allowed origin",
			origin:         "http://malicious.com",
			allowedOrigins: []string{"http://localhost:3000"},
			expected:       false,
		},
		{
			name:           "Empty origin",
			origin:         "",
			allowedOrigins: []string{"http://localhost:3000"},
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isOriginAllowed(tt.origin, tt.allowedOrigins)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}