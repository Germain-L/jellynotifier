package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"jellynotifier/models"
)

// WebhookHandler handles incoming webhook notifications
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the JSON payload
	var notification models.Notification
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Printf("Error parsing JSON: %v", err)
		return
	}

	// Log the received notification
	log.Printf("Received notification:")
	log.Printf("  Type: %s", notification.NotificationType)
	log.Printf("  Event: %s", notification.Event)
	log.Printf("  Subject: %s", notification.Subject)
	log.Printf("  Message: %s", notification.Message)
	log.Printf("  Image: %s", notification.Image)

	// Log media information if available
	if notification.Media.MediaType != "" {
		log.Printf("  Media Type: %s", notification.Media.MediaType)
		log.Printf("  TMDB ID: %s", notification.Media.TmdbId)
		log.Printf("  TVDB ID: %s", notification.Media.TvdbId)
		log.Printf("  Status: %s", notification.Media.Status)
		log.Printf("  Status 4K: %s", notification.Media.Status4k)
	}

	// Log request information if available
	if notification.Request.RequestID != "" {
		log.Printf("  Request ID: %s", notification.Request.RequestID)
		log.Printf("  Requested By: %s (%s)", notification.Request.RequestedByUsername, notification.Request.RequestedByEmail)
	}

	// Log issue information if available
	if notification.Issue.IssueID != "" {
		log.Printf("  Issue ID: %s", notification.Issue.IssueID)
		log.Printf("  Issue Type: %s", notification.Issue.IssueType)
		log.Printf("  Issue Status: %s", notification.Issue.IssueStatus)
		log.Printf("  Reported By: %s (%s)", notification.Issue.ReportedByUsername, notification.Issue.ReportedByEmail)
	}

	// Log comment information if available
	if notification.Comment.CommentMessage != "" {
		log.Printf("  Comment: %s", notification.Comment.CommentMessage)
		log.Printf("  Commented By: %s (%s)", notification.Comment.CommentedByUsername, notification.Comment.CommentedByEMail)
	}

	// Send a success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Notification received")
}

// HealthHandler provides a simple health check endpoint
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// Simple health check endpoint that accepts GET requests
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

// TestHandler provides a test endpoint for development and debugging
func TestHandler(w http.ResponseWriter, r *http.Request) {
	// Test endpoint for webhook testing - accepts both GET and POST
	log.Printf("Test endpoint hit with method: %s", r.Method)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Test successful")
}
