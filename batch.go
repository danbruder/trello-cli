package main

import (
	"fmt"
	"os"

	"github.com/adlio/trello"
	"github.com/danbruder/trello-cli/internal/batch"
	"github.com/danbruder/trello-cli/internal/client"
	"github.com/spf13/cobra"
)

var batchFileCmd = &cobra.Command{
	Use:   "file <batch-file>",
	Short: "Execute batch operations from a file",
	Long:  "Execute batch operations from a JSON or YAML file.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]

		batchFile, err := batch.LoadBatchFile(filename)
		if err != nil {
			return fmt.Errorf("failed to load batch file: %w", err)
		}

		return executeBatchOperations(cmd, batchFile)
	},
}

var batchStdinCmd = &cobra.Command{
	Use:   "stdin",
	Short: "Execute batch operations from stdin",
	Long:  "Execute batch operations from JSON or YAML piped to stdin.",
	RunE: func(cmd *cobra.Command, args []string) error {
		batchFile, err := batch.LoadBatchFromStdin()
		if err != nil {
			return fmt.Errorf("failed to load batch from stdin: %w", err)
		}

		return executeBatchOperations(cmd, batchFile)
	},
}

func executeBatchOperations(cmd *cobra.Command, batchFile *batch.BatchFile) error {
	auth := cmd.Context().Value("auth").(*client.AuthConfig)
	trelloClient := client.NewClient(auth.APIKey, auth.Token)

	processor := batch.NewBatchProcessor(batchFile.ContinueOnError)

	processor.ProcessOperations(batchFile.Operations, func(op batch.Operation) (interface{}, error) {
		return processOperation(trelloClient, op)
	})

	// Format and output results
	results, err := processor.FormatResults(format)
	if err != nil {
		return err
	}

	if !quiet {
		fmt.Println(results)
	}

	// Exit with appropriate code
	if processor.GetErrorCount() > 0 {
		os.Exit(1)
	}

	return nil
}

func processOperation(trelloClient *client.Client, op batch.Operation) (interface{}, error) {
	// Validate operation
	if err := batch.ValidateOperation(op); err != nil {
		return nil, err
	}

	// Process based on operation type and action
	switch op.Type {
	case "board":
		return processBoardOperation(trelloClient, op)
	case "list":
		return processListOperation(trelloClient, op)
	case "card":
		return processCardOperation(trelloClient, op)
	case "label":
		return processLabelOperation(trelloClient, op)
	case "checklist":
		return processChecklistOperation(trelloClient, op)
	case "member":
		return processMemberOperation(trelloClient, op)
	case "attachment":
		return processAttachmentOperation(trelloClient, op)
	default:
		return nil, fmt.Errorf("unsupported operation type: %s", op.Type)
	}
}

func processBoardOperation(trelloClient *client.Client, op batch.Operation) (interface{}, error) {
	switch op.Action {
	case "get":
		if op.ID == "" {
			return nil, fmt.Errorf("board ID is required for get action")
		}
		return trelloClient.GetBoard(op.ID, nil)
	case "create":
		if name, ok := op.Data["name"].(string); ok {
			board := trello.NewBoard(name)
			err := trelloClient.CreateBoard(&board, nil)
			return &board, err
		}
		return nil, fmt.Errorf("board name is required for create action")
	case "delete":
		if op.ID == "" {
			return nil, fmt.Errorf("board ID is required for delete action")
		}
		board, err := trelloClient.GetBoard(op.ID, nil)
		if err != nil {
			return nil, err
		}
		return nil, board.Delete(nil)
	case "add-member":
		if op.ID == "" {
			return nil, fmt.Errorf("board ID is required for add-member action")
		}
		email, emailOk := op.Data["email"].(string)
		if !emailOk || email == "" {
			return nil, fmt.Errorf("email is required for add-member action")
		}

		board, err := trelloClient.GetBoard(op.ID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get board: %w", err)
		}

		member := trello.Member{Email: email}
		_, err = board.AddMember(&member, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to add member: %w", err)
		}

		return map[string]string{"status": "success", "message": fmt.Sprintf("member %s added to board", email)}, nil
	default:
		return nil, fmt.Errorf("unsupported board action: %s", op.Action)
	}
}

