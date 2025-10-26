package batch

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestValidateOperation(t *testing.T) {
	tests := []struct {
		name        string
		operation   Operation
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid board operation",
			operation: Operation{
				Type:     "board",
				Resource: "board",
				Action:   "get",
				ID:       "test-id",
			},
			expectError: false,
		},
		{
			name: "valid card operation",
			operation: Operation{
				Type:     "card",
				Resource: "card",
				Action:   "create",
				Data:     map[string]interface{}{"name": "test"},
			},
			expectError: false,
		},
		{
			name: "missing type",
			operation: Operation{
				Resource: "board",
				Action:   "get",
			},
			expectError: true,
			errorMsg:    "operation type is required",
		},
		{
			name: "missing resource",
			operation: Operation{
				Type:   "board",
				Action: "get",
			},
			expectError: true,
			errorMsg:    "operation resource is required",
		},
		{
			name: "missing action",
			operation: Operation{
				Type:     "board",
				Resource: "board",
			},
			expectError: true,
			errorMsg:    "operation action is required",
		},
		{
			name: "invalid type",
			operation: Operation{
				Type:     "invalid",
				Resource: "board",
				Action:   "get",
			},
			expectError: true,
			errorMsg:    "invalid operation type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOperation(tt.operation)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestLoadBatchFile(t *testing.T) {
	// Test JSON file
	jsonContent := `{
		"operations": [
			{
				"type": "board",
				"resource": "board",
				"action": "get",
				"id": "test-id"
			}
		],
		"continue_on_error": true
	}`

	jsonFile := "test_batch.json"
	err := os.WriteFile(jsonFile, []byte(jsonContent), 0644)
	if err != nil {
		t.Fatalf("failed to create test JSON file: %v", err)
	}
	defer os.Remove(jsonFile)

	batchFile, err := LoadBatchFile(jsonFile)
	if err != nil {
		t.Errorf("failed to load JSON batch file: %v", err)
	}

	if len(batchFile.Operations) != 1 {
		t.Errorf("expected 1 operation, got %d", len(batchFile.Operations))
	}

	if !batchFile.ContinueOnError {
		t.Error("expected continue_on_error to be true")
	}

	// Test YAML file
	yamlContent := `operations:
  - type: card
    resource: card
    action: create
    data:
      name: "test card"
continue_on_error: false`

	yamlFile := "test_batch.yaml"
	err = os.WriteFile(yamlFile, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("failed to create test YAML file: %v", err)
	}
	defer os.Remove(yamlFile)

	batchFile, err = LoadBatchFile(yamlFile)
	if err != nil {
		t.Errorf("failed to load YAML batch file: %v", err)
	}

	if len(batchFile.Operations) != 1 {
		t.Errorf("expected 1 operation, got %d", len(batchFile.Operations))
	}

	if batchFile.ContinueOnError {
		t.Error("expected continue_on_error to be false")
	}

	// Test invalid file
	invalidFile := "test_invalid.txt"
	err = os.WriteFile(invalidFile, []byte("invalid content"), 0644)
	if err != nil {
		t.Fatalf("failed to create invalid test file: %v", err)
	}
	defer os.Remove(invalidFile)

	_, err = LoadBatchFile(invalidFile)
	if err == nil {
		t.Error("expected error for invalid file format")
	}
}

func TestLoadBatchFromStdin(t *testing.T) {
	// Test JSON input via reader
	t.Run("JSON input", func(t *testing.T) {
		jsonContent := `{
			"operations": [
				{
					"type": "board",
					"resource": "board",
					"action": "get",
					"id": "test-id"
				}
			],
			"continue_on_error": true
		}`

		reader := strings.NewReader(jsonContent)
		batchFile, err := LoadBatchFromReader(reader)
		if err != nil {
			t.Errorf("failed to load batch from JSON reader: %v", err)
		}

		if len(batchFile.Operations) != 1 {
			t.Errorf("expected 1 operation, got %d", len(batchFile.Operations))
		}

		if !batchFile.ContinueOnError {
			t.Error("expected continue_on_error to be true")
		}
	})

	// Test YAML input via reader
	t.Run("YAML input", func(t *testing.T) {
		yamlContent := `operations:
  - type: card
    resource: card
    action: create
    data:
      name: "test card"
continue_on_error: false`

		reader := strings.NewReader(yamlContent)
		batchFile, err := LoadBatchFromReader(reader)
		if err != nil {
			t.Errorf("failed to load batch from YAML reader: %v", err)
		}

		if len(batchFile.Operations) != 1 {
			t.Errorf("expected 1 operation, got %d", len(batchFile.Operations))
		}

		if batchFile.ContinueOnError {
			t.Error("expected continue_on_error to be false")
		}
	})

	// Test invalid input
	t.Run("Invalid input", func(t *testing.T) {
		invalidContent := "this is not valid JSON or YAML"
		reader := strings.NewReader(invalidContent)
		_, err := LoadBatchFromReader(reader)
		if err == nil {
			t.Error("expected error for invalid input")
		}
	})

	// Test empty input (should succeed with empty operations)
	t.Run("Empty input", func(t *testing.T) {
		reader := strings.NewReader("")
		batchFile, err := LoadBatchFromReader(reader)
		if err != nil {
			t.Errorf("unexpected error for empty input: %v", err)
		}
		if len(batchFile.Operations) != 0 {
			t.Errorf("expected 0 operations for empty input, got %d", len(batchFile.Operations))
		}
	})
}

func TestBatchProcessor(t *testing.T) {
	processor := NewBatchProcessor(true)
	if processor == nil {
		t.Fatal("expected processor to be created")
	}

	// Test empty operations
	processor.ProcessOperations([]Operation{}, func(op Operation) (interface{}, error) {
		return "test", nil
	})

	if len(processor.GetResults()) != 0 {
		t.Error("expected no results for empty operations")
	}

	// Test successful operations
	operations := []Operation{
		{
			Type:     "board",
			Resource: "board",
			Action:   "get",
			ID:       "test-id",
		},
		{
			Type:     "card",
			Resource: "card",
			Action:   "create",
			Data:     map[string]interface{}{"name": "test"},
		},
	}

	processor.ProcessOperations(operations, func(op Operation) (interface{}, error) {
		return map[string]string{"id": "test-result"}, nil
	})

	results := processor.GetResults()
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

	if processor.GetSuccessCount() != 2 {
		t.Errorf("expected 2 successful operations, got %d", processor.GetSuccessCount())
	}

	if processor.GetErrorCount() != 0 {
		t.Errorf("expected 0 errors, got %d", processor.GetErrorCount())
	}

	// Test operations with errors
	processor = NewBatchProcessor(true)
	processor.ProcessOperations(operations, func(op Operation) (interface{}, error) {
		return nil, &TestError{Message: "test error"}
	})

	results = processor.GetResults()
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

	if processor.GetSuccessCount() != 0 {
		t.Errorf("expected 0 successful operations, got %d", processor.GetSuccessCount())
	}

	if processor.GetErrorCount() != 2 {
		t.Errorf("expected 2 errors, got %d", processor.GetErrorCount())
	}

	// Test continue on error = false
	processor = NewBatchProcessor(false)
	processor.ProcessOperations(operations, func(op Operation) (interface{}, error) {
		return nil, &TestError{Message: "test error"}
	})

	results = processor.GetResults()
	if len(results) != 1 {
		t.Errorf("expected 1 result (should stop on first error), got %d", len(results))
	}
}

func TestFormatResults(t *testing.T) {
	processor := NewBatchProcessor(true)

	operations := []Operation{
		{
			Type:     "board",
			Resource: "board",
			Action:   "get",
			ID:       "test-id",
		},
	}

	processor.ProcessOperations(operations, func(op Operation) (interface{}, error) {
		return map[string]string{"id": "test-result"}, nil
	})

	// Test JSON formatting
	jsonResult, err := processor.FormatResults("json")
	if err != nil {
		t.Errorf("failed to format JSON results: %v", err)
	}

	var results []BatchResult
	err = json.Unmarshal([]byte(jsonResult), &results)
	if err != nil {
		t.Errorf("failed to unmarshal JSON results: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result in JSON, got %d", len(results))
	}

	// Test Markdown formatting
	markdownResult, err := processor.FormatResults("markdown")
	if err != nil {
		t.Errorf("failed to format Markdown results: %v", err)
	}

	if !strings.Contains(markdownResult, "# Batch Operation Results") {
		t.Error("expected markdown result to contain header")
	}

	if !strings.Contains(markdownResult, "✅") {
		t.Error("expected markdown result to contain success indicator")
	}

	// Test unsupported format
	_, err = processor.FormatResults("xml")
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestFormatResultsWithErrors(t *testing.T) {
	processor := NewBatchProcessor(true)

	operations := []Operation{
		{
			Type:     "board",
			Resource: "board",
			Action:   "get",
			ID:       "test-id",
		},
	}

	processor.ProcessOperations(operations, func(op Operation) (interface{}, error) {
		return nil, &TestError{Message: "test error"}
	})

	// Test JSON formatting with errors
	jsonResult, err := processor.FormatResults("json")
	if err != nil {
		t.Errorf("failed to format JSON results: %v", err)
	}

	var results []BatchResult
	err = json.Unmarshal([]byte(jsonResult), &results)
	if err != nil {
		t.Errorf("failed to unmarshal JSON results: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result in JSON, got %d", len(results))
	}

	if results[0].Success {
		t.Error("expected operation to be marked as failed")
	}

	if results[0].Error != "test error" {
		t.Errorf("expected error message 'test error', got '%s'", results[0].Error)
	}

	// Test Markdown formatting with errors
	markdownResult, err := processor.FormatResults("markdown")
	if err != nil {
		t.Errorf("failed to format Markdown results: %v", err)
	}

	if !strings.Contains(markdownResult, "❌") {
		t.Error("expected markdown result to contain error indicator")
	}

	if !strings.Contains(markdownResult, "test error") {
		t.Error("expected markdown result to contain error message")
	}
}

// TestError is a simple error type for testing
type TestError struct {
	Message string
}

func (e *TestError) Error() string {
	return e.Message
}
