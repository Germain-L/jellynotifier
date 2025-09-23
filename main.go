package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"jellynotifier/config"
	"jellynotifier/discord"
	"jellynotifier/handlers"
	"jellynotifier/server"
)

func main() {
	log.Println("[DEBUG] Starting JellyNotifier application...")

	// Load configuration with validation
	log.Println("[DEBUG] Loading configuration from environment variables...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[ERROR] Configuration error: %v", err)
	}
	log.Printf("[DEBUG] Configuration loaded successfully - Port: %s, EnableDiscord: %t", cfg.Port, cfg.EnableDiscord)

	log.Printf("Starting JellyNotifier on port %s", cfg.Port)

	var discordBot *discord.Bot

	// Initialize Discord bot if enabled and configured
	if cfg.EnableDiscord && cfg.DiscordToken != "" && cfg.DiscordChannel != "" {
		log.Println("[DEBUG] Discord is enabled and configured, initializing bot...")
		log.Printf("[DEBUG] Discord Channel ID: %s", cfg.DiscordChannel)
		log.Printf("[DEBUG] Discord Token length: %d characters", len(cfg.DiscordToken))

		discordBot, err = discord.NewBot(cfg.DiscordToken, cfg.DiscordChannel)
		if err != nil {
			log.Fatalf("[ERROR] Failed to create Discord bot: %v", err)
		}
		log.Println("[DEBUG] Discord bot instance created successfully")

		log.Println("[DEBUG] Starting Discord bot connection...")
		if err := discordBot.Start(); err != nil {
			log.Fatalf("[ERROR] Failed to start Discord bot: %v", err)
		}
		log.Println("Discord bot connected successfully")

		defer func() {
			log.Println("[DEBUG] Graceful shutdown: Disconnecting Discord bot...")
			if err := discordBot.Stop(); err != nil {
				log.Printf("[ERROR] Error stopping Discord bot: %v", err)
			} else {
				log.Println("[DEBUG] Discord bot disconnected successfully")
			}
		}()
	} else {
		log.Println("[DEBUG] Discord integration disabled or not configured")
		if !cfg.EnableDiscord {
			log.Println("[DEBUG] - Discord disabled by ENABLE_DISCORD=false")
		}
		if cfg.DiscordToken == "" {
			log.Println("[DEBUG] - Missing DISCORD_TOKEN environment variable")
		}
		if cfg.DiscordChannel == "" {
			log.Println("[DEBUG] - Missing DISCORD_CHANNEL_ID environment variable")
		}
	}

	// Initialize webhook handler
	log.Println("[DEBUG] Initializing webhook handler...")
	webhookHandler := handlers.NewHandler(discordBot)
	log.Println("[DEBUG] Webhook handler created successfully")

	// Set global handler for backward compatibility
	log.Println("[DEBUG] Setting global handler for backward compatibility...")
	handlers.SetGlobalHandler(webhookHandler)

	// Initialize server
	log.Printf("[DEBUG] Initializing HTTP server on port %s...", cfg.Port)
	srv := server.New(cfg.Port, webhookHandler)
	log.Println("[DEBUG] Server instance created successfully")

	// Start server in a goroutine
	go func() {
		log.Printf("[DEBUG] Starting HTTP server goroutine...")
		log.Printf("Server starting on port %s...", cfg.Port)
		if err := srv.Start(); err != nil {
			log.Fatalf("[ERROR] Server failed to start: %v", err)
		}
	}()

	log.Println("[DEBUG] Server started, waiting for shutdown signal...")
	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	receivedSignal := <-quit
	log.Printf("[DEBUG] Received shutdown signal: %v", receivedSignal)

	log.Println("Shutting down server...")
	if err := srv.Shutdown(); err != nil {
		log.Printf("[ERROR] Error during server shutdown: %v", err)
	} else {
		log.Println("[DEBUG] Server shutdown completed successfully")
	}

	log.Println("Server stopped")
	log.Println("[DEBUG] JellyNotifier application shutdown complete")
}