func processListOperation(trelloClient *client.Client, op batch.Operation) (interface{}, error) {
	switch op.Action {
	case "get":
		if op.ID == "" {
			return nil, fmt.Errorf("list ID is required for get action")
		}
		return trelloClient.GetList(op.ID, nil)
	case "create":
		if name, ok := op.Data["name"].(string); ok {
			if boardID, ok := op.Data["board_id"].(string); ok {
				board, err := trelloClient.GetBoard(boardID, trello.Defaults())
				if err != nil {
					return nil, fmt.Errorf("failed to get board: %w", err)
				}
				list, err := board.CreateList(name, trello.Defaults())
				if err != nil {
					return nil, err
				}
				return list, nil
			}
			return nil, fmt.Errorf("board_id is required for create action")
		}
		return nil, fmt.Errorf("list name is required for create action")
	case "archive":
		if op.ID == "" {
			return nil, fmt.Errorf("list ID is required for archive action")
		}
		list, err := trelloClient.GetList(op.ID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get list: %w", err)
		}

		err = list.Archive()
		if err != nil {
			return nil, fmt.Errorf("failed to archive list: %w", err)
		}

		return map[string]string{"status": "success", "message": fmt.Sprintf("list %s archived", op.ID)}, nil
	default:
		return nil, fmt.Errorf("unsupported list action: %s", op.Action)
	}
}

func processCardOperation(trelloClient *client.Client, op batch.Operation) (interface{}, error) {
	switch op.Action {
	case "get":
		if op.ID == "" {
			return nil, fmt.Errorf("card ID is required for get action")
		}
		return trelloClient.GetCard(op.ID, nil)
	case "create":
		if name, ok := op.Data["name"].(string); ok {
			if listID, ok := op.Data["list_id"].(string); ok {
				card := trello.Card{
					Name:   name,
					IDList: listID,
				}
				// Add description if provided
				if desc, ok := op.Data["desc"].(string); ok {
					card.Desc = desc
				}
				// Add position if provided (can be float64 or string "top"/"bottom")
				if pos, ok := op.Data["pos"].(float64); ok {
					card.Pos = pos
				} else if posStr, ok := op.Data["pos"].(string); ok {
					// Handle "top" and "bottom" as special values
					if posStr == "top" {
						card.Pos = 0
					}
					// "bottom" is default, no need to set
				}
				err := trelloClient.CreateCard(&card, nil)
				return &card, err
			}
			return nil, fmt.Errorf("list_id is required for create action")
		}
		return nil, fmt.Errorf("card name is required for create action")
	case "move":
		if op.ID == "" {
			return nil, fmt.Errorf("card ID is required for move action")
		}
		listID, listOk := op.Data["list_id"].(string)
		if !listOk || listID == "" {
			return nil, fmt.Errorf("list_id is required for move action")
		}

		card, err := trelloClient.GetCard(op.ID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get card: %w", err)
		}

		err = card.MoveToList(listID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to move card: %w", err)
		}

		return map[string]string{"status": "success", "message": fmt.Sprintf("card moved to list %s", listID)}, nil
	case "copy":
		if op.ID == "" {
			return nil, fmt.Errorf("card ID is required for copy action")
		}
		listID, listOk := op.Data["list_id"].(string)
		if !listOk || listID == "" {
			return nil, fmt.Errorf("list_id is required for copy action")
		}

		card, err := trelloClient.GetCard(op.ID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get card: %w", err)
		}

		copiedCard, err := card.CopyToList(listID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to copy card: %w", err)
		}

		return copiedCard, nil
	case "delete":
		if op.ID == "" {
			return nil, fmt.Errorf("card ID is required for delete action")
		}

		card, err := trelloClient.GetCard(op.ID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get card: %w", err)
		}

		err = card.Delete()
		if err != nil {
			return nil, fmt.Errorf("failed to delete card: %w", err)
		}

		return map[string]string{"status": "success", "message": "card deleted"}, nil
    
	case "archive":
		if op.ID == "" {
			return nil, fmt.Errorf("card ID is required for archive action")
		}
		card, err := trelloClient.GetCard(op.ID, nil)
		if err != nil {
			return nil, err
		}
		err = card.Archive()
		if err != nil {
			return nil, fmt.Errorf("failed to archive card: %w", err)
		}
		return map[string]string{"status": "success", "message": "card archived"}, nil
    
	default:
		return nil, fmt.Errorf("unsupported card action: %s", op.Action)
	}
}

