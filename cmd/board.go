package cmd

import (
	"fmt"

	"github.com/adlio/trello"
	"github.com/danbruder/trello-cli/internal/client"
	"github.com/danbruder/trello-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var boardListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all boards",
	Long:  "List all boards accessible to the authenticated user.",
	RunE: func(cmd *cobra.Command, args []string) error {
		auth, err := getAuthFromContext(cmd.Context())
		if err != nil {
			return err
		}
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

var boardGetCmd = &cobra.Command{
	Use:   "get <board-id>",
	Short: "Get board details",
	Long:  "Get detailed information about a specific board.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		auth, err := getAuthFromContext(cmd.Context())
		if err != nil {
			return err
		}
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

var boardCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new board",
	Long:  "Create a new Trello board with the specified name.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		auth, err := getAuthFromContext(cmd.Context())
		if err != nil {
			return err
		}
		trelloClient := client.NewClient(auth.APIKey, auth.Token)

		boardName := args[0]
		board := trello.NewBoard(boardName)

		// Set description if provided
		desc, _ := cmd.Flags().GetString("desc")
		if desc != "" {
			board.Desc = desc
		}

		err = trelloClient.CreateBoard(&board, nil)
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

var boardDeleteCmd = &cobra.Command{
	Use:   "delete <board-id>",
	Short: "Delete a board",
	Long:  "Delete a Trello board permanently.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		auth, err := getAuthFromContext(cmd.Context())
		if err != nil {
			return err
		}
		trelloClient := client.NewClient(auth.APIKey, auth.Token)

		boardID := args[0]
		board, err := trelloClient.GetBoard(boardID, nil)
		if err != nil {
			return fmt.Errorf("failed to get board: %w", err)
		}

		err = board.Delete(nil)
		if err != nil {
			return fmt.Errorf("failed to delete board: %w", err)
		}

		if !quiet {
			f, err := formatter.NewFormatter(format, fields, maxTokens, verbose)
			if err != nil {
				return err
			}
			fmt.Println(f.FormatSuccess(fmt.Sprintf("Board '%s' deleted successfully", board.Name)))
		}
		return nil
	},
}

var boardAddMemberCmd = &cobra.Command{
	Use:   "add-member <board-id> <email>",
	Short: "Add a member to a board",
	Long:  "Add a member to a board by email address.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		auth, err := getAuthFromContext(cmd.Context())
		if err != nil {
			return err
		}
		trelloClient := client.NewClient(auth.APIKey, auth.Token)

		boardID := args[0]
		email := args[1]

		board, err := trelloClient.GetBoard(boardID, nil)
		if err != nil {
			return fmt.Errorf("failed to get board: %w", err)
		}

		member := trello.Member{Email: email}
		_, err = board.AddMember(&member, nil)
		if err != nil {
			return fmt.Errorf("failed to add member: %w", err)
		}

		if !quiet {
			f, err := formatter.NewFormatter(format, fields, maxTokens, verbose)
			if err != nil {
				return err
			}
			fmt.Println(f.FormatSuccess(fmt.Sprintf("Member %s added to board '%s'", email, board.Name)))
		}
		return nil
	},
}

func init() {
	boardCmd := &cobra.Command{
		Use:   "board",
		Short: "Manage Trello boards",
		Long:  "Commands for managing Trello boards including listing, creating, updating, and deleting boards.",
	}

	boardCmd.AddCommand(boardListCmd)
	boardCmd.AddCommand(boardGetCmd)
	boardCmd.AddCommand(boardCreateCmd)
	boardCmd.AddCommand(boardDeleteCmd)
	boardCmd.AddCommand(boardAddMemberCmd)

	boardCreateCmd.Flags().String("desc", "", "Board description")

	rootCmd.AddCommand(boardCmd)
}
