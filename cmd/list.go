package main

import (
	"fmt"

	"github.com/adlio/trello"
	"github.com/danbruder/trello-cli/internal/client"
	"github.com/danbruder/trello-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var listListCmd = &cobra.Command{
	Use:   "list --board <board-id>",
	Short: "List all lists on a board",
	Long:  "List all lists on a specific board.",
	RunE: func(cmd *cobra.Command, args []string) error {
		boardID, _ := cmd.Flags().GetString("board")
		if boardID == "" {
			return fmt.Errorf("board ID is required")
		}

		auth := cmd.Context().Value("auth").(*client.AuthConfig)
		trelloClient := client.NewClient(auth.APIKey, auth.Token)

		board, err := trelloClient.GetBoard(boardID, nil)
		if err != nil {
			return fmt.Errorf("failed to get board: %w", err)
		}

		lists, err := board.GetLists(nil)
		if err != nil {
			return fmt.Errorf("failed to get lists: %w", err)
		}

		// Format output
		f, err := formatter.NewFormatter(format, fields, maxTokens, verbose)
		if err != nil {
			return err
		}

		output, err := f.FormatLists(lists)
		if err != nil {
			return err
		}

		if !quiet {
			fmt.Println(output)
		}
		return nil
	},
}

var listGetCmd = &cobra.Command{
	Use:   "get <list-id>",
	Short: "Get list details",
	Long:  "Get detailed information about a specific list.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		auth := cmd.Context().Value("auth").(*client.AuthConfig)
		trelloClient := client.NewClient(auth.APIKey, auth.Token)

		listID := args[0]
		list, err := trelloClient.GetList(listID, nil)
		if err != nil {
			return fmt.Errorf("failed to get list: %w", err)
		}

		// Format output
		f, err := formatter.NewFormatter(format, fields, maxTokens, verbose)
		if err != nil {
			return err
		}

		output, err := f.FormatList(list)
		if err != nil {
			return err
		}

		if !quiet {
			fmt.Println(output)
		}
		return nil
	},
}

var listCreateCmd = &cobra.Command{
	Use:   "create --board <board-id> <name>",
	Short: "Create a new list",
	Long:  "Create a new list on a specific board.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		boardID, _ := cmd.Flags().GetString("board")
		if boardID == "" {
			return fmt.Errorf("board ID is required")
		}

		auth := cmd.Context().Value("auth").(*client.AuthConfig)
		trelloClient := client.NewClient(auth.APIKey, auth.Token)

		listName := args[0]
		
		board, err := trelloClient.GetBoard(boardID, trello.Defaults())
		if err != nil {
			return fmt.Errorf("failed to get board: %w", err)
		}

		list, err := board.CreateList(listName, trello.Defaults())
		if err != nil {
			return fmt.Errorf("failed to create list: %w", err)
		}

		// Format output
		f, err := formatter.NewFormatter(format, fields, maxTokens, verbose)
		if err != nil {
			return err
		}

		output, err := f.FormatList(list)
		if err != nil {
			return err
		}

		if !quiet {
			fmt.Println(output)
		}
		return nil
	},
}

var listArchiveCmd = &cobra.Command{
	Use:   "archive <list-id>",
	Short: "Archive a list",
	Long:  "Archive a Trello list.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		auth := cmd.Context().Value("auth").(*client.AuthConfig)
		trelloClient := client.NewClient(auth.APIKey, auth.Token)

		listID := args[0]
		list, err := trelloClient.GetList(listID, nil)
		if err != nil {
			return fmt.Errorf("failed to get list: %w", err)
		}

		err = list.Archive()
		if err != nil {
			return fmt.Errorf("failed to archive list: %w", err)
		}

		if !quiet {
			f, _ := formatter.NewFormatter(format, fields, maxTokens, verbose)
			fmt.Println(f.FormatSuccess(fmt.Sprintf("List '%s' archived successfully", list.Name)))
		}
		return nil
	},
}

func init() {
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Manage Trello lists",
		Long:  "Commands for managing Trello lists including listing, creating, updating, and archiving lists.",
	}

	listCmd.AddCommand(listListCmd)
	listCmd.AddCommand(listGetCmd)
	listCmd.AddCommand(listCreateCmd)
	listCmd.AddCommand(listArchiveCmd)

	listListCmd.Flags().String("board", "", "Board ID")
	listCreateCmd.Flags().String("board", "", "Board ID")
}
