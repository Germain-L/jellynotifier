package discord

import (
	"fmt"
	"log"
	"strings"

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
	if token == "" {
		return nil, fmt.Errorf("discord token is required")
	}
	if channelID == "" {
		return nil, fmt.Errorf("discord channel ID is required")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %v", err)
	}

	return &Bot{
		session:   dg,
		channelID: channelID,
	}, nil
}

// Start opens the Discord connection
func (b *Bot) Start() error {
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("error opening Discord connection: %v", err)
	}
	log.Println("Discord bot connected successfully")
	return nil
}

// Stop closes the Discord connection
func (b *Bot) Stop() error {
	return b.session.Close()
}

// SendNotification sends a formatted notification to the Discord channel
func (b *Bot) SendNotification(notification models.Notification) error {
	embed := b.createEmbed(notification)

	_, err := b.session.ChannelMessageSendEmbed(b.channelID, embed)
	if err != nil {
		return fmt.Errorf("error sending message to Discord: %v", err)
	}

	log.Printf("Successfully sent notification to Discord channel %s", b.channelID)
	return nil
}

// createEmbed creates a Discord embed from the notification
func (b *Bot) createEmbed(notification models.Notification) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:       notification.Subject,
		Description: notification.Message,
		Color:       b.getColorForEvent(notification.Event),
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: notification.Image},
		Fields:      []*discordgo.MessageEmbedField{},
	}

	// Add notification type and event info
	if notification.NotificationType != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "📋 Type",
			Value:  notification.NotificationType,
			Inline: true,
		})
	}

	if notification.Event != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "🎬 Event",
			Value:  notification.Event,
			Inline: true,
		})
	}

	// Add media information
	if notification.Media.MediaType != "" {
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
		}
	}

	// Add request information
	if notification.Request.RequestID != "" {
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
	}

	// Add issue information
	if notification.Issue.IssueID != "" {
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
	}

	// Add comment information
	if notification.Comment.CommentMessage != "" {
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
	}

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
