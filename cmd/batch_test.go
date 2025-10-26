package cmd

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/danbruder/trello-cli/internal/batch"
	"github.com/danbruder/trello-cli/internal/client"
)

func TestBatchOperationsIntegration(t *testing.T) {
	// Skip if no credentials available
	apiKey := os.Getenv("TRELLO_API_KEY")
	token := os.Getenv("TRELLO_TOKEN")
	if apiKey == "" || token == "" {
		t.Skip("skipping integration test - TRELLO_API_KEY and TRELLO_TOKEN required")
	}

	auth := &client.AuthConfig{
		APIKey: apiKey,
		Token:  token,
	}
	trelloClient := client.NewClient(auth.APIKey, auth.Token)

	t.Run("Label Operations", func(t *testing.T) {
		testLabelOperations(t, trelloClient)
	})

	t.Run("Checklist Operations", func(t *testing.T) {
		testChecklistOperations(t, trelloClient)
	})

	t.Run("Member Operations", func(t *testing.T) {
		testMemberOperations(t, trelloClient)
	})

	t.Run("Attachment Operations", func(t *testing.T) {
		testAttachmentOperations(t, trelloClient)
	})
}

func testLabelOperations(t *testing.T, trelloClient *client.Client) {
	// Test label get operation
	op := batch.Operation{
		Type:     "label",
		Resource: "label",
		Action:   "get",
		ID:       "invalid-label-id", // This should fail
	}

	result, err := processLabelOperation(trelloClient, op)
	if err == nil {
		t.Error("expected error for invalid label ID")
	}
	_ = result // Suppress unused variable warning

	// Test label create operation
	op = batch.Operation{
		Type:     "label",
		Resource: "label",
		Action:   "create",
		Data: map[string]interface{}{
			"board_id": "invalid-board-id", // This should fail
			"name":     "Test Label",
			"color":    "red",
		},
	}

	result, err = processLabelOperation(trelloClient, op)
	if err == nil {
		t.Error("expected error for invalid board ID")
	}
	_ = result

	// Test label add operation
	op = batch.Operation{
		Type:     "label",
		Resource: "label",
		Action:   "add",
		Data: map[string]interface{}{
			"card_id":  "invalid-card-id", // This should fail
			"label_id": "invalid-label-id",
		},
	}

	result, err = processLabelOperation(trelloClient, op)
	if err == nil {
		t.Error("expected error for invalid card ID")
	}
	_ = result

	// Test unsupported action
	op = batch.Operation{
		Type:     "label",
		Resource: "label",
		Action:   "invalid",
	}

	result, err = processLabelOperation(trelloClient, op)
	if err == nil {
		t.Error("expected error for unsupported action")
	}
	_ = result
}

func testChecklistOperations(t *testing.T, trelloClient *client.Client) {
	// Test checklist get operation
	op := batch.Operation{
		Type:     "checklist",
		Resource: "checklist",
		Action:   "get",
		ID:       "invalid-checklist-id", // This should fail
	}

	result, err := processChecklistOperation(trelloClient, op)
	if err == nil {
		t.Error("expected error for invalid checklist ID")
	}
	_ = result

	// Test checklist create operation
	op = batch.Operation{
		Type:     "checklist",
		Resource: "checklist",
		Action:   "create",
		Data: map[string]interface{}{
			"card_id": "invalid-card-id", // This should fail
			"name":    "Test Checklist",
		},
	}

	result, err = processChecklistOperation(trelloClient, op)
	if err == nil {
		t.Error("expected error for invalid card ID")
	}
	_ = result

	// Test checklist add-item operation
	op = batch.Operation{
		Type:     "checklist",
		Resource: "checklist",
		Action:   "add-item",
		Data: map[string]interface{}{
			"checklist_id": "invalid-checklist-id", // This should fail
			"item_name":    "Test Item",
		},
	}

	result, err = processChecklistOperation(trelloClient, op)
	if err == nil {
		t.Error("expected error for invalid checklist ID")
	}
	_ = result

	// Test unsupported action
	op = batch.Operation{
		Type:     "checklist",
		Resource: "checklist",
		Action:   "invalid",
	}

	result, err = processChecklistOperation(trelloClient, op)
	if err == nil {
		t.Error("expected error for unsupported action")
	}
	_ = result
}

