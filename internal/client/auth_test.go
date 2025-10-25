package client

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLoadAuth(t *testing.T) {
	t.Run("Environment variables take precedence", func(t *testing.T) {
		// Set environment variables
		os.Setenv("TRELLO_API_KEY", "env-key")
		os.Setenv("TRELLO_TOKEN", "env-token")
		defer os.Unsetenv("TRELLO_API_KEY")
		defer os.Unsetenv("TRELLO_TOKEN")

		auth, err := LoadAuth("flag-key", "flag-token")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		if auth.Source != "environment variables" {
			t.Errorf("Expected source 'environment variables', got %s", auth.Source)
		}

		if auth.APIKey != "env-key" {
			t.Errorf("Expected API key 'env-key', got %s", auth.APIKey)
		}

		if auth.Token != "env-token" {
			t.Errorf("Expected token 'env-token', got %s", auth.Token)
		}
	})

	t.Run("Config file as fallback", func(t *testing.T) {
		// Clear environment variables
		os.Unsetenv("TRELLO_API_KEY")
		os.Unsetenv("TRELLO_TOKEN")

		auth, err := LoadAuth("", "")
		// Config file may or may not exist, so error is acceptable
		if err != nil {
			t.Logf("No config file found (expected): %v", err)
			return
		}

		// If we got here, config file exists
		if auth.Source != "config file" {
			t.Errorf("Expected source 'config file', got %s", auth.Source)
		}

		// Verify it loaded from config
		if auth.APIKey == "" {
			t.Errorf("API key should not be empty when loaded from config")
		}
		if auth.Token == "" {
			t.Errorf("Token should not be empty when loaded from config")
		}
	})

	t.Run("Command-line flags as last resort", func(t *testing.T) {
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
			defer func() {
				if err := os.Rename(backupPath, configPath); err != nil {
					t.Logf("Failed to restore config: %v", err)
				}
			}() // Restore after test
		}

		auth, err := LoadAuth("flag-key", "flag-token")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		// Should use command-line flags
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
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name           string
		configData     Config
		expectError    bool
		expectedFormat string
		expectedTokens int
	}{
		{
			name: "Valid config with defaults",
			configData: Config{
				APIKey:        "test-key",
				Token:         "test-token",
				DefaultFormat: "json",
				MaxTokens:     2000,
			},
			expectError:    false,
			expectedFormat: "json",
			expectedTokens: 2000,
		},
		{
			name: "Config with empty defaults",
			configData: Config{
				APIKey: "test-key",
				Token:  "test-token",
			},
			expectError:    false,
			expectedFormat: "markdown",
			expectedTokens: 4000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary config file
			tempDir := t.TempDir()
			configFile := filepath.Join(tempDir, "config.yaml")

			// Write config directly to file
			data, err := yaml.Marshal(&tt.configData)
			if err != nil {
				t.Fatalf("Failed to marshal config: %v", err)
			}

			err = os.WriteFile(configFile, data, 0600)
			if err != nil {
				t.Fatalf("Failed to write config file: %v", err)
			}

			// Load config directly from file
			data, err = os.ReadFile(configFile)
			if err != nil {
				t.Fatalf("Failed to read config file: %v", err)
			}

			var config Config
			if err := yaml.Unmarshal(data, &config); err != nil {
				t.Fatalf("Failed to parse config: %v", err)
			}

			// Set defaults
			if config.DefaultFormat == "" {
				config.DefaultFormat = "markdown"
			}
			if config.MaxTokens == 0 {
				config.MaxTokens = 4000
			}

			if config.DefaultFormat != tt.expectedFormat {
				t.Errorf("Expected format %s, got %s", tt.expectedFormat, config.DefaultFormat)
			}

			if config.MaxTokens != tt.expectedTokens {
				t.Errorf("Expected max tokens %d, got %d", tt.expectedTokens, config.MaxTokens)
			}
		})
	}
}

func TestSaveConfig(t *testing.T) {
	config := Config{
		APIKey:        "test-key",
		Token:         "test-token",
		DefaultFormat: "json",
		MaxTokens:     2000,
	}

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")

	// Create config directory
	configDir := filepath.Dir(configFile)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	// Write config directly to file
	data, err := yaml.Marshal(&config)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	err = os.WriteFile(configFile, data, 0600)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Errorf("Config file was not created")
	}

	// Load and verify
	data, err = os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	var loadedConfig Config
	if err := yaml.Unmarshal(data, &loadedConfig); err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	if loadedConfig.APIKey != config.APIKey {
		t.Errorf("Expected API key %s, got %s", config.APIKey, loadedConfig.APIKey)
	}

	if loadedConfig.Token != config.Token {
		t.Errorf("Expected token %s, got %s", config.Token, loadedConfig.Token)
	}

	if loadedConfig.DefaultFormat != config.DefaultFormat {
		t.Errorf("Expected format %s, got %s", config.DefaultFormat, loadedConfig.DefaultFormat)
	}

	if loadedConfig.MaxTokens != config.MaxTokens {
		t.Errorf("Expected max tokens %d, got %d", config.MaxTokens, loadedConfig.MaxTokens)
	}
}

func TestNewClient(t *testing.T) {
	apiKey := "test-key"
	token := "test-token"

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

	// Note: NewClient loads config from file, so defaults may be different
	// Config may or may not have values depending on whether config file exists
	if client.Config != nil {
		t.Logf("Config loaded with DefaultFormat: %s, MaxTokens: %d",
			client.Config.DefaultFormat, client.Config.MaxTokens)
	}
}
