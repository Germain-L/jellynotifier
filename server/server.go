package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"jellynotifier/handlers"
)

// Server represents the HTTP server configuration
type Server struct {
	Port       string
	httpServer *http.Server
	handler    *handlers.Handler
}

// New creates a new server instance with the provided webhook handler
func New(port string, webhookHandler *handlers.Handler) *Server {
	log.Printf("[DEBUG] [SERVER] Creating new server instance on port %s", port)
	return &Server{
		Port:    port,
		handler: webhookHandler,
	}
}

// SetupRoutes configures all HTTP routes for the server
func (s *Server) SetupRoutes() {
	log.Println("[DEBUG] [SERVER] Setting up HTTP routes...")
	mux := http.NewServeMux()

	if s.handler != nil {
		log.Println("[DEBUG] [SERVER] Using new handler methods")
		mux.HandleFunc("/webhook", s.handler.HandleWebhook)
		mux.HandleFunc("/health", s.handler.HealthHandler)
		mux.HandleFunc("/test", s.handler.TestHandler)
	} else {
		log.Println("[DEBUG] [SERVER] Using legacy handlers for backward compatibility")
		mux.HandleFunc("/webhook", handlers.WebhookHandler)
		mux.HandleFunc("/health", handlers.HealthHandler)
		mux.HandleFunc("/test", handlers.TestHandler)
	}

	s.httpServer = &http.Server{
		Addr:    ":" + s.Port,
		Handler: mux,
	}

	log.Printf("[DEBUG] [SERVER] HTTP server configured on address %s", s.httpServer.Addr)
	log.Println("[DEBUG] [SERVER] Registered routes: /webhook, /health, /test")
}

// Start starts the HTTP server on the configured port
func (s *Server) Start() error {
	log.Printf("[DEBUG] [SERVER] Starting HTTP server setup...")
	s.SetupRoutes()
	log.Printf("[DEBUG] [SERVER] Server setup completed, listening on port %s", s.Port)
	fmt.Printf("Server starting on port %s...\n", s.Port)

	err := s.httpServer.ListenAndServe()
	if err != nil {
		log.Printf("[ERROR] [SERVER] Server failed to start: %v", err)
	}
	return err
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	if s.httpServer == nil {
		log.Println("[DEBUG] [SERVER] No HTTP server to shutdown")
		return nil
	}

	log.Println("[DEBUG] [SERVER] Initiating graceful server shutdown...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // Increased from 5s to 15s
	defer cancel()

	fmt.Println("Shutting down server...")
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		log.Printf("[ERROR] [SERVER] Error during server shutdown: %v", err)
	} else {
		log.Println("[DEBUG] [SERVER] Server shutdown completed successfully")
	}
	return err
}
