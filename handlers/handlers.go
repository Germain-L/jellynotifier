package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"jellynotifier/models"
)

// Handler handles incoming webhook notifications
type Handler struct {
	discordBot DiscordBot
}

// DiscordBot interface for Discord bot functionality
type DiscordBot interface {
	SendNotification(notification models.Notification) error
}

// NewHandler creates a new webhook handler with optional Discord bot
func NewHandler(discordBot DiscordBot) *Handler {
	log.Println("[DEBUG] [HANDLERS] Creating new webhook handler...")

	if discordBot != nil {
		log.Println("[DEBUG] [HANDLERS] Discord bot integration enabled")
	} else {
		log.Println("[DEBUG] [HANDLERS] No Discord bot provided - running without Discord integration")
	}

	return &Handler{
		discordBot: discordBot,
	}
}

// HandleWebhook processes incoming webhook notifications
func (h *Handler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DEBUG] [HANDLERS] Incoming request - Method: %s, URL: %s, RemoteAddr: %s", r.Method, r.URL.Path, r.RemoteAddr)
	log.Printf("[DEBUG] [HANDLERS] Request headers - Content-Type: %s, User-Agent: %s", r.Header.Get("Content-Type"), r.Header.Get("User-Agent"))

	defer r.Body.Close() // Ensure request body is closed to prevent resource leaks

	// Only allow POST requests
	if r.Method != http.MethodPost {
		log.Printf("[WARN] [HANDLERS] Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Validate content type
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" && !strings.HasPrefix(contentType, "application/json") {
		log.Printf("[ERROR] [HANDLERS] Invalid content type: %s", contentType)
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}
	log.Println("[DEBUG] [HANDLERS] Content-Type validation passed")

	// Parse the JSON payload
	log.Println("[DEBUG] [HANDLERS] Parsing JSON payload...")
	var notification models.Notification
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		log.Printf("[ERROR] [HANDLERS] Error parsing JSON: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	log.Println("[DEBUG] [HANDLERS] JSON payload parsed successfully")

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

	// Send to Discord if bot is available
	if h.discordBot != nil {
		log.Println("[DEBUG] [HANDLERS] Discord bot available, sending notification...")
		if err := h.discordBot.SendNotification(notification); err != nil {
			log.Printf("[ERROR] [HANDLERS] Error sending notification to Discord: %v", err)
		} else {
			log.Println("[DEBUG] [HANDLERS] Discord notification sent successfully")
		}
	} else {
		log.Println("[DEBUG] [HANDLERS] No Discord bot configured, skipping Discord notification")
	}

	// Send a success response
	log.Println("[DEBUG] [HANDLERS] Sending success response")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Notification received")
	log.Println("[DEBUG] [HANDLERS] Webhook processing completed successfully")
}

// HealthHandler provides a simple health check endpoint
func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DEBUG] [HANDLERS] Health check requested from %s", r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
	log.Println("[DEBUG] [HANDLERS] Health check response sent")
}

// TestHandler provides a test endpoint for development and debugging
func (h *Handler) TestHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DEBUG] [HANDLERS] Test endpoint hit with method: %s from %s", r.Method, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Test successful")
	log.Println("[DEBUG] [HANDLERS] Test response sent")
}

// Legacy function handlers for backward compatibility
var globalHandler *Handler

// SetGlobalHandler sets the global handler instance for legacy functions
func SetGlobalHandler(handler *Handler) {
	globalHandler = handler
}

// WebhookHandler legacy function wrapper
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	if globalHandler != nil {
		globalHandler.HandleWebhook(w, r)
	} else {
		// Fallback to original behavior
		handleWebhookLegacy(w, r)
	}
}

// HealthHandler legacy function wrapper
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if globalHandler != nil {
		globalHandler.HealthHandler(w, r)
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	}
}

// TestHandler legacy function wrapper
func TestHandler(w http.ResponseWriter, r *http.Request) {
	if globalHandler != nil {
		globalHandler.TestHandler(w, r)
	} else {
		log.Printf("Test endpoint hit with method: %s", r.Method)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Test successful")
	}
}

// handleWebhookLegacy provides the original webhook handling without Discord
func handleWebhookLegacy(w http.ResponseWriter, r *http.Request) {
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
