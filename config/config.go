package config

import (
	"os"
	"strconv"
)

// Config holds all configuration values for the application
type Config struct {
	Port           string
	DiscordToken   string
	DiscordChannel string
	EnableDiscord  bool
}

// Load reads configuration from environment variables with sensible defaults
func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		DiscordToken:   getEnv("DISCORD_TOKEN", ""),
		DiscordChannel: getEnv("DISCORD_CHANNEL_ID", ""),
		EnableDiscord:  getBoolEnv("ENABLE_DISCORD", true),
	}
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getBoolEnv gets a boolean environment variable with a fallback default value
func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
