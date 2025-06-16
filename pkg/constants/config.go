package constants

const (
	DefaultServerPort    = "8080"
	DefaultPostgreSQLPort = "5432"
	
	DefaultAccessTokenDurationMinutes = 15
	DefaultRefreshTokenDurationDays   = 7
	HoursPerDay                       = 24
	
	RefreshTokenRandomBytes = 32
	
	// CORS defaults
	DefaultCORSAllowedOrigins   = "*"
	DefaultCORSAllowedMethods   = "GET,POST,PUT,DELETE,OPTIONS"
	DefaultCORSAllowedHeaders   = "Content-Type,Authorization,X-Requested-With"
	DefaultCORSAllowCredentials = true
)