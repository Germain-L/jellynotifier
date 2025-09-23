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
	// Load configuration
	cfg := config.Load()
	log.Printf("Starting JellyNotifier on port %s", cfg.Port)

	var discordBot *discord.Bot
	var err error

	// Initialize Discord bot if enabled and configured
	if cfg.EnableDiscord && cfg.DiscordToken != "" && cfg.DiscordChannel != "" {
		log.Println("Initializing Discord bot...")
		discordBot, err = discord.NewBot(cfg.DiscordToken, cfg.DiscordChannel)
		if err != nil {
			log.Fatalf("Failed to create Discord bot: %v", err)
		}

		if err := discordBot.Start(); err != nil {
			log.Fatalf("Failed to start Discord bot: %v", err)
		}
		log.Println("Discord bot connected successfully")

		defer func() {
			log.Println("Disconnecting Discord bot...")
			if err := discordBot.Stop(); err != nil {
				log.Printf("Error stopping Discord bot: %v", err)
			}
		}()
	} else {
		log.Println("Discord integration disabled or not configured")
		if !cfg.EnableDiscord {
			log.Println("  - Discord disabled by ENABLE_DISCORD=false")
		}
		if cfg.DiscordToken == "" {
			log.Println("  - Missing DISCORD_TOKEN environment variable")
		}
		if cfg.DiscordChannel == "" {
			log.Println("  - Missing DISCORD_CHANNEL_ID environment variable")
		}
	}

	// Initialize webhook handler
	webhookHandler := handlers.NewHandler(discordBot)

	// Set global handler for backward compatibility
	handlers.SetGlobalHandler(webhookHandler)

	// Initialize server
	srv := server.New(cfg.Port, webhookHandler)

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s...", cfg.Port)
		if err := srv.Start(); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := srv.Shutdown(); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	log.Println("Server stopped")
}
