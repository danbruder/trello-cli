package main

import (
	"fmt"

	"github.com/adlio/trello"
	"github.com/danbruder/trello-cli/internal/client"
	"github.com/danbruder/trello-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var labelListCmd = &cobra.Command{
	Use:   "list --board <board-id>",
	Short: "List all labels on a board",
	Long:  "List all labels available on a specific board.",
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

		labels, err := board.GetLabels(nil)
		if err != nil {
			return fmt.Errorf("failed to get labels: %w", err)
		}

		// Format output
		f, err := formatter.NewFormatter(format, fields, maxTokens, verbose)
		if err != nil {
			return err
		}

		output, err := f.FormatLabels(labels)
		if err != nil {
			return err
		}

		if !quiet {
			fmt.Println(output)
		}
		return nil
	},
}

var labelCreateCmd = &cobra.Command{
	Use:   "create --board <board-id> --name <name> --color <color>",
	Short: "Create a new label",
	Long:  "Create a new label on a specific board.",
	RunE: func(cmd *cobra.Command, args []string) error {
		boardID, _ := cmd.Flags().GetString("board")
		name, _ := cmd.Flags().GetString("name")
		color, _ := cmd.Flags().GetString("color")

		if boardID == "" {
			return fmt.Errorf("board ID is required")
		}
		if name == "" {
			return fmt.Errorf("label name is required")
		}
		if color == "" {
			return fmt.Errorf("label color is required")
		}

		auth := cmd.Context().Value("auth").(*client.AuthConfig)
		trelloClient := client.NewClient(auth.APIKey, auth.Token)

		board, err := trelloClient.GetBoard(boardID, nil)
		if err != nil {
			return fmt.Errorf("failed to get board: %w", err)
		}

		label := trello.Label{
			Name:  name,
			Color: color,
		}

		err = board.CreateLabel(&label, nil)
		if err != nil {
			return fmt.Errorf("failed to create label: %w", err)
		}

		// Format output
		f, err := formatter.NewFormatter(format, fields, maxTokens, verbose)
		if err != nil {
			return err
		}

		output, err := f.FormatLabel(&label)
		if err != nil {
			return err
		}

		if !quiet {
			fmt.Println(output)
		}
		return nil
	},
}

var labelAddCmd = &cobra.Command{
	Use:   "add <card-id> <label-id>",
	Short: "Add a label to a card",
	Long:  "Add an existing label to a specific card.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		auth := cmd.Context().Value("auth").(*client.AuthConfig)
		trelloClient := client.NewClient(auth.APIKey, auth.Token)

		cardID := args[0]
		labelID := args[1]

		card, err := trelloClient.GetCard(cardID, nil)
		if err != nil {
			return fmt.Errorf("failed to get card: %w", err)
		}

		err = card.AddIDLabel(labelID)
		if err != nil {
			return fmt.Errorf("failed to add label: %w", err)
		}

		if !quiet {
			f, _ := formatter.NewFormatter(format, fields, maxTokens, verbose)
			fmt.Println(f.FormatSuccess(fmt.Sprintf("Label %s added to card '%s'", labelID, card.Name)))
		}
		return nil
	},
}

func init() {
	labelCmd = &cobra.Command{
		Use:   "label",
		Short: "Manage Trello labels",
		Long:  "Commands for managing Trello labels including listing, creating, and adding labels to cards.",
	}

	labelCmd.AddCommand(labelListCmd)
	labelCmd.AddCommand(labelCreateCmd)
	labelCmd.AddCommand(labelAddCmd)

	labelListCmd.Flags().String("board", "", "Board ID")
	labelCreateCmd.Flags().String("board", "", "Board ID")
	labelCreateCmd.Flags().String("name", "", "Label name")
	labelCreateCmd.Flags().String("color", "", "Label color (red, yellow, orange, green, blue, purple, pink, lime, sky, grey)")
}
