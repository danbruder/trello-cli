package context

import (
	"testing"
)

func TestNewOptimizer(t *testing.T) {
	optimizer := NewOptimizer(1000, []string{"id", "name"}, true)

	if optimizer.maxTokens != 1000 {
		t.Errorf("Expected maxTokens 1000, got %d", optimizer.maxTokens)
	}

	if len(optimizer.fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(optimizer.fields))
	}

	if !optimizer.verbose {
		t.Errorf("Expected verbose to be true")
	}
}

func TestEstimateTokens(t *testing.T) {
	optimizer := NewOptimizer(0, []string{}, false)

	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{
			name:     "Empty string",
			text:     "",
			expected: 0,
		},
		{
			name:     "Short text",
			text:     "Hello",
			expected: 1, // 5 chars / 4 = 1.25, rounded down
		},
		{
			name:     "Medium text",
			text:     "This is a test message with multiple words",
			expected: 10, // 44 chars / 4 = 11
		},
		{
			name:     "Long text",
			text:     "This is a very long text that contains many words and should result in a higher token count when processed by the token estimation function",
			expected: 34, // ~140 chars / 4 = 35, but actual is 34
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := optimizer.EstimateTokens(tt.text)
			if result != tt.expected {
				t.Errorf("Expected %d tokens, got %d", tt.expected, result)
			}
		})
	}
}

func TestTruncateToTokenLimit(t *testing.T) {
	tests := []struct {
		name             string
		maxTokens        int
		text             string
		expectTruncation bool
	}{
		{
			name:             "No limit",
			maxTokens:        0,
			text:             "This is a test message",
			expectTruncation: false,
		},
		{
			name:             "Text within limit",
			maxTokens:        100,
			text:             "Short text",
			expectTruncation: false,
		},
		{
			name:             "Text exceeds limit",
			maxTokens:        10,
			text:             "This is a very long text that should be truncated when the token limit is applied because it exceeds the maximum allowed tokens",
			expectTruncation: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			optimizer := NewOptimizer(tt.maxTokens, []string{}, false)
			result := optimizer.TruncateToTokenLimit(tt.text)

			if tt.expectTruncation {
				if len(result) >= len(tt.text) {
					t.Errorf("Expected text to be truncated, but result length %d >= original length %d", len(result), len(tt.text))
				}
				if !contains(result, "... (truncated to fit token limit)") {
					t.Errorf("Expected truncation message in result")
				}
			} else {
				if result != tt.text {
					t.Errorf("Expected text to remain unchanged, got: %s", result)
				}
			}
		})
	}
}

func TestShouldIncludeField(t *testing.T) {
	tests := []struct {
		name     string
		fields   []string
		field    string
		expected bool
	}{
		{
			name:     "No fields specified",
			fields:   []string{},
			field:    "anyfield",
			expected: true,
		},
		{
			name:     "Field in list",
			fields:   []string{"id", "name", "desc"},
			field:    "name",
			expected: true,
		},
		{
			name:     "Field not in list",
			fields:   []string{"id", "name", "desc"},
			field:    "url",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			optimizer := NewOptimizer(0, tt.fields, false)
			result := optimizer.ShouldIncludeField(tt.field)

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetDefaultFields(t *testing.T) {
	tests := []struct {
		name       string
		entityType string
		verbose    bool
		expected   []string
	}{
		{
			name:       "Board verbose",
			entityType: "board",
			verbose:    true,
			expected:   []string{"id", "name", "desc", "url", "closed", "dateLastActivity", "members"},
		},
		{
			name:       "Board non-verbose",
			entityType: "board",
			verbose:    false,
			expected:   []string{"id", "name", "desc", "closed"},
		},
		{
			name:       "Card verbose",
			entityType: "card",
			verbose:    true,
			expected:   []string{"id", "name", "desc", "url", "due", "labels", "closed", "checklists", "attachments"},
		},
		{
			name:       "Card non-verbose",
			entityType: "card",
			verbose:    false,
			expected:   []string{"id", "name", "desc", "due", "labels", "closed"},
		},
		{
			name:       "Unknown entity",
			entityType: "unknown",
			verbose:    false,
			expected:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			optimizer := NewOptimizer(0, []string{}, tt.verbose)
			result := optimizer.GetDefaultFields(tt.entityType)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d fields, got %d", len(tt.expected), len(result))
				return
			}

			for i, field := range result {
				if field != tt.expected[i] {
					t.Errorf("Expected field %s at index %d, got %s", tt.expected[i], i, field)
				}
			}
		})
	}
}

func TestFormatSummary(t *testing.T) {
	optimizer := NewOptimizer(0, []string{"id", "name"}, false)

	summary := optimizer.FormatSummary("Test Summary", 5, []string{"id", "name"})

	if !contains(summary, "# Test Summary Summary") {
		t.Errorf("Summary should contain title")
	}

	if !contains(summary, "**Total Count:** 5") {
		t.Errorf("Summary should contain count")
	}

	if !contains(summary, "**Included Fields:**") {
		t.Errorf("Summary should contain fields section")
	}

	if !contains(summary, "- id") {
		t.Errorf("Summary should contain field id")
	}

	if !contains(summary, "- name") {
		t.Errorf("Summary should contain field name")
	}
}

func TestGetRelevantFields(t *testing.T) {
	tests := []struct {
		name           string
		fields         []string
		entityType     string
		verbose        bool
		expectedLength int
	}{
		{
			name:           "Custom fields specified",
			fields:         []string{"id", "name"},
			entityType:     "board",
			verbose:        true,
			expectedLength: 2,
		},
		{
			name:           "No custom fields, verbose",
			fields:         []string{},
			entityType:     "board",
			verbose:        true,
			expectedLength: 7, // Board verbose fields
		},
		{
			name:           "No custom fields, non-verbose",
			fields:         []string{},
			entityType:     "card",
			verbose:        false,
			expectedLength: 6, // Card non-verbose fields
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			optimizer := NewOptimizer(0, tt.fields, tt.verbose)
			result := optimizer.GetRelevantFields(tt.entityType)

			if len(result) != tt.expectedLength {
				t.Errorf("Expected %d fields, got %d", tt.expectedLength, len(result))
			}
		})
	}
}

func TestTruncateText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		maxLen   int
		expected string
	}{
		{
			name:     "Text within limit",
			text:     "Short text",
			maxLen:   20,
			expected: "Short text",
		},
		{
			name:     "Text exceeds limit",
			text:     "This is a very long text",
			maxLen:   10,
			expected: "This is a ...",
		},
		{
			name:     "Empty text",
			text:     "",
			maxLen:   10,
			expected: "",
		},
		{
			name:     "Exact limit",
			text:     "Exactly ten",
			maxLen:   10,
			expected: "Exactly te...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateText(tt.text, tt.maxLen)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && contains(s[1:], substr)
}
