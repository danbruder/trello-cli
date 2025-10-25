package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/danbruder/trello-cli/internal/client"
	"github.com/danbruder/trello-cli/internal/context"
	"github.com/spf13/cobra"
)

func TestBoardCommands(t *testing.T) {
	// Skip integration tests that require Trello credentials
	t.Skip("Skipping integration tests - requires Trello API credentials")
}

func TestConfigCommands(t *testing.T) {
	t.Run("Config Show Command", func(t *testing.T) {
		configShowCmd := &cobra.Command{
			Use: "show",
			RunE: func(cmd *cobra.Command, args []string) error {
				config, err := client.LoadConfig()
				if err != nil {
					return err
				}

				// Verify config has default values
				if config.DefaultFormat == "" {
					t.Errorf("Default format should not be empty")
				}

				if config.MaxTokens == 0 {
					t.Errorf("Max tokens should not be zero")
				}

				return nil
			},
		}

		err := configShowCmd.RunE(configShowCmd, []string{})
		if err != nil {
			t.Errorf("Config show command failed: %v", err)
		}
	})
}

func TestAuthenticationFlow(t *testing.T) {
	t.Run("Environment Variable Authentication", func(t *testing.T) {
		// Set test environment variables
		os.Setenv("TRELLO_API_KEY", "test-key")
		os.Setenv("TRELLO_TOKEN", "test-token")
		defer os.Unsetenv("TRELLO_API_KEY")
		defer os.Unsetenv("TRELLO_TOKEN")

		auth, err := client.LoadAuth("", "")
		if err != nil {
			t.Errorf("Failed to load auth from environment: %v", err)
		}

		if auth.Source != "environment variables" {
			t.Errorf("Expected source 'environment variables', got %s", auth.Source)
		}

		if auth.APIKey != "test-key" {
			t.Errorf("Expected API key 'test-key', got %s", auth.APIKey)
		}

		if auth.Token != "test-token" {
			t.Errorf("Expected token 'test-token', got %s", auth.Token)
		}
	})

	t.Run("Command Line Flag Authentication", func(t *testing.T) {
		// Clear environment variables
		os.Unsetenv("TRELLO_API_KEY")
		os.Unsetenv("TRELLO_TOKEN")

		// Temporarily rename the config file to test flags-only scenario
		home, err := os.UserHomeDir()
		if err != nil {
			t.Fatalf("Failed to get home directory: %v", err)
		}
		configPath := filepath.Join(home, ".trello-cli", "config.yaml")
		backupPath := configPath + ".backup"

		// Backup existing config
		if _, err := os.Stat(configPath); err == nil {
			err = os.Rename(configPath, backupPath)
			if err != nil {
				t.Fatalf("Failed to backup config: %v", err)
			}
			defer os.Rename(backupPath, configPath) // Restore after test
		}

		auth, err := client.LoadAuth("flag-key", "flag-token")
		if err != nil {
			t.Errorf("Failed to load auth from flags: %v", err)
		}

		if auth.Source != "command-line flags" {
			t.Errorf("Expected source 'command-line flags', got %s", auth.Source)
		}

		if auth.APIKey != "flag-key" {
			t.Errorf("Expected API key 'flag-key', got %s", auth.APIKey)
		}

		if auth.Token != "flag-token" {
			t.Errorf("Expected token 'flag-token', got %s", auth.Token)
		}
	})

	t.Run("No Authentication", func(t *testing.T) {
		// Clear environment variables
		os.Unsetenv("TRELLO_API_KEY")
		os.Unsetenv("TRELLO_TOKEN")

		// Temporarily rename the config file to test no-auth scenario
		home, err := os.UserHomeDir()
		if err != nil {
			t.Fatalf("Failed to get home directory: %v", err)
		}
		configPath := filepath.Join(home, ".trello-cli", "config.yaml")
		backupPath := configPath + ".backup"

		// Backup existing config
		if _, err := os.Stat(configPath); err == nil {
			err = os.Rename(configPath, backupPath)
			if err != nil {
				t.Fatalf("Failed to backup config: %v", err)
			}
			defer os.Rename(backupPath, configPath) // Restore after test
		}

		_, err = client.LoadAuth("", "")
		if err == nil {
			t.Errorf("Expected error when no authentication provided")
		}
	})
}

func TestFormatterIntegration(t *testing.T) {
	// Skip integration tests that require Trello credentials
	t.Skip("Skipping integration tests - requires Trello API credentials")
}

func TestContextOptimization(t *testing.T) {
	t.Run("Token Limiting", func(t *testing.T) {
		optimizer := context.NewOptimizer(20, []string{}, false) // 20 tokens = 80 chars

		longText := "This is a very long text that should be truncated when the token limit is applied because it exceeds the maximum allowed tokens and should result in truncation"

		result := optimizer.TruncateToTokenLimit(longText)

		if len(result) >= len(longText) {
			t.Errorf("Text should be truncated")
		}

		if !strings.Contains(result, "... (truncated to fit token limit)") {
			t.Errorf("Should contain truncation message")
		}
	})

	t.Run("Field Filtering", func(t *testing.T) {
		optimizer := context.NewOptimizer(0, []string{"id", "name"}, false)

		if !optimizer.ShouldIncludeField("id") {
			t.Errorf("Should include 'id' field")
		}

		if !optimizer.ShouldIncludeField("name") {
			t.Errorf("Should include 'name' field")
		}

		if optimizer.ShouldIncludeField("url") {
			t.Errorf("Should not include 'url' field")
		}
	})
}
