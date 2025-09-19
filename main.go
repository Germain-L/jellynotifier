package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Database instance
var db *sql.DB

// UserMapping represents a username to Discord ID mapping
type UserMapping struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	DiscordID string `json:"discord_id"`
}

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
	// Initialize database
	if err := initDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Set up the HTTP server
	http.HandleFunc("/webhook", webhookHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/test", testHandler)

	// User mapping endpoints
	http.HandleFunc("/users", usersHandler)      // GET: list all, POST: create new
	http.HandleFunc("/users/", userHandler)      // GET: get by username, PUT: update, DELETE: delete
	http.HandleFunc("/resolve/", resolveHandler) // GET: resolve username to Discord ID

	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Database initialization
func initDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./data/userdb.sqlite")
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Test connection
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	// Create table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS user_mappings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		discord_id TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err = db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

// Database operations
func createUserMapping(username, discordID string) error {
	_, err := db.Exec("INSERT INTO user_mappings (username, discord_id) VALUES (?, ?)", username, discordID)
	return err
}

func getUserMapping(username string) (*UserMapping, error) {
	row := db.QueryRow("SELECT id, username, discord_id FROM user_mappings WHERE username = ?", username)

	var mapping UserMapping
	err := row.Scan(&mapping.ID, &mapping.Username, &mapping.DiscordID)
	if err != nil {
		return nil, err
	}
	return &mapping, nil
}

func getAllUserMappings() ([]UserMapping, error) {
	rows, err := db.Query("SELECT id, username, discord_id FROM user_mappings ORDER BY username")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mappings []UserMapping
	for rows.Next() {
		var mapping UserMapping
		if err := rows.Scan(&mapping.ID, &mapping.Username, &mapping.DiscordID); err != nil {
			return nil, err
		}
		mappings = append(mappings, mapping)
	}
	return mappings, nil
}

func updateUserMapping(username, discordID string) error {
	_, err := db.Exec("UPDATE user_mappings SET discord_id = ? WHERE username = ?", discordID, username)
	return err
}

func deleteUserMapping(username string) error {
	_, err := db.Exec("DELETE FROM user_mappings WHERE username = ?", username)
	return err
}

func resolveDiscordID(username string) string {
	mapping, err := getUserMapping(username)
	if err != nil {
		log.Printf("Failed to resolve Discord ID for username %s: %v", username, err)
		return ""
	}
	return mapping.DiscordID
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

// User management endpoints
func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// List all users
		mappings, err := getAllUserMappings()
		if err != nil {
			http.Error(w, "Failed to get users", http.StatusInternalServerError)
			log.Printf("Error getting users: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mappings)

	case http.MethodPost:
		// Create new user mapping
		var mapping UserMapping
		if err := json.NewDecoder(r.Body).Decode(&mapping); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if mapping.Username == "" || mapping.DiscordID == "" {
			http.Error(w, "Username and discord_id are required", http.StatusBadRequest)
			return
		}

		if err := createUserMapping(mapping.Username, mapping.DiscordID); err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				http.Error(w, "Username already exists", http.StatusConflict)
			} else {
				http.Error(w, "Failed to create user", http.StatusInternalServerError)
				log.Printf("Error creating user: %v", err)
			}
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	// Extract username from URL path
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	username := strings.Split(path, "/")[0]

	if username == "" {
		http.Error(w, "Username required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Get specific user
		mapping, err := getUserMapping(username)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to get user", http.StatusInternalServerError)
				log.Printf("Error getting user: %v", err)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mapping)

	case http.MethodPut:
		// Update user
		var mapping UserMapping
		if err := json.NewDecoder(r.Body).Decode(&mapping); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if mapping.DiscordID == "" {
			http.Error(w, "discord_id is required", http.StatusBadRequest)
			return
		}

		if err := updateUserMapping(username, mapping.DiscordID); err != nil {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			log.Printf("Error updating user: %v", err)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})

	case http.MethodDelete:
		// Delete user
		if err := deleteUserMapping(username); err != nil {
			http.Error(w, "Failed to delete user", http.StatusInternalServerError)
			log.Printf("Error deleting user: %v", err)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func resolveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract username from URL path
	path := strings.TrimPrefix(r.URL.Path, "/resolve/")
	username := strings.Split(path, "/")[0]

	if username == "" {
		http.Error(w, "Username required", http.StatusBadRequest)
		return
	}

	discordID := resolveDiscordID(username)
	if discordID == "" {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"username":   username,
		"discord_id": discordID,
	})
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

	// Log request information if available and resolve Discord ID
	if notification.Request.RequestID != "" {
		log.Printf("  Request ID: %s", notification.Request.RequestID)
		log.Printf("  Requested By: %s (%s)", notification.Request.RequestedByUsername, notification.Request.RequestedByEmail)

		// Try to resolve Discord ID from username if not provided
		if notification.Request.RequestedBySettingsDiscordID == "" && notification.Request.RequestedByUsername != "" {
			if resolvedID := resolveDiscordID(notification.Request.RequestedByUsername); resolvedID != "" {
				log.Printf("  Resolved Discord ID: %s", resolvedID)
				notification.Request.RequestedBySettingsDiscordID = resolvedID
			}
		}
		if notification.Request.RequestedBySettingsDiscordID != "" {
			log.Printf("  Discord ID: %s", notification.Request.RequestedBySettingsDiscordID)
		}
	}

	// Log issue information if available and resolve Discord ID
	if notification.Issue.IssueID != "" {
		log.Printf("  Issue ID: %s", notification.Issue.IssueID)
		log.Printf("  Issue Type: %s", notification.Issue.IssueType)
		log.Printf("  Issue Status: %s", notification.Issue.IssueStatus)
		log.Printf("  Reported By: %s (%s)", notification.Issue.ReportedByUsername, notification.Issue.ReportedByEmail)

		// Try to resolve Discord ID from username if not provided
		if notification.Issue.ReportedBySettingsDiscordID == "" && notification.Issue.ReportedByUsername != "" {
			if resolvedID := resolveDiscordID(notification.Issue.ReportedByUsername); resolvedID != "" {
				log.Printf("  Resolved Discord ID: %s", resolvedID)
				notification.Issue.ReportedBySettingsDiscordID = resolvedID
			}
		}
		if notification.Issue.ReportedBySettingsDiscordID != "" {
			log.Printf("  Discord ID: %s", notification.Issue.ReportedBySettingsDiscordID)
		}
	}

	// Log comment information if available and resolve Discord ID
	if notification.Comment.CommentMessage != "" {
		log.Printf("  Comment: %s", notification.Comment.CommentMessage)
		log.Printf("  Commented By: %s (%s)", notification.Comment.CommentedByUsername, notification.Comment.CommentedByEMail)

		// Try to resolve Discord ID from username if not provided
		if notification.Comment.CommentedBySettingsDiscordID == "" && notification.Comment.CommentedByUsername != "" {
			if resolvedID := resolveDiscordID(notification.Comment.CommentedByUsername); resolvedID != "" {
				log.Printf("  Resolved Discord ID: %s", resolvedID)
				notification.Comment.CommentedBySettingsDiscordID = resolvedID
			}
		}
		if notification.Comment.CommentedBySettingsDiscordID != "" {
			log.Printf("  Discord ID: %s", notification.Comment.CommentedBySettingsDiscordID)
		}
	}

	// Send a success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Notification received")
}
