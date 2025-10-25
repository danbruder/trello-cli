package batch

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Operation represents a single batch operation
type Operation struct {
	Type       string                 `json:"type" yaml:"type"`
	Resource   string                 `json:"resource" yaml:"resource"`
	Action     string                 `json:"action" yaml:"action"`
	ID         string                 `json:"id,omitempty" yaml:"id,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty" yaml:"data,omitempty"`
	Parameters map[string]string      `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// BatchFile represents a batch operations file
type BatchFile struct {
	Operations      []Operation `json:"operations" yaml:"operations"`
	ContinueOnError bool        `json:"continue_on_error" yaml:"continue_on_error"`
}

// BatchProcessor handles batch operations
type BatchProcessor struct {
	continueOnError bool
	results         []BatchResult
}

// BatchResult represents the result of a batch operation
type BatchResult struct {
	Operation Operation   `json:"operation"`
	Success   bool        `json:"success"`
	Error     string      `json:"error,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor(continueOnError bool) *BatchProcessor {
	return &BatchProcessor{
		continueOnError: continueOnError,
		results:         make([]BatchResult, 0),
	}
}

// LoadBatchFile loads batch operations from a file
func LoadBatchFile(filename string) (*BatchFile, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read batch file: %w", err)
	}

	var batchFile BatchFile

	// Try JSON first
	if err := json.Unmarshal(data, &batchFile); err == nil {
		return &batchFile, nil
	}

	// Try YAML
	if err := yaml.Unmarshal(data, &batchFile); err == nil {
		return &batchFile, nil
	}

	return nil, fmt.Errorf("failed to parse batch file as JSON or YAML")
}

// LoadBatchFromStdin loads batch operations from stdin
func LoadBatchFromStdin() (*BatchFile, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("failed to read from stdin: %w", err)
	}

	var batchFile BatchFile

	// Try JSON first
	if err := json.Unmarshal(data, &batchFile); err == nil {
		return &batchFile, nil
	}

	// Try YAML
	if err := yaml.Unmarshal(data, &batchFile); err == nil {
		return &batchFile, nil
	}

	return nil, fmt.Errorf("failed to parse stdin as JSON or YAML")
}

// ProcessOperations processes a list of operations
func (bp *BatchProcessor) ProcessOperations(operations []Operation, processor func(Operation) (interface{}, error)) {
	for _, op := range operations {
		result := BatchResult{
			Operation: op,
		}

		data, err := processor(op)
		if err != nil {
			result.Success = false
			result.Error = err.Error()

			if !bp.continueOnError {
				bp.results = append(bp.results, result)
				return
			}
		} else {
			result.Success = true
			result.Data = data
		}

		bp.results = append(bp.results, result)
	}
}

// GetResults returns the processing results
func (bp *BatchProcessor) GetResults() []BatchResult {
	return bp.results
}

// GetSuccessCount returns the number of successful operations
func (bp *BatchProcessor) GetSuccessCount() int {
	count := 0
	for _, result := range bp.results {
		if result.Success {
			count++
		}
	}
	return count
}

// GetErrorCount returns the number of failed operations
func (bp *BatchProcessor) GetErrorCount() int {
	count := 0
	for _, result := range bp.results {
		if !result.Success {
			count++
		}
	}
	return count
}

// FormatResults formats the batch results for output
func (bp *BatchProcessor) FormatResults(format string) (string, error) {
	switch strings.ToLower(format) {
	case "json":
		data, err := json.MarshalIndent(bp.results, "", "  ")
		if err != nil {
			return "", err
		}
		return string(data), nil

	case "markdown", "md":
		return bp.formatMarkdown(), nil

	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

// formatMarkdown formats results as markdown
func (bp *BatchProcessor) formatMarkdown() string {
	var sb strings.Builder

	sb.WriteString("# Batch Operation Results\n\n")
	sb.WriteString(fmt.Sprintf("**Total Operations:** %d\n", len(bp.results)))
	sb.WriteString(fmt.Sprintf("**Successful:** %d\n", bp.GetSuccessCount()))
	sb.WriteString(fmt.Sprintf("**Failed:** %d\n\n", bp.GetErrorCount()))

	for i, result := range bp.results {
		status := "✅"
		if !result.Success {
			status = "❌"
		}

		sb.WriteString(fmt.Sprintf("## Operation %d %s\n", i+1, status))
		sb.WriteString(fmt.Sprintf("- **Type:** %s\n", result.Operation.Type))
		sb.WriteString(fmt.Sprintf("- **Resource:** %s\n", result.Operation.Resource))
		sb.WriteString(fmt.Sprintf("- **Action:** %s\n", result.Operation.Action))

		if result.Operation.ID != "" {
			sb.WriteString(fmt.Sprintf("- **ID:** %s\n", result.Operation.ID))
		}

		if result.Success {
			sb.WriteString("- **Status:** Success\n")
		} else {
			sb.WriteString("- **Status:** Failed\n")
			sb.WriteString(fmt.Sprintf("- **Error:** %s\n", result.Error))
		}

		sb.WriteString("\n")
	}

	return sb.String()
}

// ValidateOperation validates a batch operation
func ValidateOperation(op Operation) error {
	if op.Type == "" {
		return fmt.Errorf("operation type is required")
	}

	if op.Resource == "" {
		return fmt.Errorf("operation resource is required")
	}

	if op.Action == "" {
		return fmt.Errorf("operation action is required")
	}

	// Validate operation types
	validTypes := []string{"board", "list", "card", "label", "checklist", "member", "attachment"}
	validType := false
	for _, vt := range validTypes {
		if op.Type == vt {
			validType = true
			break
		}
	}

	if !validType {
		return fmt.Errorf("invalid operation type: %s (valid: %s)", op.Type, strings.Join(validTypes, ", "))
	}

	return nil
}
