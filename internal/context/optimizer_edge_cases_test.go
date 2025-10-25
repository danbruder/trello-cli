package context

import (
	"strings"
	"testing"
)

// TestEstimateTokensWithEmptyString tests token estimation with empty string
func TestEstimateTokensWithEmptyString(t *testing.T) {
	optimizer := NewOptimizer(0, []string{}, false)

	result := optimizer.EstimateTokens("")
	if result != 0 {
		t.Errorf("Expected 0 tokens for empty string, got %d", result)
	}
}

// TestEstimateTokensWithUnicodeCharacters tests token estimation with unicode
func TestEstimateTokensWithUnicodeCharacters(t *testing.T) {
	optimizer := NewOptimizer(0, []string{}, false)

	// Unicode characters may take more bytes
	text := "Hello 世界 🌍"
	result := optimizer.EstimateTokens(text)

	// Should still return a reasonable estimate
	if result < 0 {
		t.Errorf("Token estimate should not be negative")
	}
}

// TestTruncateToTokenLimitWithZeroLimit tests that zero limit means no truncation
func TestTruncateToTokenLimitWithZeroLimit(t *testing.T) {
	optimizer := NewOptimizer(0, []string{}, false)

	longText := strings.Repeat("This is a very long text. ", 100)
	result := optimizer.TruncateToTokenLimit(longText)

	if result != longText {
		t.Error("Text should not be truncated when maxTokens is 0")
	}
}

// TestTruncateToTokenLimitWithNegativeLimit tests negative token limit
func TestTruncateToTokenLimitWithNegativeLimit(t *testing.T) {
	optimizer := NewOptimizer(-100, []string{}, false)

	text := "This is some text"
	result := optimizer.TruncateToTokenLimit(text)

	if result != text {
		t.Error("Text should not be truncated when maxTokens is negative")
	}
}

// TestTruncateToTokenLimitExactBoundary tests truncation at exact boundary
func TestTruncateToTokenLimitExactBoundary(t *testing.T) {
	optimizer := NewOptimizer(10, []string{}, false)

	// 10 tokens * 4 chars = 40 chars
	text := strings.Repeat("a", 40)
	result := optimizer.TruncateToTokenLimit(text)

	if result != text {
		t.Error("Text should not be truncated when exactly at limit")
	}
}

// TestShouldIncludeFieldWithEmptyFieldList tests include all when empty
func TestShouldIncludeFieldWithEmptyFieldList(t *testing.T) {
	optimizer := NewOptimizer(0, []string{}, false)

	// Should include any field when no filter is specified
	fields := []string{"id", "name", "desc", "url", "anything"}
	for _, field := range fields {
		if !optimizer.ShouldIncludeField(field) {
			t.Errorf("Should include field %s when no filter specified", field)
		}
	}
}

// TestShouldIncludeFieldCaseSensitivity tests case sensitivity
func TestShouldIncludeFieldCaseSensitivity(t *testing.T) {
	optimizer := NewOptimizer(0, []string{"id", "Name"}, false)

	// Should be case-sensitive
	if !optimizer.ShouldIncludeField("id") {
		t.Error("Should include 'id' field")
	}

	if !optimizer.ShouldIncludeField("Name") {
		t.Error("Should include 'Name' field")
	}

	// Different case should not match
	if optimizer.ShouldIncludeField("ID") {
		t.Error("Should not include 'ID' (different case)")
	}

	if optimizer.ShouldIncludeField("name") {
		t.Error("Should not include 'name' (different case)")
	}
}

// TestGetDefaultFieldsWithUnknownEntity tests default fields for unknown entity
func TestGetDefaultFieldsWithUnknownEntity(t *testing.T) {
	optimizer := NewOptimizer(0, []string{}, false)

	result := optimizer.GetDefaultFields("unknown-entity-type")

	if len(result) != 0 {
		t.Errorf("Expected empty fields for unknown entity, got %d fields", len(result))
	}
}

// TestGetDefaultFieldsVerboseVsNonVerbose tests verbose mode differences
func TestGetDefaultFieldsVerboseVsNonVerbose(t *testing.T) {
	verboseOpt := NewOptimizer(0, []string{}, true)
	nonVerboseOpt := NewOptimizer(0, []string{}, false)

	entities := []string{"board", "list", "card", "member"}

	for _, entity := range entities {
		verboseFields := verboseOpt.GetDefaultFields(entity)
		nonVerboseFields := nonVerboseOpt.GetDefaultFields(entity)

		if len(verboseFields) <= len(nonVerboseFields) {
			t.Errorf("Verbose mode should return more fields for %s (verbose: %d, non-verbose: %d)",
				entity, len(verboseFields), len(nonVerboseFields))
		}
	}
}

// TestTruncateTextWithZeroLength tests truncate with zero max length
func TestTruncateTextWithZeroLength(t *testing.T) {
	result := TruncateText("Some text", 0)

	if result != "..." {
		t.Errorf("Expected '...', got %s", result)
	}
}

