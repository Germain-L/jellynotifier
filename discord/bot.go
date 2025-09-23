package discord

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"jellynotifier/models"
)

// Bot represents the Discord bot instance
type Bot struct {
	session   *discordgo.Session
	channelID string
}

// NewBot creates a new Discord bot instance
func NewBot(token, channelID string) (*Bot, error) {
	log.Println("[DEBUG] [DISCORD] Creating new Discord bot instance...")

	if token == "" {
		log.Println("[ERROR] [DISCORD] Discord token is empty")
		return nil, fmt.Errorf("discord token is required")
	}
	if channelID == "" {
		log.Println("[ERROR] [DISCORD] Discord channel ID is empty")
		return nil, fmt.Errorf("discord channel ID is required")
	}

	log.Printf("[DEBUG] [DISCORD] Token length: %d characters", len(token))
	log.Printf("[DEBUG] [DISCORD] Target channel ID: %s", channelID)

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Printf("[ERROR] [DISCORD] Failed to create Discord session: %v", err)
		return nil, fmt.Errorf("error creating Discord session: %v", err)
	}

	log.Println("[DEBUG] [DISCORD] Discord session created successfully")
	return &Bot{
		session:   dg,
		channelID: channelID,
	}, nil
}

// Start opens the Discord connection and waits for it to be ready
func (b *Bot) Start() error {
	log.Println("[DEBUG] [DISCORD] Opening Discord connection...")

	if err := b.session.Open(); err != nil {
		log.Printf("[ERROR] [DISCORD] Failed to open connection: %v", err)
		return fmt.Errorf("error opening Discord connection: %v", err)
	}

	// Wait for the connection to be ready
	log.Println("[DEBUG] [DISCORD] Waiting for Discord bot to be ready...")
	time.Sleep(2 * time.Second)

	log.Println("[DEBUG] [DISCORD] Connection established and ready")
	log.Println("Discord bot connected successfully")
	return nil
}

// Stop closes the Discord connection gracefully
func (b *Bot) Stop() error {
	log.Println("[DEBUG] [DISCORD] Attempting to close Discord connection...")

	if b.session != nil && b.session.DataReady {
		log.Println("[DEBUG] [DISCORD] Session is active, closing...")
		if err := b.session.Close(); err != nil {
			log.Printf("[ERROR] [DISCORD] Error closing session: %v", err)
			return err
		}
		log.Println("[DEBUG] [DISCORD] Session closed successfully")
	} else {
		log.Println("[DEBUG] [DISCORD] Session is not active or not ready, skipping close")
	}
	return nil
}

// SendNotification sends a formatted notification to the Discord channel
func (b *Bot) SendNotification(notification models.Notification) error {
	log.Printf("[DEBUG] [DISCORD] Preparing to send notification - Type: %s, Event: %s", notification.NotificationType, notification.Event)

	embed := b.createEmbed(notification)
	log.Printf("[DEBUG] [DISCORD] Created embed with %d fields", len(embed.Fields))

	_, err := b.session.ChannelMessageSendEmbed(b.channelID, embed)
	if err != nil {
		log.Printf("[ERROR] [DISCORD] Failed to send message: %v", err)
		return fmt.Errorf("error sending message to Discord: %v", err)
	}

	log.Printf("[DEBUG] [DISCORD] Message sent successfully to channel %s", b.channelID)
	log.Printf("Successfully sent notification to Discord channel %s", b.channelID)
	return nil
}

