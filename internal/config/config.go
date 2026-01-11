package config

import (
	"os"
)

// Config holds the agent configuration
type Config struct {
	Port        string
	AdminToken  string
	TenantsRoot string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Port:        getEnv("Q8_AGENT_PORT", "8080"),
		AdminToken:  getEnv("Q8_AGENT_ADMIN_TOKEN", "change-me"),
		TenantsRoot: getEnv("Q8_TENANTS_ROOT", "/opt/tenants"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
