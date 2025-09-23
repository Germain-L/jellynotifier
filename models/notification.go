package models

// Notification represents the webhook payload structure with template-style field names
type Notification struct {
	NotificationType string        `json:"notification_type"`
	Event            string        `json:"event"`
	Subject          string        `json:"subject"`
	Message          string        `json:"message"`
	Image            string        `json:"image"`
	Media            Media         `json:"{{media}}"`
	Request          Request       `json:"{{request}}"`
	Issue            Issue         `json:"{{issue}}"`
	Comment          Comment       `json:"{{comment}}"`
	Extra            []interface{} `json:"{{extra}}"`
}

// Media contains media-related information from the notification
type Media struct {
	MediaType string `json:"media_type"`
	TmdbId    string `json:"tmdbId"`
	TvdbId    string `json:"tvdbId"`
	Status    string `json:"status"`
	Status4k  string `json:"status4k"`
}

// Request contains request-related information from the notification
type Request struct {
	RequestID                         string `json:"request_id"`
	RequestedByEmail                  string `json:"requestedBy_email"`
	RequestedByUsername               string `json:"requestedBy_username"`
	RequestedByAvatar                 string `json:"requestedBy_avatar"`
	RequestedBySettingsDiscordID      string `json:"requestedBy_settings_discordId"`
	RequestedBySettingsTelegramChatID string `json:"requestedBy_settings_telegramChatId"`
}

// Issue contains issue-related information from the notification
type Issue struct {
	IssueID                          string `json:"issue_id"`
	IssueType                        string `json:"issue_type"`
	IssueStatus                      string `json:"issue_status"`
	ReportedByEmail                  string `json:"reportedBy_email"`
	ReportedByUsername               string `json:"reportedBy_username"`
	ReportedByAvatar                 string `json:"reportedBy_avatar"`
	ReportedBySettingsDiscordID      string `json:"reportedBy_settings_discordId"`
	ReportedBySettingsTelegramChatID string `json:"reportedBy_settings_telegramChatId"`
}

// Comment contains comment-related information from the notification
type Comment struct {
	CommentMessage                    string `json:"comment_message"`
	CommentedByEMail                  string `json:"commentedBy_email"`
	CommentedByUsername               string `json:"commentedBy_username"`
	CommentedByAvatar                 string `json:"commentedBy_avatar"`
	CommentedBySettingsDiscordID      string `json:"commentedBy_settings_discordId"`
	CommentedBySettingsTelegramChatID string `json:"commentedBy_settings_telegramChatId"`
}
