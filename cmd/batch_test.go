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
