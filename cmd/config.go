package main

import (
	"fmt"

	"github.com/danbruder/trello-cli/internal/client"
	"github.com/spf13/cobra"
)

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values",
	Long:  "Set configuration values for API credentials and default settings.",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, _ := cmd.Flags().GetString("api-key")
		token, _ := cmd.Flags().GetString("token")
		defaultFormat, _ := cmd.Flags().GetString("default-format")
		maxTokens, _ := cmd.Flags().GetInt("max-tokens")

		config := &client.Config{
			APIKey:        apiKey,
			Token:         token,
			DefaultFormat: defaultFormat,
			MaxTokens:     maxTokens,
		}

		err := client.SaveConfig(config)
		if err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		if !quiet {
			fmt.Println("Configuration saved successfully")
		}
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  "Display the current configuration settings.",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := client.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if !quiet {
			fmt.Printf("Configuration file: %s\n", getConfigPath())
			fmt.Printf("API Key: %s\n", maskString(config.APIKey))
			fmt.Printf("Token: %s\n", maskString(config.Token))
			fmt.Printf("Default Format: %s\n", config.DefaultFormat)
			fmt.Printf("Max Tokens: %d\n", config.MaxTokens)
		}
		return nil
	},
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show configuration file path",
	Long:  "Display the path to the configuration file.",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := client.GetConfigPath()
		if err != nil {
			return fmt.Errorf("failed to get config path: %w", err)
		}

		if !quiet {
			fmt.Println(path)
		}
		return nil
	},
}

func getConfigPath() string {
	path, err := client.GetConfigPath()
	if err != nil {
		return "unknown"
	}
	return path
}

func maskString(s string) string {
	if s == "" {
		return "(not set)"
	}
	if len(s) <= 8 {
		return "***"
	}
	return s[:4] + "***" + s[len(s)-4:]
}

func init() {
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long:  "Commands for managing Trello CLI configuration including setting credentials and defaults.",
	}

	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configPathCmd)

	configSetCmd.Flags().String("api-key", "", "Trello API key")
	configSetCmd.Flags().String("token", "", "Trello token")
	configSetCmd.Flags().String("default-format", "markdown", "Default output format")
	configSetCmd.Flags().Int("max-tokens", 4000, "Default maximum tokens")
}
