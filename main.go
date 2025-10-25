package main

import (
	"context"
	"fmt"
	"os"

	"github.com/adlio/trello"
	"github.com/danbruder/trello-cli/internal/client"
	"github.com/danbruder/trello-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var (
	apiKey    string
	token     string
	format    string
	fields    []string
	maxTokens int
	verbose   bool
	quiet     bool
	debug     bool
)

var rootCmd = &cobra.Command{
	Use:   "trello-cli",
	Short: "A Trello CLI optimized for LLM use",
	Long: `A comprehensive Trello CLI tool built in Go that provides full access to Trello's API
with features optimized for LLM integration including context optimization, batch operations,
and flexible output formats.`,
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
		ctx := context.WithValue(cmd.Context(), "auth", auth)
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
	// Board commands
	boardCmd := &cobra.Command{
		Use:   "board",
		Short: "Manage Trello boards",
		Long:  "Commands for managing Trello boards including listing, creating, updating, and deleting boards.",
	}

	boardListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all boards",
		Long:  "List all boards accessible to the authenticated user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			auth := cmd.Context().Value("auth").(*client.AuthConfig)
			trelloClient := client.NewClient(auth.APIKey, auth.Token)

			// Get current member
			member, err := trelloClient.GetMember("me", nil)
			if err != nil {
				return fmt.Errorf("failed to get current member: %w", err)
			}

			// Get boards
			boards, err := member.GetBoards(nil)
			if err != nil {
				return fmt.Errorf("failed to get boards: %w", err)
			}

			// Format output
			f, err := formatter.NewFormatter(format, fields, maxTokens, verbose)
			if err != nil {
				return err
			}

			output, err := f.FormatBoards(boards)
			if err != nil {
				return err
			}

			if !quiet {
				fmt.Println(output)
			}
			return nil
		},
	}

	boardGetCmd := &cobra.Command{
		Use:   "get <board-id>",
		Short: "Get board details",
		Long:  "Get detailed information about a specific board.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			auth := cmd.Context().Value("auth").(*client.AuthConfig)
			trelloClient := client.NewClient(auth.APIKey, auth.Token)

			boardID := args[0]
			board, err := trelloClient.GetBoard(boardID, nil)
			if err != nil {
				return fmt.Errorf("failed to get board: %w", err)
			}

			// Format output
			f, err := formatter.NewFormatter(format, fields, maxTokens, verbose)
			if err != nil {
				return err
			}

			output, err := f.FormatBoard(board)
			if err != nil {
				return err
			}

			if !quiet {
				fmt.Println(output)
			}
			return nil
		},
	}

	boardCreateCmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new board",
		Long:  "Create a new Trello board with the specified name.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			auth := cmd.Context().Value("auth").(*client.AuthConfig)
			trelloClient := client.NewClient(auth.APIKey, auth.Token)

			boardName := args[0]
			board := trello.NewBoard(boardName)

			err := trelloClient.CreateBoard(&board, nil)
			if err != nil {
				return fmt.Errorf("failed to create board: %w", err)
			}

			// Format output
			f, err := formatter.NewFormatter(format, fields, maxTokens, verbose)
			if err != nil {
				return err
			}

			output, err := f.FormatBoard(&board)
			if err != nil {
				return err
			}

			if !quiet {
				fmt.Println(output)
			}
			return nil
		},
	}

	boardCmd.AddCommand(boardListCmd)
	boardCmd.AddCommand(boardGetCmd)
	boardCmd.AddCommand(boardCreateCmd)
	rootCmd.AddCommand(boardCmd)

	// Config commands
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long:  "Commands for managing Trello CLI configuration including setting credentials and defaults.",
	}

	configShowCmd := &cobra.Command{
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

	configCmd.AddCommand(configShowCmd)
	rootCmd.AddCommand(configCmd)
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

func main() {
	// Initialize all commands
	initCommands()

	if err := rootCmd.Execute(); err != nil {
		if !quiet {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
}
