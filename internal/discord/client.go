package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Client handles Discord API interactions
type Client struct {
	// Discord webhook URL
	webhookURL string
}

// Message represents a Discord message
type Message struct {
	Content   string   `json:"content,omitempty"`
	Username  string   `json:"username,omitempty"`
	AvatarURL string   `json:"avatar_url,omitempty"`
	Embeds    []*Embed `json:"embeds,omitempty"`
}

// Embed represents a Discord embed
type Embed struct {
	Title       string      `json:"title,omitempty"`
	Description string      `json:"description,omitempty"`
	Color       int         `json:"color,omitempty"`
	Fields      []EmbedField `json:"fields,omitempty"`
}

// EmbedField represents a field in a Discord embed
type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// NewClient creates a new Discord client
func NewClient(webhookURL string) (*Client, error) {
	if webhookURL == "" {
		return nil, errors.New("discord webhook URL is required")
	}

	// Check if webhookURL is actually a webhook URL
	if !strings.HasPrefix(webhookURL, "https://discord.com/api/webhooks/") {
		return nil, errors.New("invalid webhook URL format, must start with https://discord.com/api/webhooks/")
	}

	client := &Client{
		webhookURL: webhookURL,
	}

	return client, nil
}

// SendMessage sends a message to Discord via webhook
func (c *Client) SendMessage(content string) error {
	return c.sendWebhookMessage(content)
}

// sendWebhookMessage sends a message via webhook
func (c *Client) sendWebhookMessage(content string) error {
	message := Message{
		Content:  content,
		Username: "Delivr",
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	resp, err := http.Post(c.webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err == nil {
			return fmt.Errorf("error sending message to Discord: HTTP %d %s, %v", 
				resp.StatusCode, resp.Status, response)
		}
		return fmt.Errorf("error sending message to Discord: HTTP %d %s", 
			resp.StatusCode, resp.Status)
	}

	return nil
}

// SendEmbed sends a rich embed message to Discord
func (c *Client) SendEmbed(title, description string, fields []EmbedField, color int) error {
	embed := &Embed{
		Title:       title,
		Description: description,
		Color:       color,
		Fields:      fields,
	}

	message := Message{
		Username: "Delivr",
		Embeds:   []*Embed{embed},
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	resp, err := http.Post(c.webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err == nil {
			return fmt.Errorf("error sending embed to Discord: HTTP %d %s, %v", 
				resp.StatusCode, resp.Status, response)
		}
		return fmt.Errorf("error sending embed to Discord: HTTP %d %s", 
			resp.StatusCode, resp.Status)
	}

	return nil
}

// Close is a no-op for webhook clients
func (c *Client) Close() error {
	return nil
}
