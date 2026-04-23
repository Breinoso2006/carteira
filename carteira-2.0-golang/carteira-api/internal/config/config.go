package config

import (
	"fmt"
	"os"
)

// Config holds application configuration
type Config struct {
	DatabasePath  string
	CacheTTlHours int
	CacheEnabled  bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		DatabasePath:  "./portfolio.db",
		CacheTTlHours: 24,
		CacheEnabled:  true,
	}

	// Load database path
	if path := os.Getenv("DATABASE_PATH"); path != "" {
		config.DatabasePath = path
	}

	// Load cache TTL
	if ttlStr := os.Getenv("CACHE_TTL_HOURS"); ttlStr != "" {
		// Parse as integer
		fmt.Sscanf(ttlStr, "%d", &config.CacheTTlHours)
		if config.CacheTTlHours <= 0 {
			fmt.Printf("Warning: Invalid CACHE_TTL_HOURS (%s), using default 24\n", ttlStr)
			config.CacheTTlHours = 24
		}
	}

	// Load cache enabled flag
	if enabledStr := os.Getenv("CACHE_ENABLED"); enabledStr != "" {
		config.CacheEnabled = enabledStr != "false"
	}

	return config, nil
}
