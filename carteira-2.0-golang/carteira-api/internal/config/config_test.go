package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Test default values
	os.Unsetenv("DATABASE_PATH")
	os.Unsetenv("CACHE_TTL_HOURS")
	os.Unsetenv("CACHE_ENABLED")

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if config.DatabasePath != "./portfolio.db" {
		t.Errorf("Expected default DATABASE_PATH './portfolio.db', got '%s'", config.DatabasePath)
	}

	if config.CacheTTlHours != 24 {
		t.Errorf("Expected default CACHE_TTL_HOURS 24, got %d", config.CacheTTlHours)
	}

	if config.CacheEnabled != true {
		t.Errorf("Expected default CACHE_ENABLED true, got %v", config.CacheEnabled)
	}
}

func TestLoadConfigWithEnvironment(t *testing.T) {
	// Set environment variables
	os.Setenv("DATABASE_PATH", "/test/path.db")
	os.Setenv("CACHE_TTL_HOURS", "48")
	os.Setenv("CACHE_ENABLED", "false")

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if config.DatabasePath != "/test/path.db" {
		t.Errorf("Expected DATABASE_PATH '/test/path.db', got '%s'", config.DatabasePath)
	}

	if config.CacheTTlHours != 48 {
		t.Errorf("Expected CACHE_TTL_HOURS 48, got %d", config.CacheTTlHours)
	}

	if config.CacheEnabled != false {
		t.Errorf("Expected CACHE_ENABLED false, got %v", config.CacheEnabled)
	}
}

func TestLoadConfigInvalidTTL(t *testing.T) {
	os.Setenv("CACHE_TTL_HOURS", "-5")

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if config.CacheTTlHours != 24 {
		t.Errorf("Expected default CACHE_TTL_HOURS 24 for invalid value, got %d", config.CacheTTlHours)
	}
}