func testMemberOperations(t *testing.T, trelloClient *client.Client) {
	// Test member get operation
	op := batch.Operation{
		Type:     "member",
		Resource: "member",
		Action:   "get",
		ID:       "invalid-member-id", // This should fail
	}

	result, err := processMemberOperation(trelloClient, op)
	if err == nil {
		t.Error("expected error for invalid member ID")
	}
	_ = result

	// Test member boards operation
	op = batch.Operation{
		Type:     "member",
		Resource: "member",
		Action:   "boards",
		ID:       "invalid-member-id", // This should fail
	}

	result, err = processMemberOperation(trelloClient, op)
	if err == nil {
		t.Error("expected error for invalid member ID")
	}
	_ = result

	// Test unsupported action
	op = batch.Operation{
		Type:     "member",
		Resource: "member",
		Action:   "invalid",
	}

	result, err = processMemberOperation(trelloClient, op)
	if err == nil {
		t.Error("expected error for unsupported action")
	}
	_ = result
}

func testAttachmentOperations(t *testing.T, trelloClient *client.Client) {
	// Test attachment list operation
	op := batch.Operation{
		Type:     "attachment",
		Resource: "attachment",
		Action:   "list",
		Data: map[string]interface{}{
			"card_id": "invalid-card-id", // This should fail
		},
	}

	result, err := processAttachmentOperation(trelloClient, op)
	if err == nil {
		t.Error("expected error for invalid card ID")
	}
	_ = result

	// Test attachment add operation
	op = batch.Operation{
		Type:     "attachment",
		Resource: "attachment",
		Action:   "add",
		Data: map[string]interface{}{
			"card_id": "invalid-card-id", // This should fail
			"url":     "https://example.com/test.jpg",
		},
	}

	result, err = processAttachmentOperation(trelloClient, op)
	if err == nil {
		t.Error("expected error for invalid card ID")
	}
	_ = result

	// Test unsupported action
	op = batch.Operation{
		Type:     "attachment",
		Resource: "attachment",
		Action:   "invalid",
	}

	result, err = processAttachmentOperation(trelloClient, op)
	if err == nil {
		t.Error("expected error for unsupported action")
	}
	_ = result
}

func TestBatchFileIntegration(t *testing.T) {
	// Skip if no credentials available
	apiKey := os.Getenv("TRELLO_API_KEY")
	token := os.Getenv("TRELLO_TOKEN")
	if apiKey == "" || token == "" {
		t.Skip("skipping integration test - TRELLO_API_KEY and TRELLO_TOKEN required")
	}

	// Test JSON batch file
	jsonBatch := batch.BatchFile{
		Operations: []batch.Operation{
			{
				Type:     "board",
				Resource: "board",
				Action:   "get",
				ID:       "invalid-board-id", // This should fail
			},
			{
				Type:     "card",
				Resource: "card",
				Action:   "create",
				Data: map[string]interface{}{
					"name":    "Test Card",
					"list_id": "invalid-list-id", // This should fail
				},
			},
		},
		ContinueOnError: true,
	}

	// Create temporary batch file
	jsonData, err := json.Marshal(jsonBatch)
	if err != nil {
		t.Fatalf("failed to marshal batch file: %v", err)
	}

	jsonFile := "test_batch_integration.json"
	err = os.WriteFile(jsonFile, jsonData, 0644)
	if err != nil {
		t.Fatalf("failed to create test batch file: %v", err)
	}
	defer os.Remove(jsonFile)

	// Test batch file loading
	loadedBatch, err := batch.LoadBatchFile(jsonFile)
	if err != nil {
		t.Errorf("failed to load batch file: %v", err)
	}

	if len(loadedBatch.Operations) != 2 {
		t.Errorf("expected 2 operations, got %d", len(loadedBatch.Operations))
	}

	if !loadedBatch.ContinueOnError {
		t.Error("expected continue_on_error to be true")
	}

	// Test batch processing with errors
	auth := &client.AuthConfig{
		APIKey: apiKey,
		Token:  token,
	}
	trelloClient := client.NewClient(auth.APIKey, auth.Token)

	processor := batch.NewBatchProcessor(loadedBatch.ContinueOnError)

	processor.ProcessOperations(loadedBatch.Operations, func(op batch.Operation) (interface{}, error) {
		return processOperation(trelloClient, op)
	})

	results := processor.GetResults()
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

	if processor.GetSuccessCount() != 0 {
		t.Errorf("expected 0 successful operations, got %d", processor.GetSuccessCount())
	}

	if processor.GetErrorCount() != 2 {
		t.Errorf("expected 2 errors, got %d", processor.GetErrorCount())
	}

	// Test result formatting
	jsonResult, err := processor.FormatResults("json")
	if err != nil {
		t.Errorf("failed to format results as JSON: %v", err)
	}

	var testData interface{}
	if err := json.Unmarshal([]byte(jsonResult), &testData); err != nil {
		t.Errorf("expected valid JSON output, got error: %v", err)
	}

	markdownResult, err := processor.FormatResults("markdown")
	if err != nil {
		t.Errorf("failed to format results as Markdown: %v", err)
	}

	if len(markdownResult) == 0 {
		t.Error("expected non-empty markdown result")
	}
}

