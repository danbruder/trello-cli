package client

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Config holds the Trello API credentials and default settings
type Config struct {
	APIKey        string `yaml:"api_key" mapstructure:"api_key"`
	Token         string `yaml:"token" mapstructure:"token"`
	DefaultFormat string `yaml:"default_format" mapstructure:"default_format"`
	MaxTokens     int    `yaml:"max_tokens" mapstructure:"max_tokens"`
}

// AuthConfig holds authentication credentials with their sources
type AuthConfig struct {
	APIKey string
	Token  string
	Source string // "env", "config", or "flags"
}

// LoadAuth loads authentication credentials with precedence order:
// 1. Environment variables
// 2. Config file
// 3. Command-line flags
func LoadAuth(flagAPIKey, flagToken string) (*AuthConfig, error) {
	// First, try environment variables
	envAPIKey := os.Getenv("TRELLO_API_KEY")
	envToken := os.Getenv("TRELLO_TOKEN")

	if envAPIKey != "" && envToken != "" {
		return &AuthConfig{
			APIKey: envAPIKey,
			Token:  envToken,
			Source: "environment variables",
		}, nil
	}

	// Second, try config file
	config, err := LoadConfig()
	if err == nil && config.APIKey != "" && config.Token != "" {
		return &AuthConfig{
			APIKey: config.APIKey,
			Token:  config.Token,
			Source: "config file",
		}, nil
	}

	// Third, try command-line flags
	if flagAPIKey != "" && flagToken != "" {
		return &AuthConfig{
			APIKey: flagAPIKey,
			Token:  flagToken,
			Source: "command-line flags",
		}, nil
	}

	// No valid credentials found
	return nil, fmt.Errorf("no valid Trello credentials found. Please set TRELLO_API_KEY and TRELLO_TOKEN environment variables, create a config file, or use --api-key and --token flags")
}

// LoadConfig loads configuration from the config file
func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(home, ".trello-cli")
	configFile := filepath.Join(configDir, "config.yaml")

	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return &Config{}, nil // Return empty config, not an error
	}

	// Read config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults
	if config.DefaultFormat == "" {
		config.DefaultFormat = "json"
	}
	if config.MaxTokens == 0 {
		config.MaxTokens = 4000
	}

	return &config, nil
}

// SaveConfig saves the configuration to the config file
func SaveConfig(config *Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, ".trello-cli")
	configFile := filepath.Join(configDir, "config.yaml")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write config file
	if err := os.WriteFile(configFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".trello-cli", "config.yaml"), nil
}

// InitViper initializes viper for configuration management
func InitViper() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	viper.AddConfigPath(filepath.Join(home, ".trello-cli"))
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Set defaults
	viper.SetDefault("default_format", "json")
	viper.SetDefault("max_tokens", 4000)

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	return nil
}
