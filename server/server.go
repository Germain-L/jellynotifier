package server

import (
	"context"
	"fmt"
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
	return &Server{
		Port:    port,
		handler: webhookHandler,
	}
}

// NewLegacy creates a new server instance with legacy configuration (backward compatibility)
func NewLegacy() *Server {
	return &Server{
		Port: "8080",
	}
}

// SetupRoutes configures all HTTP routes for the server
func (s *Server) SetupRoutes() {
	mux := http.NewServeMux()

	if s.handler != nil {
		// Use new handler methods
		mux.HandleFunc("/webhook", s.handler.HandleWebhook)
		mux.HandleFunc("/health", s.handler.HealthHandler)
		mux.HandleFunc("/test", s.handler.TestHandler)
	} else {
		// Use legacy handlers for backward compatibility
		mux.HandleFunc("/webhook", handlers.WebhookHandler)
		mux.HandleFunc("/health", handlers.HealthHandler)
		mux.HandleFunc("/test", handlers.TestHandler)
	}

	s.httpServer = &http.Server{
		Addr:    ":" + s.Port,
		Handler: mux,
	}
}

// Start starts the HTTP server on the configured port
func (s *Server) Start() error {
	s.SetupRoutes()
	fmt.Printf("Server starting on port %s...\n", s.Port)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	if s.httpServer == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}