// TestTruncateTextWithNegativeLength tests truncate with negative max length
func TestTruncateTextWithNegativeLength(t *testing.T) {
	result := TruncateText("Some text", -5)

	if result != "..." {
		t.Errorf("Expected '...', got %s", result)
	}
}

// TestTruncateTextPreservesLength tests that truncation doesn't exceed maxLen + 3
func TestTruncateTextPreservesLength(t *testing.T) {
	text := "This is a very long text that needs to be truncated"
	maxLen := 20

	result := TruncateText(text, maxLen)

	// Result should be maxLen + 3 (for "...")
	expectedLen := maxLen + 3
	if len(result) != expectedLen {
		t.Errorf("Expected length %d, got %d", expectedLen, len(result))
	}
}

// TestFormatSummaryWithEmptyFields tests summary with no fields
func TestFormatSummaryWithEmptyFields(t *testing.T) {
	optimizer := NewOptimizer(0, []string{}, false)

	summary := optimizer.FormatSummary("Test", 10, []string{})

	if !strings.Contains(summary, "# Test Summary") {
		t.Error("Summary should contain title")
	}

	if !strings.Contains(summary, "**Total Count:** 10") {
		t.Error("Summary should contain count")
	}

	// Should not have fields section
	if strings.Contains(summary, "**Included Fields:**") {
		t.Error("Summary should not have fields section when fields are empty")
	}
}

// TestFormatSummaryWithManyFields tests summary with multiple fields
func TestFormatSummaryWithManyFields(t *testing.T) {
	optimizer := NewOptimizer(0, []string{}, false)

	fields := []string{"id", "name", "desc", "url", "created"}
	summary := optimizer.FormatSummary("Boards", 50, fields)

	if !strings.Contains(summary, "# Boards Summary") {
		t.Error("Summary should contain title")
	}

	if !strings.Contains(summary, "**Total Count:** 50") {
		t.Error("Summary should contain count")
	}

	if !strings.Contains(summary, "**Included Fields:**") {
		t.Error("Summary should have fields section")
	}

	// Check all fields are listed
	for _, field := range fields {
		if !strings.Contains(summary, "- "+field) {
			t.Errorf("Summary should contain field %s", field)
		}
	}
}

// TestGetRelevantFieldsWithCustomFields tests relevant fields with custom fields
func TestGetRelevantFieldsWithCustomFields(t *testing.T) {
	customFields := []string{"id", "name", "custom1"}
	optimizer := NewOptimizer(0, customFields, false)

	result := optimizer.GetRelevantFields("board")

	// Should return custom fields, not default fields
	if len(result) != len(customFields) {
		t.Errorf("Expected %d fields, got %d", len(customFields), len(result))
	}

	for i, field := range customFields {
		if result[i] != field {
			t.Errorf("Expected field %s at index %d, got %s", field, i, result[i])
		}
	}
}

// TestGetRelevantFieldsWithoutCustomFields tests relevant fields without custom fields
func TestGetRelevantFieldsWithoutCustomFields(t *testing.T) {
	optimizer := NewOptimizer(0, []string{}, true)

	result := optimizer.GetRelevantFields("card")

	// Should return default fields for card in verbose mode
	expectedFields := []string{"id", "name", "desc", "url", "due", "labels", "closed", "checklists", "attachments"}

	if len(result) != len(expectedFields) {
		t.Errorf("Expected %d fields, got %d", len(expectedFields), len(result))
	}
}

// TestSummarizeCards tests the card summarization function
func TestSummarizeCards(t *testing.T) {
	optimizer := NewOptimizer(0, []string{}, false)

	// This is currently a placeholder function
	result := optimizer.SummarizeCards(nil, 10)

	if result == "" {
		t.Error("SummarizeCards should return a non-empty string")
	}
}

// TestOptimizerWithExtremeValues tests optimizer with extreme values
func TestOptimizerWithExtremeValues(t *testing.T) {
	// Very high token limit
	opt1 := NewOptimizer(1000000, []string{}, false)
	if opt1.maxTokens != 1000000 {
		t.Error("Should accept very high token limit")
	}

	// Many fields
	manyFields := make([]string, 100)
	for i := 0; i < 100; i++ {
		manyFields[i] = string(rune('a' + i%26))
	}
	opt2 := NewOptimizer(0, manyFields, false)
	if len(opt2.fields) != 100 {
		t.Error("Should accept many fields")
	}
}

// TestTruncateToTokenLimitWithVeryLongText tests truncation with extremely long text
func TestTruncateToTokenLimitWithVeryLongText(t *testing.T) {
	optimizer := NewOptimizer(100, []string{}, false)

	// Create a very long text (1MB)
	longText := strings.Repeat("a", 1024*1024)

	result := optimizer.TruncateToTokenLimit(longText)

	// Should be truncated to ~400 chars (100 tokens * 4)
	if len(result) > 500 {
		t.Errorf("Text should be truncated, got length %d", len(result))
	}

	if !strings.Contains(result, "... (truncated to fit token limit)") {
		t.Error("Truncated text should contain truncation message")
	}
}
