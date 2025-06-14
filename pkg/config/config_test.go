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
	})
}
