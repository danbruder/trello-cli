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

		return executeBatchOperations(batchFile)
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

		return executeBatchOperations(batchFile)
	},
}

func executeBatchOperations(batchFile *batch.BatchFile) error {
	auth := &client.AuthConfig{} // This would be loaded from context in real implementation
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
				err := trelloClient.CreateCard(&card, nil)
				return &card, err
			}
			return nil, fmt.Errorf("list_id is required for create action")
		}
		return nil, fmt.Errorf("card name is required for create action")
	default:
		return nil, fmt.Errorf("unsupported card action: %s", op.Action)
	}
}

func processLabelOperation(trelloClient *client.Client, op batch.Operation) (interface{}, error) {
	// Placeholder implementation
	return nil, fmt.Errorf("label operations not yet implemented")
}

func processChecklistOperation(trelloClient *client.Client, op batch.Operation) (interface{}, error) {
	// Placeholder implementation
	return nil, fmt.Errorf("checklist operations not yet implemented")
}

func processMemberOperation(trelloClient *client.Client, op batch.Operation) (interface{}, error) {
	// Placeholder implementation
	return nil, fmt.Errorf("member operations not yet implemented")
}

func processAttachmentOperation(trelloClient *client.Client, op batch.Operation) (interface{}, error) {
	// Placeholder implementation
	return nil, fmt.Errorf("attachment operations not yet implemented")
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