func processLabelOperation(trelloClient *client.Client, op batch.Operation) (interface{}, error) {
	switch op.Action {
	case "get":
		if op.ID == "" {
			return nil, fmt.Errorf("label ID is required for get action")
		}
		return trelloClient.GetLabel(op.ID, nil)
	case "create":
		name, nameOk := op.Data["name"].(string)
		color, colorOk := op.Data["color"].(string)
		boardID, boardOk := op.Data["board_id"].(string)

		if !nameOk || name == "" {
			return nil, fmt.Errorf("label name is required for create action")
		}
		if !colorOk || color == "" {
			return nil, fmt.Errorf("label color is required for create action")
		}
		if !boardOk || boardID == "" {
			return nil, fmt.Errorf("board_id is required for create action")
		}

		board, err := trelloClient.GetBoard(boardID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get board: %w", err)
		}

		label := trello.Label{
			Name:  name,
			Color: color,
		}

		err = board.CreateLabel(&label, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create label: %w", err)
		}
		return &label, nil
	case "add":
		cardID, cardOk := op.Data["card_id"].(string)
		labelID, labelOk := op.Data["label_id"].(string)

		if !cardOk || cardID == "" {
			return nil, fmt.Errorf("card_id is required for add action")
		}
		if !labelOk || labelID == "" {
			return nil, fmt.Errorf("label_id is required for add action")
		}

		card, err := trelloClient.GetCard(cardID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get card: %w", err)
		}

		err = card.AddIDLabel(labelID)
		if err != nil {
			return nil, fmt.Errorf("failed to add label: %w", err)
		}

		return map[string]string{"status": "success", "message": "label added to card"}, nil
	default:
		return nil, fmt.Errorf("unsupported label action: %s", op.Action)
	}
}

func processChecklistOperation(trelloClient *client.Client, op batch.Operation) (interface{}, error) {
	switch op.Action {
	case "get":
		if op.ID == "" {
			return nil, fmt.Errorf("checklist ID is required for get action")
		}
		return trelloClient.GetChecklist(op.ID, nil)
	case "create":
		name, nameOk := op.Data["name"].(string)
		cardID, cardOk := op.Data["card_id"].(string)

		if !nameOk || name == "" {
			return nil, fmt.Errorf("checklist name is required for create action")
		}
		if !cardOk || cardID == "" {
			return nil, fmt.Errorf("card_id is required for create action")
		}

		card, err := trelloClient.GetCard(cardID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get card: %w", err)
		}

		checklist, err := trelloClient.CreateChecklist(card, name, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create checklist: %w", err)
		}
		return checklist, nil
	case "add-item":
		checklistID, checklistOk := op.Data["checklist_id"].(string)
		itemName, itemOk := op.Data["item_name"].(string)

		if !checklistOk || checklistID == "" {
			return nil, fmt.Errorf("checklist_id is required for add-item action")
		}
		if !itemOk || itemName == "" {
			return nil, fmt.Errorf("item_name is required for add-item action")
		}

		checklist, err := trelloClient.GetChecklist(checklistID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get checklist: %w", err)
		}

		item, err := checklist.CreateCheckItem(itemName, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to add item: %w", err)
		}

		return item, nil
	default:
		return nil, fmt.Errorf("unsupported checklist action: %s", op.Action)
	}
}

func processMemberOperation(trelloClient *client.Client, op batch.Operation) (interface{}, error) {
	switch op.Action {
	case "get":
		if op.ID == "" {
			return nil, fmt.Errorf("member ID or username is required for get action")
		}
		return trelloClient.GetMember(op.ID, nil)
	case "boards":
		if op.ID == "" {
			return nil, fmt.Errorf("member ID or username is required for boards action")
		}
		member, err := trelloClient.GetMember(op.ID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get member: %w", err)
		}

		boards, err := member.GetBoards(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get member boards: %w", err)
		}
		return boards, nil
	default:
		return nil, fmt.Errorf("unsupported member action: %s", op.Action)
	}
}

func processAttachmentOperation(trelloClient *client.Client, op batch.Operation) (interface{}, error) {
	switch op.Action {
	case "list":
		cardID, cardOk := op.Data["card_id"].(string)

		if !cardOk || cardID == "" {
			return nil, fmt.Errorf("card_id is required for list action")
		}

		card, err := trelloClient.GetCard(cardID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get card: %w", err)
		}

		attachments, err := card.GetAttachments(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get attachments: %w", err)
		}
		return attachments, nil
	case "add":
		cardID, cardOk := op.Data["card_id"].(string)
		url, urlOk := op.Data["url"].(string)

		if !cardOk || cardID == "" {
			return nil, fmt.Errorf("card_id is required for add action")
		}
		if !urlOk || url == "" {
			return nil, fmt.Errorf("url is required for add action")
		}

		card, err := trelloClient.GetCard(cardID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get card: %w", err)
		}

		attachment := trello.Attachment{
			URL: url,
		}

		err = card.AddURLAttachment(&attachment)
		if err != nil {
			return nil, fmt.Errorf("failed to add attachment: %w", err)
		}

		return &attachment, nil
	default:
		return nil, fmt.Errorf("unsupported attachment action: %s", op.Action)
	}
}

func init() {
	batchCmd = &cobra.Command{
		Use:   "batch",
		Short: "Execute batch operations",
		Long:  "Execute multiple Trello operations from a file or stdin for automation and scripting.",
	}

	batchCmd.AddCommand(batchFileCmd)
	batchCmd.AddCommand(batchStdinCmd)

}
