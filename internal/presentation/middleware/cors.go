package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"layered-architecture-template/pkg/config"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

// NewCORSConfig creates a new CORS configuration from environment variables
func NewCORSConfig(cfg *config.Config) *CORSConfig {
	return &CORSConfig{
		AllowedOrigins:   parseCommaSeparated(cfg.CORS.AllowedOrigins),
		AllowedMethods:   parseCommaSeparated(cfg.CORS.AllowedMethods),
		AllowedHeaders:   parseCommaSeparated(cfg.CORS.AllowedHeaders),
		AllowCredentials: cfg.CORS.AllowCredentials,
	}
}

// CORS returns a middleware function that handles CORS
func CORS(corsConfig *CORSConfig) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Check if origin is allowed
		if isOriginAllowed(origin, corsConfig.AllowedOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		
		// Set allowed methods
		c.Header("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowedMethods, ","))
		
		// Set allowed headers
		c.Header("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowedHeaders, ","))
		
		// Set allow credentials
		if corsConfig.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		
		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})
}

// parseCommaSeparated parses a comma-separated string into a slice
func parseCommaSeparated(s string) []string {
	if s == "" {
		return []string{}
	}
	
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	
	return result
}

// isOriginAllowed checks if the origin is in the allowed list
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return false
	}
	
	for _, allowed := range allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	
	return false
}