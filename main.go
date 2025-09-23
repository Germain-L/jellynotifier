package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Define structures that match your JSON payload
type Notification struct {
	NotificationType string `json:"notification_type"`
	Event            string `json:"event"`
	Subject          string `json:"subject"`
	Message          string `json:"message"`
	Image            string `json:"image"`
	Media            struct {
		MediaType string `json:"media_type"`
		TmdbId    string `json:"tmdbId"`
		TvdbId    string `json:"tvdbId"`
		Status    string `json:"status"`
		Status4k  string `json:"status4k"`
	} `json:"{{media}}"`
	Request struct {
		RequestID                         string `json:"request_id"`
		RequestedByEmail                  string `json:"requestedBy_email"`
		RequestedByUsername               string `json:"requestedBy_username"`
		RequestedByAvatar                 string `json:"requestedBy_avatar"`
		RequestedBySettingsDiscordID      string `json:"requestedBy_settings_discordId"`
		RequestedBySettingsTelegramChatID string `json:"requestedBy_settings_telegramChatId"`
	} `json:"{{request}}"`
	Issue struct {
		IssueID                          string `json:"issue_id"`
		IssueType                        string `json:"issue_type"`
		IssueStatus                      string `json:"issue_status"`
		ReportedByEmail                  string `json:"reportedBy_email"`
		ReportedByUsername               string `json:"reportedBy_username"`
		ReportedByAvatar                 string `json:"reportedBy_avatar"`
		ReportedBySettingsDiscordID      string `json:"reportedBy_settings_discordId"`
		ReportedBySettingsTelegramChatID string `json:"reportedBy_settings_telegramChatId"`
	} `json:"{{issue}}"`
	Comment struct {
		CommentMessage                    string `json:"comment_message"`
		CommentedByEMail                  string `json:"commentedBy_email"`
		CommentedByUsername               string `json:"commentedBy_username"`
		CommentedByAvatar                 string `json:"commentedBy_avatar"`
		CommentedBySettingsDiscordID      string `json:"commentedBy_settings_discordId"`
		CommentedBySettingsTelegramChatID string `json:"commentedBy_settings_telegramChatId"`
	} `json:"{{comment}}"`
	Extra []interface{} `json:"{{extra}}"`
}

func main() {
	// Set up the HTTP server
	http.HandleFunc("/webhook", webhookHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/test", testHandler)
	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Simple health check endpoint that accepts GET requests
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	// Test endpoint for webhook testing - accepts both GET and POST
	log.Printf("Test endpoint hit with method: %s", r.Method)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Test successful")
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the JSON payload
	var notification Notification
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
