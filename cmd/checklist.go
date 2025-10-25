package main

import (
	"fmt"

	"github.com/adlio/trello"
	"github.com/danbruder/trello-cli/internal/client"
	"github.com/danbruder/trello-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var checklistListCmd = &cobra.Command{
	Use:   "list --card <card-id>",
	Short: "List all checklists on a card",
	Long:  "List all checklists on a specific card.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID, _ := cmd.Flags().GetString("card")
		if cardID == "" {
			return fmt.Errorf("card ID is required")
		}

		auth := cmd.Context().Value("auth").(*client.AuthConfig)
		trelloClient := client.NewClient(auth.APIKey, auth.Token)

		card, err := trelloClient.GetCard(cardID, trello.Defaults())
		if err != nil {
			return fmt.Errorf("failed to get card: %w", err)
		}

		// Get checklists from card fields
		checklists := card.Checklists
		if checklists == nil {
			checklists = []*trello.Checklist{}
		}

		// Format output
		f, err := formatter.NewFormatter(format, fields, maxTokens, verbose)
		if err != nil {
			return err
		}

		output, err := f.FormatChecklists(checklists)
		if err != nil {
			return err
		}

		if !quiet {
			fmt.Println(output)
		}
		return nil
	},
}

var checklistCreateCmd = &cobra.Command{
	Use:   "create --card <card-id> <name>",
	Short: "Create a new checklist",
	Long:  "Create a new checklist on a specific card.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID, _ := cmd.Flags().GetString("card")
		if cardID == "" {
			return fmt.Errorf("card ID is required")
		}

		auth := cmd.Context().Value("auth").(*client.AuthConfig)
		trelloClient := client.NewClient(auth.APIKey, auth.Token)

		checklistName := args[0]
		card, err := trelloClient.GetCard(cardID, nil)
		if err != nil {
			return fmt.Errorf("failed to get card: %w", err)
		}

		checklist, err := card.CreateChecklist(checklistName, trello.Defaults())
		if err != nil {
			return fmt.Errorf("failed to create checklist: %w", err)
		}

		// Format output
		f, err := formatter.NewFormatter(format, fields, maxTokens, verbose)
		if err != nil {
			return err
		}

		output, err := f.FormatChecklist(checklist)
		if err != nil {
			return err
		}

		if !quiet {
			fmt.Println(output)
		}
		return nil
	},
}

var checklistAddItemCmd = &cobra.Command{
	Use:   "add-item <checklist-id> <name>",
	Short: "Add an item to a checklist",
	Long:  "Add a new item to an existing checklist.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		auth := cmd.Context().Value("auth").(*client.AuthConfig)
		trelloClient := client.NewClient(auth.APIKey, auth.Token)

		checklistID := args[0]
		itemName := args[1]

		checklist, err := trelloClient.GetChecklist(checklistID, trello.Defaults())
		if err != nil {
			return fmt.Errorf("failed to get checklist: %w", err)
		}

		item, err := checklist.CreateCheckItem(itemName, trello.Defaults())
		if err != nil {
			return fmt.Errorf("failed to add item: %w", err)
		}
		_ = item // Suppress unused variable warning

		if !quiet {
			f, _ := formatter.NewFormatter(format, fields, maxTokens, verbose)
			fmt.Println(f.FormatSuccess(fmt.Sprintf("Item '%s' added to checklist '%s'", itemName, checklist.Name)))
		}
		return nil
	},
}

func init() {
	checklistCmd = &cobra.Command{
		Use:   "checklist",
		Short: "Manage Trello checklists",
		Long:  "Commands for managing Trello checklists including listing, creating, and adding items to checklists.",
	}

	checklistCmd.AddCommand(checklistListCmd)
	checklistCmd.AddCommand(checklistCreateCmd)
	checklistCmd.AddCommand(checklistAddItemCmd)

	checklistListCmd.Flags().String("card", "", "Card ID")
	checklistCreateCmd.Flags().String("card", "", "Card ID")
}
