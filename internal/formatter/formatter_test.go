package formatter

import (
	"fmt"
	"strings"
	"testing"

	"github.com/adlio/trello"
)

func TestNewFormatter(t *testing.T) {
	tests := []struct {
		name        string
		format      string
		fields      []string
		maxTokens   int
		verbose     bool
		expectError bool
	}{
		{
			name:        "Valid JSON formatter",
			format:      "json",
			fields:      []string{"id", "name"},
			maxTokens:   1000,
			verbose:     true,
			expectError: false,
		},
		{
			name:        "Valid Markdown formatter",
			format:      "markdown",
			fields:      []string{"id", "name"},
			maxTokens:   1000,
			verbose:     true,
			expectError: false,
		},
		{
			name:        "Valid MD formatter",
			format:      "md",
			fields:      []string{"id", "name"},
			maxTokens:   1000,
			verbose:     true,
			expectError: false,
		},
		{
			name:        "Invalid format",
			format:      "xml",
			fields:      []string{"id", "name"},
			maxTokens:   1000,
			verbose:     true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter, err := NewFormatter(tt.format, tt.fields, tt.maxTokens, tt.verbose)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if formatter == nil {
				t.Errorf("Formatter is nil")
			}
		})
	}
}

func TestJSONFormatter(t *testing.T) {
	formatter := NewJSONFormatter([]string{"id", "name"}, 100, false)

	// Test board formatting
	board := &trello.Board{
		ID:   "test-board-id",
		Name: "Test Board",
		Desc: "Test Description",
	}

	output, err := formatter.FormatBoard(board)
	if err != nil {
		t.Fatalf("Failed to format board: %v", err)
	}

	if !strings.Contains(output, "test-board-id") {
		t.Errorf("Output should contain board ID")
	}

	if !strings.Contains(output, "Test Board") {
		t.Errorf("Output should contain board name")
	}

	// Test error formatting
	errorOutput := formatter.FormatError(fmt.Errorf("test error"))
	if !strings.Contains(errorOutput, "error") {
		t.Errorf("Error output should contain 'error' field")
	}

	// Test success formatting
	successOutput := formatter.FormatSuccess("Test message")
	if !strings.Contains(successOutput, "success") {
		t.Errorf("Success output should contain 'success' field")
	}
}

func TestMarkdownFormatter(t *testing.T) {
	formatter := NewMarkdownFormatter([]string{"id", "name"}, 100, false)

	// Test board formatting
	board := &trello.Board{
		ID:   "test-board-id",
		Name: "Test Board",
		Desc: "Test Description",
	}

	output, err := formatter.FormatBoard(board)
	if err != nil {
		t.Fatalf("Failed to format board: %v", err)
	}

	if !strings.Contains(output, "# Board: Test Board") {
		t.Errorf("Output should contain board title")
	}

	if !strings.Contains(output, "test-board-id") {
		t.Errorf("Output should contain board ID")
	}

	// Test boards formatting
	boards := []*trello.Board{
		{ID: "board1", Name: "Board 1"},
		{ID: "board2", Name: "Board 2"},
	}

	boardsOutput, err := formatter.FormatBoards(boards)
	if err != nil {
		t.Fatalf("Failed to format boards: %v", err)
	}

	if !strings.Contains(boardsOutput, "# Boards (2)") {
		t.Errorf("Output should contain boards count")
	}

	if !strings.Contains(boardsOutput, "Board 1") {
		t.Errorf("Output should contain first board name")
	}

	// Test error formatting
	errorOutput := formatter.FormatError(fmt.Errorf("test error"))
	if !strings.Contains(errorOutput, "❌") {
		t.Errorf("Error output should contain error emoji")
	}

	// Test success formatting
	successOutput := formatter.FormatSuccess("Test message")
	if !strings.Contains(successOutput, "✅") {
		t.Errorf("Success output should contain success emoji")
	}
}

