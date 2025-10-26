package main

import (
	"context"
	"fmt"
	"os"

	"github.com/danbruder/trello-cli/internal/client"
	"github.com/spf13/cobra"
)

// contextKey is a custom type for context keys to prevent collisions
type contextKey string

const authContextKey contextKey = "auth"

// getAuthFromContext safely retrieves authentication from context
func getAuthFromContext(ctx context.Context) (*client.AuthConfig, error) {
	auth, ok := ctx.Value(authContextKey).(*client.AuthConfig)
	if !ok {
		return nil, fmt.Errorf("authentication not found in context")
	}
	return auth, nil
}

var (
	apiKey    string
	token     string
	debug     bool
	format    string
	fields    []string
	maxTokens int
	verbose   bool
	quiet     bool
)

// Version information set during build
var (
	version   = "1.0.4"
	buildTime = "unknown"
	goVersion = "unknown"
)

// Command variables
var (
	attachmentCmd *cobra.Command
	batchCmd      *cobra.Command
	boardCmd      *cobra.Command
	cardCmd       *cobra.Command
	checklistCmd  *cobra.Command
	configCmd     *cobra.Command
	labelCmd      *cobra.Command
	listCmd       *cobra.Command
	memberCmd     *cobra.Command
)

var rootCmd = &cobra.Command{
	Use:   "trello-cli",
	Short: "A Trello CLI optimized for LLM use",
	Long: `A comprehensive Trello CLI tool built in Go that provides full access to Trello's API
with features optimized for LLM integration including context optimization, batch operations,
and flexible output formats.`,
	Version: version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Load authentication
		auth, err := client.LoadAuth(apiKey, token)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		if debug && !quiet {
			fmt.Printf("Using credentials from: %s\n", auth.Source)
		}

		// Store auth in command context for subcommands
		ctx := context.WithValue(cmd.Context(), authContextKey, auth)
		cmd.SetContext(ctx)
		return nil
	},
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "Trello API key (overrides env/config)")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "Trello token (overrides env/config)")
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "markdown", "Output format (json, markdown)")
	rootCmd.PersistentFlags().StringSliceVar(&fields, "fields", []string{}, "Specific fields to include in output")
	rootCmd.PersistentFlags().IntVar(&maxTokens, "max-tokens", 0, "Maximum tokens in output (0 = unlimited)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Quiet mode (minimal output)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug mode (show API calls)")
}

func initCommands() {
	// Add all commands from the cmd package to the root command
	if boardCmd != nil {
		rootCmd.AddCommand(boardCmd)
	}
	if cardCmd != nil {
		rootCmd.AddCommand(cardCmd)
	}
	if listCmd != nil {
		rootCmd.AddCommand(listCmd)
	}
	if labelCmd != nil {
		rootCmd.AddCommand(labelCmd)
	}
	if checklistCmd != nil {
		rootCmd.AddCommand(checklistCmd)
	}
	if memberCmd != nil {
		rootCmd.AddCommand(memberCmd)
	}
	if attachmentCmd != nil {
		rootCmd.AddCommand(attachmentCmd)
	}
	if batchCmd != nil {
		rootCmd.AddCommand(batchCmd)
	}
	if configCmd != nil {
		rootCmd.AddCommand(configCmd)
	}
}

func main() {
	// Initialize all commands after all init() functions have run
	initCommands()

	if err := rootCmd.Execute(); err != nil {
		if !quiet {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
}
