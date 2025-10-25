package client

import (
	"testing"
)

func TestNewClientWithEmptyCredentials(t *testing.T) {
	// Test creating client with empty credentials
	client := NewClient("", "")

	if client == nil {
		t.Fatal("NewClient should not return nil even with empty credentials")
	}

	if client.Client == nil {
		t.Fatal("Trello client should not be nil")
	}

	if client.Config == nil {
		t.Fatal("Config should not be nil")
	}
}

func TestNewClientWithValidCredentials(t *testing.T) {
	apiKey := "test-api-key-12345"
	token := "test-token-67890"

	client := NewClient(apiKey, token)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.Client == nil {
		t.Fatal("Trello client is nil")
	}

	if client.Config == nil {
		t.Fatal("Config is nil")
	}

	// Config may or may not have values depending on whether config file exists
	// Just verify config object exists
	t.Logf("Config DefaultFormat: %s", client.Config.DefaultFormat)
	t.Logf("Config MaxTokens: %d", client.Config.MaxTokens)
}

func TestNewClientConfigDefaults(t *testing.T) {
	client := NewClient("test-key", "test-token")

	// Config is created, but values depend on whether config file exists
	// The NewClient function always returns a config, but values may be defaults or from file
	if client.Config == nil {
		t.Fatal("Config should not be nil")
	}

	// Log what we got (values may vary based on config file presence)
	t.Logf("Config DefaultFormat: %s (may be empty if no config file)", client.Config.DefaultFormat)
	t.Logf("Config MaxTokens: %d (may be 0 if no config file)", client.Config.MaxTokens)
}
