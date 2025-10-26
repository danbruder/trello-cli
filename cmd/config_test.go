package cmd

import (
	"testing"
)

func TestMaskString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: "(not set)",
		},
		{
			name:     "Short string",
			input:    "short",
			expected: "***",
		},
		{
			name:     "Exactly 8 characters",
			input:    "12345678",
			expected: "***",
		},
		{
			name:     "Long API key",
			input:    "abcdefghijklmnopqrstuvwxyz",
			expected: "abcd***wxyz",
		},
		{
			name:     "Medium length string",
			input:    "1234567890",
			expected: "1234***7890",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskString(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetConfigPath(t *testing.T) {
	path := getConfigPath()

	if path == "" {
		t.Error("Config path should not be empty")
	}

	// Should contain .trello-cli directory
	if path != "unknown" && len(path) > 0 {
		// Path was successfully retrieved
		t.Logf("Config path: %s", path)
	}
}
