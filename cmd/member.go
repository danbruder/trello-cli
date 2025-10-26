package cmd

import (
	"fmt"

	"github.com/danbruder/trello-cli/internal/client"
	"github.com/danbruder/trello-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var memberGetCmd = &cobra.Command{
	Use:   "get <username-or-id>",
	Short: "Get member information",
	Long:  "Get detailed information about a specific member.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		auth, err := getAuthFromContext(cmd.Context())
		if err != nil {
			return err
		}
		trelloClient := client.NewClient(auth.APIKey, auth.Token)

		usernameOrID := args[0]
		member, err := trelloClient.GetMember(usernameOrID, nil)
		if err != nil {
			return fmt.Errorf("failed to get member: %w", err)
		}

		// Format output
		f, err := formatter.NewFormatter(format, fields, maxTokens, verbose)
		if err != nil {
			return err
		}

		output, err := f.FormatMember(member)
		if err != nil {
			return err
		}

		if !quiet {
			fmt.Println(output)
		}
		return nil
	},
}

var memberBoardsCmd = &cobra.Command{
	Use:   "boards <username-or-id>",
	Short: "List member's boards",
	Long:  "List all boards that a specific member has access to.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		auth, err := getAuthFromContext(cmd.Context())
		if err != nil {
			return err
		}
		trelloClient := client.NewClient(auth.APIKey, auth.Token)

		usernameOrID := args[0]
		member, err := trelloClient.GetMember(usernameOrID, nil)
		if err != nil {
			return fmt.Errorf("failed to get member: %w", err)
		}

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

func init() {
	memberCmd := &cobra.Command{
		Use:   "member",
		Short: "Manage Trello members",
		Long:  "Commands for managing Trello members including getting member information and listing boards.",
	}

	memberCmd.AddCommand(memberGetCmd)
	memberCmd.AddCommand(memberBoardsCmd)

	rootCmd.AddCommand(memberCmd)
}
