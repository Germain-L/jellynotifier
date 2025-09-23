package config

import (
	"fmt"
	"log"
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

// Load reads configuration from environment variables with validation
func Load() (*Config, error) {
	log.Println("[DEBUG] [CONFIG] Starting configuration loading...")

	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		DiscordToken:   getEnv("DISCORD_TOKEN", ""),
		DiscordChannel: getEnv("DISCORD_CHANNEL_ID", ""),
		EnableDiscord:  getBoolEnv("ENABLE_DISCORD", true),
	}

	log.Printf("[DEBUG] [CONFIG] PORT: %s", cfg.Port)
	log.Printf("[DEBUG] [CONFIG] ENABLE_DISCORD: %t", cfg.EnableDiscord)
	log.Printf("[DEBUG] [CONFIG] DISCORD_TOKEN present: %t", cfg.DiscordToken != "")
	log.Printf("[DEBUG] [CONFIG] DISCORD_CHANNEL_ID present: %t", cfg.DiscordChannel != "")

	// Validate required Discord configuration if Discord is enabled
	if cfg.EnableDiscord {
		log.Println("[DEBUG] [CONFIG] Discord is enabled, validating required configuration...")
		if cfg.DiscordToken == "" {
			log.Println("[ERROR] [CONFIG] DISCORD_TOKEN is required when Discord is enabled")
			return nil, fmt.Errorf("DISCORD_TOKEN environment variable is required when Discord is enabled")
		}
		if cfg.DiscordChannel == "" {
			log.Println("[ERROR] [CONFIG] DISCORD_CHANNEL_ID is required when Discord is enabled")
			return nil, fmt.Errorf("DISCORD_CHANNEL_ID environment variable is required when Discord is enabled")
		}
		log.Println("[DEBUG] [CONFIG] Discord configuration validation passed")
	} else {
		log.Println("[DEBUG] [CONFIG] Discord is disabled, skipping Discord configuration validation")
	}

	log.Println("[DEBUG] [CONFIG] Configuration loading completed successfully")
	return cfg, nil
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		log.Printf("[DEBUG] [CONFIG] Environment variable %s found", key)
		return value
	}
	log.Printf("[DEBUG] [CONFIG] Environment variable %s not found, using default: %s", key, defaultValue)
	return defaultValue
}

// getBoolEnv gets a boolean environment variable with a fallback default value
func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		log.Printf("[DEBUG] [CONFIG] Environment variable %s found with value: %s", key, value)
		if parsed, err := strconv.ParseBool(value); err == nil {
			log.Printf("[DEBUG] [CONFIG] Successfully parsed %s as boolean: %t", key, parsed)
			return parsed
		}
		log.Printf("[DEBUG] [CONFIG] Failed to parse %s as boolean, using default: %t", key, defaultValue)
	} else {
		log.Printf("[DEBUG] [CONFIG] Environment variable %s not found, using default: %t", key, defaultValue)
	}
	return defaultValue
}
