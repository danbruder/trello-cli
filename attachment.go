package main

import (
	"fmt"

	"github.com/adlio/trello"
	"github.com/danbruder/trello-cli/internal/client"
	"github.com/danbruder/trello-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var attachmentListCmd = &cobra.Command{
	Use:   "list --card <card-id>",
	Short: "List all attachments on a card",
	Long:  "List all attachments on a specific card.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID, _ := cmd.Flags().GetString("card")
		if cardID == "" {
			return fmt.Errorf("card ID is required")
		}

		auth := cmd.Context().Value("auth").(*client.AuthConfig)
		trelloClient := client.NewClient(auth.APIKey, auth.Token)

		card, err := trelloClient.GetCard(cardID, nil)
		if err != nil {
			return fmt.Errorf("failed to get card: %w", err)
		}

		attachments, err := card.GetAttachments(nil)
		if err != nil {
			return fmt.Errorf("failed to get attachments: %w", err)
		}

		// Format output
		f, err := formatter.NewFormatter(format, fields, maxTokens, verbose)
		if err != nil {
			return err
		}

		output, err := f.FormatAttachments(attachments)
		if err != nil {
			return err
		}

		if !quiet {
			fmt.Println(output)
		}
		return nil
	},
}

var attachmentAddCmd = &cobra.Command{
	Use:   "add --card <card-id> <url>",
	Short: "Add an attachment to a card",
	Long:  "Add an attachment to a specific card by URL.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID, _ := cmd.Flags().GetString("card")
		if cardID == "" {
			return fmt.Errorf("card ID is required")
		}

		auth := cmd.Context().Value("auth").(*client.AuthConfig)
		trelloClient := client.NewClient(auth.APIKey, auth.Token)

		url := args[0]
		card, err := trelloClient.GetCard(cardID, nil)
		if err != nil {
			return fmt.Errorf("failed to get card: %w", err)
		}

		attachment := trello.Attachment{
			URL: url,
		}

		err = card.AddURLAttachment(&attachment)
		if err != nil {
			return fmt.Errorf("failed to add attachment: %w", err)
		}

		if !quiet {
			f, _ := formatter.NewFormatter(format, fields, maxTokens, verbose)
			fmt.Println(f.FormatSuccess(fmt.Sprintf("Attachment added to card '%s'", card.Name)))
		}
		return nil
	},
}

func init() {
	attachmentCmd = &cobra.Command{
		Use:   "attachment",
		Short: "Manage Trello attachments",
		Long:  "Commands for managing Trello attachments including listing and adding attachments to cards.",
	}

	attachmentCmd.AddCommand(attachmentListCmd)
	attachmentCmd.AddCommand(attachmentAddCmd)

	attachmentListCmd.Flags().String("card", "", "Card ID")
	attachmentAddCmd.Flags().String("card", "", "Card ID")
	
}
