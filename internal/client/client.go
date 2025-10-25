package client

import (
	"github.com/adlio/trello"
)

// Client wraps the Trello client with additional context
type Client struct {
	*trello.Client
	Config *Config
}

// NewClient creates a new Trello client with authentication
func NewClient(apiKey, token string) *Client {
	trelloClient := trello.NewClient(apiKey, token)

	config, _ := LoadConfig()
	if config == nil {
		config = &Config{
			DefaultFormat: "markdown",
			MaxTokens:     4000,
		}
	}

	return &Client{
		Client: trelloClient,
		Config: config,
	}
}