func TestTokenLimiting(t *testing.T) {
	// Test JSON formatter with token limit
	jsonFormatter := NewJSONFormatter([]string{}, 50, false)

	board := &trello.Board{
		ID:   "test-board-id",
		Name: "Test Board",
		Desc: "This is a very long description that should be truncated when the token limit is applied",
	}

	output, err := jsonFormatter.FormatBoard(board)
	if err != nil {
		t.Fatalf("Failed to format board: %v", err)
	}

	// Rough check that output is truncated (50 tokens * 4 chars = 200 chars max)
	if len(output) > 250 { // Allow some buffer
		t.Errorf("Output should be truncated to token limit, got length %d", len(output))
	}

	// Test Markdown formatter with token limit
	mdFormatter := NewMarkdownFormatter([]string{}, 50, false)

	mdOutput, err := mdFormatter.FormatBoard(board)
	if err != nil {
		t.Fatalf("Failed to format board: %v", err)
	}

	if len(mdOutput) > 250 { // Allow some buffer
		t.Errorf("Markdown output should be truncated to token limit, got length %d", len(mdOutput))
	}
}

func TestFieldFiltering(t *testing.T) {
	// Test JSON formatter with field filtering
	jsonFormatter := NewJSONFormatter([]string{"id", "name"}, 0, false)

	board := &trello.Board{
		ID:   "test-board-id",
		Name: "Test Board",
		Desc: "Test Description",
		URL:  "https://trello.com/test",
	}

	output, err := jsonFormatter.FormatBoard(board)
	if err != nil {
		t.Fatalf("Failed to format board: %v", err)
	}

	// Should contain id and name
	if !strings.Contains(output, "test-board-id") {
		t.Errorf("Output should contain board ID")
	}

	if !strings.Contains(output, "Test Board") {
		t.Errorf("Output should contain board name")
	}

	// Should not contain desc and url (filtered out)
	if strings.Contains(output, "Test Description") {
		t.Errorf("Output should not contain description when filtered")
	}

	if strings.Contains(output, "https://trello.com/test") {
		t.Errorf("Output should not contain URL when filtered")
	}
}

func TestVerboseMode(t *testing.T) {
	// Test verbose mode
	verboseFormatter := NewMarkdownFormatter([]string{}, 0, true)
	nonVerboseFormatter := NewMarkdownFormatter([]string{}, 0, false)

	board := &trello.Board{
		ID:     "test-board-id",
		Name:   "Test Board",
		Desc:   "Test Description",
		URL:    "https://trello.com/test",
		Closed: true,
	}

	verboseOutput, err := verboseFormatter.FormatBoard(board)
	if err != nil {
		t.Fatalf("Failed to format board: %v", err)
	}

	nonVerboseOutput, err := nonVerboseFormatter.FormatBoard(board)
	if err != nil {
		t.Fatalf("Failed to format board: %v", err)
	}

	// Verbose output should be longer
	if len(verboseOutput) <= len(nonVerboseOutput) {
		t.Errorf("Verbose output should be longer than non-verbose output")
	}

	// Verbose output should contain more fields
	if !strings.Contains(verboseOutput, "https://trello.com/test") {
		t.Errorf("Verbose output should contain URL")
	}

	if !strings.Contains(verboseOutput, "Test Description") {
		t.Errorf("Verbose output should contain description")
	}

	if !strings.Contains(verboseOutput, "Closed") {
		t.Errorf("Verbose output should contain closed status")
	}
}

func TestCardFormatting(t *testing.T) {
	formatter := NewMarkdownFormatter([]string{}, 0, false)

	card := &trello.Card{
		ID:   "test-card-id",
		Name: "Test Card",
		Desc: "Test Description",
		URL:  "https://trello.com/test",
	}

	output, err := formatter.FormatCard(card)
	if err != nil {
		t.Fatalf("Failed to format card: %v", err)
	}

	if !strings.Contains(output, "# Card: Test Card") {
		t.Errorf("Output should contain card title")
	}

	if !strings.Contains(output, "test-card-id") {
		t.Errorf("Output should contain card ID")
	}
}

func TestListFormatting(t *testing.T) {
	formatter := NewMarkdownFormatter([]string{}, 0, false)

	list := &trello.List{
		ID:     "test-list-id",
		Name:   "Test List",
		Closed: false,
		Pos:    123.45,
	}

	output, err := formatter.FormatList(list)
	if err != nil {
		t.Fatalf("Failed to format list: %v", err)
	}

	if !strings.Contains(output, "# List: Test List") {
		t.Errorf("Output should contain list title")
	}

	if !strings.Contains(output, "test-list-id") {
		t.Errorf("Output should contain list ID")
	}
}
