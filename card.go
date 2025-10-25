package main

import (
	"fmt"

	"github.com/adlio/trello"
	"github.com/danbruder/trello-cli/internal/client"
	"github.com/danbruder/trello-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var cardListCmd = &cobra.Command{
	Use:   "list --list <list-id>",
	Short: "List all cards in a list",
	Long:  "List all cards in a specific list.",
	RunE: func(cmd *cobra.Command, args []string) error {
		listID, _ := cmd.Flags().GetString("list")
		if listID == "" {
			return fmt.Errorf("list ID is required")
		}

		auth := cmd.Context().Value("auth").(*client.AuthConfig)
		trelloClient := client.NewClient(auth.APIKey, auth.Token)

		list, err := trelloClient.GetList(listID, nil)
		if err != nil {
			return fmt.Errorf("failed to get list: %w", err)
		}

		cards, err := list.GetCards(nil)
		if err != nil {
			return fmt.Errorf("failed to get cards: %w", err)
		}

		// Format output
		f, err := formatter.NewFormatter(format, fields, maxTokens, verbose)
		if err != nil {
			return err
		}

		output, err := f.FormatCards(cards)
		if err != nil {
			return err
		}

		if !quiet {
			fmt.Println(output)
		}
		return nil
	},
}

var cardGetCmd = &cobra.Command{
	Use:   "get <card-id>",
	Short: "Get card details",
	Long:  "Get detailed information about a specific card.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		auth := cmd.Context().Value("auth").(*client.AuthConfig)
		trelloClient := client.NewClient(auth.APIKey, auth.Token)

		cardID := args[0]
		card, err := trelloClient.GetCard(cardID, nil)
		if err != nil {
			return fmt.Errorf("failed to get card: %w", err)
		}

		// Format output
		f, err := formatter.NewFormatter(format, fields, maxTokens, verbose)
		if err != nil {
			return err
		}

		output, err := f.FormatCard(card)
		if err != nil {
			return err
		}

		if !quiet {
			fmt.Println(output)
		}
		return nil
	},
}

var cardCreateCmd = &cobra.Command{
	Use:   "create --list <list-id> <name>",
	Short: "Create a new card",
	Long:  "Create a new card in a specific list.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		listID, _ := cmd.Flags().GetString("list")
		if listID == "" {
			return fmt.Errorf("list ID is required")
		}

		auth := cmd.Context().Value("auth").(*client.AuthConfig)
		trelloClient := client.NewClient(auth.APIKey, auth.Token)

		cardName := args[0]
		card := trello.Card{
			Name:   cardName,
			IDList: listID,
		}

		err := trelloClient.CreateCard(&card, nil)
		if err != nil {
			return fmt.Errorf("failed to create card: %w", err)
		}

		// Format output
		f, err := formatter.NewFormatter(format, fields, maxTokens, verbose)
		if err != nil {
			return err
		}

		output, err := f.FormatCard(&card)
		if err != nil {
			return err
		}

		if !quiet {
			fmt.Println(output)
		}
		return nil
	},
}

var cardMoveCmd = &cobra.Command{
	Use:   "move <card-id> --list <list-id>",
	Short: "Move a card to another list",
	Long:  "Move a card from its current list to another list.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		listID, _ := cmd.Flags().GetString("list")
		if listID == "" {
			return fmt.Errorf("target list ID is required")
		}

		auth := cmd.Context().Value("auth").(*client.AuthConfig)
		trelloClient := client.NewClient(auth.APIKey, auth.Token)

		cardID := args[0]
		card, err := trelloClient.GetCard(cardID, nil)
		if err != nil {
			return fmt.Errorf("failed to get card: %w", err)
		}

		err = card.MoveToList(listID, nil)
		if err != nil {
			return fmt.Errorf("failed to move card: %w", err)
		}

		if !quiet {
			f, _ := formatter.NewFormatter(format, fields, maxTokens, verbose)
			fmt.Println(f.FormatSuccess(fmt.Sprintf("Card '%s' moved to list %s", card.Name, listID)))
		}
		return nil
	},
}

var cardCopyCmd = &cobra.Command{
	Use:   "copy <card-id> --list <list-id>",
	Short: "Copy a card to another list",
	Long:  "Create a copy of a card in another list.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		listID, _ := cmd.Flags().GetString("list")
		if listID == "" {
			return fmt.Errorf("target list ID is required")
		}

		auth := cmd.Context().Value("auth").(*client.AuthConfig)
		trelloClient := client.NewClient(auth.APIKey, auth.Token)

		cardID := args[0]
		card, err := trelloClient.GetCard(cardID, nil)
		if err != nil {
			return fmt.Errorf("failed to get card: %w", err)
		}

		_, err = card.CopyToList(listID, nil)
		if err != nil {
			return fmt.Errorf("failed to copy card: %w", err)
		}

		if !quiet {
			f, _ := formatter.NewFormatter(format, fields, maxTokens, verbose)
			fmt.Println(f.FormatSuccess(fmt.Sprintf("Card '%s' copied to list %s", card.Name, listID)))
		}
		return nil
	},
}

var cardDeleteCmd = &cobra.Command{
	Use:   "delete <card-id>",
	Short: "Delete a card",
	Long:  "Delete a Trello card permanently.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		auth := cmd.Context().Value("auth").(*client.AuthConfig)
		trelloClient := client.NewClient(auth.APIKey, auth.Token)

		cardID := args[0]
		card, err := trelloClient.GetCard(cardID, nil)
		if err != nil {
			return fmt.Errorf("failed to get card: %w", err)
		}

		err = card.Delete()
		if err != nil {
			return fmt.Errorf("failed to delete card: %w", err)
		}

		if !quiet {
			f, _ := formatter.NewFormatter(format, fields, maxTokens, verbose)
			fmt.Println(f.FormatSuccess(fmt.Sprintf("Card '%s' deleted successfully", card.Name)))
		}
		return nil
	},
}

var cardArchiveCmd = &cobra.Command{
	Use:   "archive <card-id>",
	Short: "Archive a card",
	Long:  "Archive a Trello card.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		auth := cmd.Context().Value("auth").(*client.AuthConfig)
		trelloClient := client.NewClient(auth.APIKey, auth.Token)

		cardID := args[0]
		card, err := trelloClient.GetCard(cardID, nil)
		if err != nil {
			return fmt.Errorf("failed to get card: %w", err)
		}

		err = card.Archive()
		if err != nil {
			return fmt.Errorf("failed to archive card: %w", err)
		}

		if !quiet {
			f, _ := formatter.NewFormatter(format, fields, maxTokens, verbose)
			fmt.Println(f.FormatSuccess(fmt.Sprintf("Card '%s' archived successfully", card.Name)))
		}
		return nil
	},
}

func init() {
	cardCmd = &cobra.Command{
		Use:   "card",
		Short: "Manage Trello cards",
		Long:  "Commands for managing Trello cards including listing, creating, updating, moving, copying, and deleting cards.",
	}

	cardCmd.AddCommand(cardListCmd)
	cardCmd.AddCommand(cardGetCmd)
	cardCmd.AddCommand(cardCreateCmd)
	cardCmd.AddCommand(cardMoveCmd)
	cardCmd.AddCommand(cardCopyCmd)
	cardCmd.AddCommand(cardDeleteCmd)
	cardCmd.AddCommand(cardArchiveCmd)

	cardListCmd.Flags().String("list", "", "List ID")
	cardCreateCmd.Flags().String("list", "", "List ID")
	cardMoveCmd.Flags().String("list", "", "Target list ID")
	cardCopyCmd.Flags().String("list", "", "Target list ID")

}
