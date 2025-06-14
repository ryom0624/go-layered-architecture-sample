package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("load config with default values", func(t *testing.T) {
		// Clear environment variables to test defaults
		os.Clearenv()

		config, err := LoadConfig()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if config == nil {
			t.Error("expected config, got nil")
		}

		// Test default server config
		if config.Server.Port != "8080" {
			t.Errorf("expected default server port '8080', got %s", config.Server.Port)
		}

		// Test default database config
		if config.Database.Driver != "postgres" {
			t.Errorf("expected default driver 'postgres', got %s", config.Database.Driver)
		}
		if config.Database.Host != "localhost" {
			t.Errorf("expected default host 'localhost', got %s", config.Database.Host)
		}
		if config.Database.Port != "5432" {
			t.Errorf("expected default port '5432', got %s", config.Database.Port)
		}
		if config.Database.User != "user" {
			t.Errorf("expected default user 'user', got %s", config.Database.User)
		}
		if config.Database.Password != "password" {
			t.Errorf("expected default password 'password', got %s", config.Database.Password)
		}
		if config.Database.Name != "database" {
			t.Errorf("expected default name 'database', got %s", config.Database.Name)
		}
	})

	t.Run("load config with environment variables", func(t *testing.T) {
		// Set environment variables
		os.Setenv("SERVER_PORT", "3000")
		os.Setenv("DB_DRIVER", "mysql")
		os.Setenv("DB_HOST", "db.example.com")
		os.Setenv("DB_PORT", "3306")
		os.Setenv("DB_USER", "testuser")
		os.Setenv("DB_PASSWORD", "testpass")
		os.Setenv("DB_NAME", "testdb")

		defer func() {
			// Clean up environment variables
			os.Clearenv()
		}()

		config, err := LoadConfig()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Test environment variable values
		if config.Server.Port != "3000" {
			t.Errorf("expected server port '3000', got %s", config.Server.Port)
		}
		if config.Database.Driver != "mysql" {
			t.Errorf("expected driver 'mysql', got %s", config.Database.Driver)
		}
		if config.Database.Host != "db.example.com" {
			t.Errorf("expected host 'db.example.com', got %s", config.Database.Host)
		}
		if config.Database.Port != "3306" {
			t.Errorf("expected port '3306', got %s", config.Database.Port)
		}
		if config.Database.User != "testuser" {
			t.Errorf("expected user 'testuser', got %s", config.Database.User)
		}
		if config.Database.Password != "testpass" {
			t.Errorf("expected password 'testpass', got %s", config.Database.Password)
		}
		if config.Database.Name != "testdb" {
			t.Errorf("expected name 'testdb', got %s", config.Database.Name)
		}
	})

	t.Run("load config with partial environment variables", func(t *testing.T) {
		// Set only some environment variables
		os.Clearenv()
		os.Setenv("SERVER_PORT", "9000")
		os.Setenv("DB_HOST", "partial.example.com")

		defer func() {
			os.Clearenv()
		}()

		config, err := LoadConfig()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Test mixed values (env vars + defaults)
		if config.Server.Port != "9000" {
			t.Errorf("expected server port '9000', got %s", config.Server.Port)
		}
		if config.Database.Host != "partial.example.com" {
			t.Errorf("expected host 'partial.example.com', got %s", config.Database.Host)
		}
		// These should still be defaults
		if config.Database.Driver != "postgres" {
			t.Errorf("expected default driver 'postgres', got %s", config.Database.Driver)
		}
		if config.Database.Port != "5432" {
			t.Errorf("expected default port '5432', got %s", config.Database.Port)
		}
	})
}

func TestGetEnv(t *testing.T) {
	t.Run("get existing environment variable", func(t *testing.T) {
		os.Setenv("TEST_VAR", "test_value")
		defer os.Unsetenv("TEST_VAR")

		value := getEnv("TEST_VAR", "default")
		if value != "test_value" {
			t.Errorf("expected 'test_value', got %s", value)
		}
	})

	t.Run("get non-existent environment variable", func(t *testing.T) {
		os.Unsetenv("NON_EXISTENT_VAR")

		value := getEnv("NON_EXISTENT_VAR", "default_value")
		if value != "default_value" {
			t.Errorf("expected 'default_value', got %s", value)
		}
	})

	t.Run("get empty environment variable", func(t *testing.T) {
		os.Setenv("EMPTY_VAR", "")
		defer os.Unsetenv("EMPTY_VAR")

		value := getEnv("EMPTY_VAR", "default_value")
		if value != "default_value" {
			t.Errorf("expected 'default_value', got %s", value)
		}
	})
}

func TestGetEnvInt(t *testing.T) {
	t.Run("get valid integer environment variable", func(t *testing.T) {
		os.Setenv("TEST_INT", "42")
		defer os.Unsetenv("TEST_INT")

		value := getEnvInt("TEST_INT", 10)
		if value != 42 {
			t.Errorf("expected 42, got %d", value)
		}
	})

	t.Run("get non-existent integer environment variable", func(t *testing.T) {
		os.Unsetenv("NON_EXISTENT_INT")

		value := getEnvInt("NON_EXISTENT_INT", 100)
		if value != 100 {
			t.Errorf("expected 100, got %d", value)
		}
	})

	t.Run("get invalid integer environment variable", func(t *testing.T) {
		os.Setenv("INVALID_INT", "not_a_number")
		defer os.Unsetenv("INVALID_INT")

		value := getEnvInt("INVALID_INT", 50)
		if value != 50 {
			t.Errorf("expected default value 50, got %d", value)
		}
	})

	t.Run("get empty integer environment variable", func(t *testing.T) {
		os.Setenv("EMPTY_INT", "")
		defer os.Unsetenv("EMPTY_INT")

		value := getEnvInt("EMPTY_INT", 25)
		if value != 25 {
			t.Errorf("expected default value 25, got %d", value)
		}
	})

	t.Run("get zero integer environment variable", func(t *testing.T) {
		os.Setenv("ZERO_INT", "0")
		defer os.Unsetenv("ZERO_INT")

		value := getEnvInt("ZERO_INT", 99)
		if value != 0 {
			t.Errorf("expected 0, got %d", value)
		}
	})

	t.Run("get negative integer environment variable", func(t *testing.T) {
		os.Setenv("NEGATIVE_INT", "-123")
		defer os.Unsetenv("NEGATIVE_INT")

		value := getEnvInt("NEGATIVE_INT", 1)
		if value != -123 {
			t.Errorf("expected -123, got %d", value)
		}
	})
}