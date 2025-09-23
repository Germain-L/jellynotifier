package server

import (
	"fmt"
	"net/http"

	"jellynotifier/handlers"
)

// Server represents the HTTP server configuration
type Server struct {
	Port string
}

// New creates a new server instance with default configuration
func New() *Server {
	return &Server{
		Port: "8080",
	}
}

// SetupRoutes configures all HTTP routes for the server
func (s *Server) SetupRoutes() {
	http.HandleFunc("/webhook", handlers.WebhookHandler)
	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/test", handlers.TestHandler)
}

// Start starts the HTTP server on the configured port
func (s *Server) Start() error {
	s.SetupRoutes()
	fmt.Printf("Server starting on port %s...\n", s.Port)
	return http.ListenAndServe(":"+s.Port, nil)
}