func TestBatchOperationValidation(t *testing.T) {
	tests := []struct {
		name      string
		operation batch.Operation
		expectErr bool
	}{
		{
			name: "valid board operation",
			operation: batch.Operation{
				Type:     "board",
				Resource: "board",
				Action:   "get",
				ID:       "test-id",
			},
			expectErr: false,
		},
		{
			name: "valid label operation",
			operation: batch.Operation{
				Type:     "label",
				Resource: "label",
				Action:   "create",
				Data: map[string]interface{}{
					"board_id": "test-board",
					"name":     "test-label",
					"color":    "red",
				},
			},
			expectErr: false,
		},
		{
			name: "valid checklist operation",
			operation: batch.Operation{
				Type:     "checklist",
				Resource: "checklist",
				Action:   "create",
				Data: map[string]interface{}{
					"card_id": "test-card",
					"name":    "test-checklist",
				},
			},
			expectErr: false,
		},
		{
			name: "valid member operation",
			operation: batch.Operation{
				Type:     "member",
				Resource: "member",
				Action:   "get",
				ID:       "test-member",
			},
			expectErr: false,
		},
		{
			name: "valid attachment operation",
			operation: batch.Operation{
				Type:     "attachment",
				Resource: "attachment",
				Action:   "add",
				Data: map[string]interface{}{
					"card_id": "test-card",
					"url":     "https://example.com/test.jpg",
				},
			},
			expectErr: false,
		},
		{
			name: "invalid operation type",
			operation: batch.Operation{
				Type:     "invalid",
				Resource: "board",
				Action:   "get",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := batch.ValidateOperation(tt.operation)
			if tt.expectErr && err == nil {
				t.Error("expected validation error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected validation error: %v", err)
			}
		})
	}
}