// createEmbed creates a Discord embed from the notification
func (b *Bot) createEmbed(notification models.Notification) *discordgo.MessageEmbed {
	log.Println("[DEBUG] [DISCORD] Creating Discord embed...")

	embed := &discordgo.MessageEmbed{
		Title:       notification.Subject,
		Description: notification.Message,
		Color:       b.getColorForEvent(notification.Event),
		Timestamp:   time.Now().Format(time.RFC3339),
		Fields:      []*discordgo.MessageEmbedField{},
	}

	log.Printf("[DEBUG] [DISCORD] Embed base created - Title: %s, Color: %d", embed.Title, embed.Color)

	// Add thumbnail only if image URL is provided and not empty
	if notification.Image != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: notification.Image}
		log.Printf("[DEBUG] [DISCORD] Added thumbnail: %s", notification.Image)
	}

	// Add notification type and event info
	if notification.NotificationType != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "📋 Type",
			Value:  notification.NotificationType,
			Inline: true,
		})
		log.Printf("[DEBUG] [DISCORD] Added notification type field: %s", notification.NotificationType)
	}

	if notification.Event != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "🎬 Event",
			Value:  notification.Event,
			Inline: true,
		})
		log.Printf("[DEBUG] [DISCORD] Added event field: %s", notification.Event)
	}

	// Add media information
	if notification.Media.MediaType != "" {
		log.Println("[DEBUG] [DISCORD] Processing media information...")
		mediaInfo := []string{}
		if notification.Media.MediaType != "" {
			mediaInfo = append(mediaInfo, fmt.Sprintf("Type: %s", notification.Media.MediaType))
		}
		if notification.Media.Status != "" {
			mediaInfo = append(mediaInfo, fmt.Sprintf("Status: %s", notification.Media.Status))
		}
		if notification.Media.Status4k != "" {
			mediaInfo = append(mediaInfo, fmt.Sprintf("4K Status: %s", notification.Media.Status4k))
		}
		if notification.Media.TmdbId != "" {
			mediaInfo = append(mediaInfo, fmt.Sprintf("TMDB: %s", notification.Media.TmdbId))
		}

		if len(mediaInfo) > 0 {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "🎭 Media Info",
				Value:  strings.Join(mediaInfo, "\n"),
				Inline: false,
			})
			log.Printf("[DEBUG] [DISCORD] Added media info field with %d items", len(mediaInfo))
		}
	}

	// Add request information
	if notification.Request.RequestID != "" {
		log.Println("[DEBUG] [DISCORD] Processing request information...")
		requestInfo := []string{
			fmt.Sprintf("ID: %s", notification.Request.RequestID),
		}
		if notification.Request.RequestedByUsername != "" {
			requestInfo = append(requestInfo, fmt.Sprintf("Requested by: %s", notification.Request.RequestedByUsername))
		}
		if notification.Request.RequestedByEmail != "" {
			requestInfo = append(requestInfo, fmt.Sprintf("Email: %s", notification.Request.RequestedByEmail))
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "📝 Request Info",
			Value:  strings.Join(requestInfo, "\n"),
			Inline: false,
		})
		log.Printf("[DEBUG] [DISCORD] Added request info field for request ID: %s", notification.Request.RequestID)
	}

	// Add issue information
	if notification.Issue.IssueID != "" {
		log.Println("[DEBUG] [DISCORD] Processing issue information...")
		issueInfo := []string{
			fmt.Sprintf("ID: %s", notification.Issue.IssueID),
		}
		if notification.Issue.IssueType != "" {
			issueInfo = append(issueInfo, fmt.Sprintf("Type: %s", notification.Issue.IssueType))
		}
		if notification.Issue.IssueStatus != "" {
			issueInfo = append(issueInfo, fmt.Sprintf("Status: %s", notification.Issue.IssueStatus))
		}
		if notification.Issue.ReportedByUsername != "" {
			issueInfo = append(issueInfo, fmt.Sprintf("Reported by: %s", notification.Issue.ReportedByUsername))
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "🐛 Issue Info",
			Value:  strings.Join(issueInfo, "\n"),
			Inline: false,
		})
		log.Printf("[DEBUG] [DISCORD] Added issue info field for issue ID: %s", notification.Issue.IssueID)
	}

	// Add comment information
	if notification.Comment.CommentMessage != "" {
		log.Println("[DEBUG] [DISCORD] Processing comment information...")
		commentInfo := []string{
			fmt.Sprintf("Message: %s", notification.Comment.CommentMessage),
		}
		if notification.Comment.CommentedByUsername != "" {
			commentInfo = append(commentInfo, fmt.Sprintf("By: %s", notification.Comment.CommentedByUsername))
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "💬 Comment",
			Value:  strings.Join(commentInfo, "\n"),
			Inline: false,
		})
		log.Printf("[DEBUG] [DISCORD] Added comment field from user: %s", notification.Comment.CommentedByUsername)
	}

	log.Printf("[DEBUG] [DISCORD] Embed creation completed with %d total fields", len(embed.Fields))
	return embed
}

// getColorForEvent returns an appropriate color for the notification event
func (b *Bot) getColorForEvent(event string) int {
	switch strings.ToLower(event) {
	case "media.available":
		return 0x00FF00 // Green
	case "media.requested":
		return 0x0099FF // Blue
	case "media.approved":
		return 0x00FF99 // Teal
	case "media.declined":
		return 0xFF0000 // Red
	case "issue.created":
		return 0xFF6600 // Orange
	case "issue.resolved":
		return 0x00FF00 // Green
	case "comment.created":
		return 0x9900FF // Purple
	default:
		return 0x999999 // Gray
	}
}