// TestProcessBoardOperationValidation tests board operation validation
func TestProcessBoardOperationValidation(t *testing.T) {
	auth := &client.AuthConfig{
		APIKey: "test-key",
		Token:  "test-token",
	}
	trelloClient := client.NewClient(auth.APIKey, auth.Token)

	tests := []struct {
		name        string
		operation   batch.Operation
		expectError bool
		errorMsg    string
	}{
		{
			name: "Get board without ID",
			operation: batch.Operation{
				Type:     "board",
				Resource: "board",
				Action:   "get",
				ID:       "",
			},
			expectError: true,
			errorMsg:    "board ID is required",
		},
		{
			name: "Create board without name",
			operation: batch.Operation{
				Type:     "board",
				Resource: "board",
				Action:   "create",
				Data:     map[string]interface{}{},
			},
			expectError: true,
			errorMsg:    "board name is required",
		},
		{
			name: "Delete board without ID",
			operation: batch.Operation{
				Type:     "board",
				Resource: "board",
				Action:   "delete",
				ID:       "",
			},
			expectError: true,
			errorMsg:    "board ID is required",
		},
		{
			name: "Add member without board ID",
			operation: batch.Operation{
				Type:     "board",
				Resource: "board",
				Action:   "add-member",
				ID:       "",
				Data: map[string]interface{}{
					"email": "test@example.com",
				},
			},
			expectError: true,
			errorMsg:    "board ID is required",
		},
		{
			name: "Add member without email",
			operation: batch.Operation{
				Type:     "board",
				Resource: "board",
				Action:   "add-member",
				ID:       "test-board-id",
				Data:     map[string]interface{}{},
			},
			expectError: true,
			errorMsg:    "email is required",
		},
		{
			name: "Unsupported board action",
			operation: batch.Operation{
				Type:     "board",
				Resource: "board",
				Action:   "invalid-action",
			},
			expectError: true,
			errorMsg:    "unsupported board action",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := processBoardOperation(trelloClient, tt.operation)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestProcessListOperationValidation tests list operation validation
func TestProcessListOperationValidation(t *testing.T) {
	auth := &client.AuthConfig{
		APIKey: "test-key",
		Token:  "test-token",
	}
	trelloClient := client.NewClient(auth.APIKey, auth.Token)

	tests := []struct {
		name        string
		operation   batch.Operation
		expectError bool
		errorMsg    string
	}{
		{
			name: "Get list without ID",
			operation: batch.Operation{
				Type:     "list",
				Resource: "list",
				Action:   "get",
				ID:       "",
			},
			expectError: true,
			errorMsg:    "list ID is required",
		},
		{
			name: "Create list without name",
			operation: batch.Operation{
				Type:     "list",
				Resource: "list",
				Action:   "create",
				Data: map[string]interface{}{
					"board_id": "test-board",
				},
			},
			expectError: true,
			errorMsg:    "list name is required",
		},
		{
			name: "Create list without board_id",
			operation: batch.Operation{
				Type:     "list",
				Resource: "list",
				Action:   "create",
				Data: map[string]interface{}{
					"name": "Test List",
				},
			},
			expectError: true,
			errorMsg:    "board_id is required",
		},
		{
			name: "Archive list without ID",
			operation: batch.Operation{
				Type:     "list",
				Resource: "list",
				Action:   "archive",
				ID:       "",
			},
			expectError: true,
			errorMsg:    "list ID is required",
		},
		{
			name: "Unsupported list action",
			operation: batch.Operation{
				Type:     "list",
				Resource: "list",
				Action:   "delete",
			},
			expectError: true,
			errorMsg:    "unsupported list action",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := processListOperation(trelloClient, tt.operation)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestProcessCardOperationValidation tests card operation validation
func TestProcessCardOperationValidation(t *testing.T) {
	auth := &client.AuthConfig{
		APIKey: "test-key",
		Token:  "test-token",
	}
	trelloClient := client.NewClient(auth.APIKey, auth.Token)

	tests := []struct {
		name        string
		operation   batch.Operation
		expectError bool
		errorMsg    string
	}{
		{
			name: "Get card without ID",
			operation: batch.Operation{
				Type:     "card",
				Resource: "card",
				Action:   "get",
				ID:       "",
			},
			expectError: true,
			errorMsg:    "card ID is required",
		},
		{
			name: "Create card without name",
			operation: batch.Operation{
				Type:     "card",
				Resource: "card",
				Action:   "create",
				Data: map[string]interface{}{
					"list_id": "test-list",
				},
			},
			expectError: true,
			errorMsg:    "card name is required",
		},
		{
			name: "Create card without list_id",
			operation: batch.Operation{
				Type:     "card",
				Resource: "card",
				Action:   "create",
				Data: map[string]interface{}{
					"name": "Test Card",
				},
			},
			expectError: true,
			errorMsg:    "list_id is required",
		},
		{
			name: "Move card without ID",
			operation: batch.Operation{
				Type:     "card",
				Resource: "card",
				Action:   "move",
				ID:       "",
				Data: map[string]interface{}{
					"list_id": "test-list",
				},
			},
			expectError: true,
			errorMsg:    "card ID is required",
		},
		{
			name: "Move card without list_id",
			operation: batch.Operation{
				Type:     "card",
				Resource: "card",
				Action:   "move",
				ID:       "test-card-id",
				Data:     map[string]interface{}{},
			},
			expectError: true,
			errorMsg:    "list_id is required",
		},
		{
			name: "Copy card without ID",
			operation: batch.Operation{
				Type:     "card",
				Resource: "card",
				Action:   "copy",
				ID:       "",
				Data: map[string]interface{}{
					"list_id": "test-list",
				},
			},
			expectError: true,
			errorMsg:    "card ID is required",
		},
		{
			name: "Copy card without list_id",
			operation: batch.Operation{
				Type:     "card",
				Resource: "card",
				Action:   "copy",
				ID:       "test-card-id",
				Data:     map[string]interface{}{},
			},
			expectError: true,
			errorMsg:    "list_id is required",
		},
		{
			name: "Delete card without ID",
			operation: batch.Operation{
				Type:     "card",
				Resource: "card",
				Action:   "delete",
				ID:       "",
			},
			expectError: true,
			errorMsg:    "card ID is required",
		},
		{
			name: "Archive card without ID",
			operation: batch.Operation{
				Type:     "card",
				Resource: "card",
				Action:   "archive",
				ID:       "",
			},
			expectError: true,
			errorMsg:    "card ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := processCardOperation(trelloClient, tt.operation)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestProcessLabelOperationValidation tests label operation validation
func TestProcessLabelOperationValidation(t *testing.T) {
	auth := &client.AuthConfig{
		APIKey: "test-key",
		Token:  "test-token",
	}
	trelloClient := client.NewClient(auth.APIKey, auth.Token)

	tests := []struct {
		name        string
		operation   batch.Operation
		expectError bool
		errorMsg    string
	}{
		{
			name: "Get label without ID",
			operation: batch.Operation{
				Type:     "label",
				Resource: "label",
				Action:   "get",
				ID:       "",
			},
			expectError: true,
			errorMsg:    "label ID is required",
		},
		{
			name: "Create label without name",
			operation: batch.Operation{
				Type:     "label",
				Resource: "label",
				Action:   "create",
				Data: map[string]interface{}{
					"board_id": "test-board",
					"color":    "red",
				},
			},
			expectError: true,
			errorMsg:    "label name is required",
		},
		{
			name: "Create label without color",
			operation: batch.Operation{
				Type:     "label",
				Resource: "label",
				Action:   "create",
				Data: map[string]interface{}{
					"board_id": "test-board",
					"name":     "Test Label",
				},
			},
			expectError: true,
			errorMsg:    "label color is required",
		},
		{
			name: "Create label without board_id",
			operation: batch.Operation{
				Type:     "label",
				Resource: "label",
				Action:   "create",
				Data: map[string]interface{}{
					"name":  "Test Label",
					"color": "red",
				},
			},
			expectError: true,
			errorMsg:    "board_id is required",
		},
		{
			name: "Add label without card_id",
			operation: batch.Operation{
				Type:     "label",
				Resource: "label",
				Action:   "add",
				Data: map[string]interface{}{
					"label_id": "test-label",
				},
			},
			expectError: true,
			errorMsg:    "card_id is required",
		},
		{
			name: "Add label without label_id",
			operation: batch.Operation{
				Type:     "label",
				Resource: "label",
				Action:   "add",
				Data: map[string]interface{}{
					"card_id": "test-card",
				},
			},
			expectError: true,
			errorMsg:    "label_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := processLabelOperation(trelloClient, tt.operation)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestProcessChecklistOperationValidation tests checklist operation validation
func TestProcessChecklistOperationValidation(t *testing.T) {
	auth := &client.AuthConfig{
		APIKey: "test-key",
		Token:  "test-token",
	}
	trelloClient := client.NewClient(auth.APIKey, auth.Token)

	tests := []struct {
		name        string
		operation   batch.Operation
		expectError bool
		errorMsg    string
	}{
		{
			name: "Get checklist without ID",
			operation: batch.Operation{
				Type:     "checklist",
				Resource: "checklist",
				Action:   "get",
				ID:       "",
			},
			expectError: true,
			errorMsg:    "checklist ID is required",
		},
		{
			name: "Create checklist without name",
			operation: batch.Operation{
				Type:     "checklist",
				Resource: "checklist",
				Action:   "create",
				Data: map[string]interface{}{
					"card_id": "test-card",
				},
			},
			expectError: true,
			errorMsg:    "checklist name is required",
		},
		{
			name: "Create checklist without card_id",
			operation: batch.Operation{
				Type:     "checklist",
				Resource: "checklist",
				Action:   "create",
				Data: map[string]interface{}{
					"name": "Test Checklist",
				},
			},
			expectError: true,
			errorMsg:    "card_id is required",
		},
		{
			name: "Add item without checklist_id",
			operation: batch.Operation{
				Type:     "checklist",
				Resource: "checklist",
				Action:   "add-item",
				Data: map[string]interface{}{
					"item_name": "Test Item",
				},
			},
			expectError: true,
			errorMsg:    "checklist_id is required",
		},
		{
			name: "Add item without item_name",
			operation: batch.Operation{
				Type:     "checklist",
				Resource: "checklist",
				Action:   "add-item",
				Data: map[string]interface{}{
					"checklist_id": "test-checklist",
				},
			},
			expectError: true,
			errorMsg:    "item_name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := processChecklistOperation(trelloClient, tt.operation)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestProcessMemberOperationValidation tests member operation validation
func TestProcessMemberOperationValidation(t *testing.T) {
	auth := &client.AuthConfig{
		APIKey: "test-key",
		Token:  "test-token",
	}
	trelloClient := client.NewClient(auth.APIKey, auth.Token)

	tests := []struct {
		name        string
		operation   batch.Operation
		expectError bool
		errorMsg    string
	}{
		{
			name: "Get member without ID",
			operation: batch.Operation{
				Type:     "member",
				Resource: "member",
				Action:   "get",
				ID:       "",
			},
			expectError: true,
			errorMsg:    "member ID or username is required",
		},
		{
			name: "Get member boards without ID",
			operation: batch.Operation{
				Type:     "member",
				Resource: "member",
				Action:   "boards",
				ID:       "",
			},
			expectError: true,
			errorMsg:    "member ID or username is required",
		},
		{
			name: "Unsupported member action",
			operation: batch.Operation{
				Type:     "member",
				Resource: "member",
				Action:   "delete",
			},
			expectError: true,
			errorMsg:    "unsupported member action",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := processMemberOperation(trelloClient, tt.operation)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestProcessAttachmentOperationValidation tests attachment operation validation
func TestProcessAttachmentOperationValidation(t *testing.T) {
	auth := &client.AuthConfig{
		APIKey: "test-key",
		Token:  "test-token",
	}
	trelloClient := client.NewClient(auth.APIKey, auth.Token)

	tests := []struct {
		name        string
		operation   batch.Operation
		expectError bool
		errorMsg    string
	}{
		{
			name: "List attachments without card_id",
			operation: batch.Operation{
				Type:     "attachment",
				Resource: "attachment",
				Action:   "list",
				Data:     map[string]interface{}{},
			},
			expectError: true,
			errorMsg:    "card_id is required",
		},
		{
			name: "Add attachment without card_id",
			operation: batch.Operation{
				Type:     "attachment",
				Resource: "attachment",
				Action:   "add",
				Data: map[string]interface{}{
					"url": "https://example.com/test.jpg",
				},
			},
			expectError: true,
			errorMsg:    "card_id is required",
		},
		{
			name: "Add attachment without url",
			operation: batch.Operation{
				Type:     "attachment",
				Resource: "attachment",
				Action:   "add",
				Data: map[string]interface{}{
					"card_id": "test-card",
				},
			},
			expectError: true,
			errorMsg:    "url is required",
		},
		{
			name: "Unsupported attachment action",
			operation: batch.Operation{
				Type:     "attachment",
				Resource: "attachment",
				Action:   "delete",
			},
			expectError: true,
			errorMsg:    "unsupported attachment action",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := processAttachmentOperation(trelloClient, tt.operation)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestProcessOperationRouting tests that operations are routed to correct processors
func TestProcessOperationRouting(t *testing.T) {
	auth := &client.AuthConfig{
		APIKey: "test-key",
		Token:  "test-token",
	}
	trelloClient := client.NewClient(auth.APIKey, auth.Token)

	tests := []struct {
		name        string
		operation   batch.Operation
		expectError bool
	}{
		{
			name: "Route to board processor",
			operation: batch.Operation{
				Type:     "board",
				Resource: "board",
				Action:   "get",
				ID:       "test-id",
			},
			expectError: true, // Will fail due to invalid ID, but routing works
		},
		{
			name: "Route to card processor",
			operation: batch.Operation{
				Type:     "card",
				Resource: "card",
				Action:   "get",
				ID:       "test-id",
			},
			expectError: true,
		},
		{
			name: "Route to list processor",
			operation: batch.Operation{
				Type:     "list",
				Resource: "list",
				Action:   "get",
				ID:       "test-id",
			},
			expectError: true,
		},
		{
			name: "Invalid operation type",
			operation: batch.Operation{
				Type:     "invalid-type",
				Resource: "invalid",
				Action:   "get",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := processOperation(trelloClient, tt.operation)

			// We expect errors for all these tests since we're not making real API calls
			// We're just testing that the routing logic works
			if err == nil && tt.expectError {
				t.Errorf("Expected error but got none")
			}
		})
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
