package context

import (
	"fmt"
	"strings"
)

// Optimizer provides LLM context optimization features
type Optimizer struct {
	maxTokens int
	fields    []string
	verbose   bool
}

// NewOptimizer creates a new context optimizer
func NewOptimizer(maxTokens int, fields []string, verbose bool) *Optimizer {
	return &Optimizer{
		maxTokens: maxTokens,
		fields:    fields,
		verbose:   verbose,
	}
}

// EstimateTokens provides a rough estimation of token count
// Uses ~4 characters per token as a rough approximation
func (o *Optimizer) EstimateTokens(text string) int {
	return len(text) / 4
}

// TruncateToTokenLimit truncates text to fit within token limit
func (o *Optimizer) TruncateToTokenLimit(text string) string {
	if o.maxTokens <= 0 {
		return text
	}

	maxChars := o.maxTokens * 4
	if len(text) <= maxChars {
		return text
	}

	return text[:maxChars] + "\n\n... (truncated to fit token limit)"
}

// ShouldIncludeField checks if a field should be included based on field filter
func (o *Optimizer) ShouldIncludeField(field string) bool {
	if len(o.fields) == 0 {
		return true
	}

	for _, f := range o.fields {
		if f == field {
			return true
		}
	}
	return false
}

// GetDefaultFields returns default field sets based on verbosity
func (o *Optimizer) GetDefaultFields(entityType string) []string {
	switch entityType {
	case "board":
		if o.verbose {
			return []string{"id", "name", "desc", "url", "closed", "dateLastActivity", "members"}
		}
		return []string{"id", "name", "desc", "closed"}

	case "list":
		if o.verbose {
			return []string{"id", "name", "closed", "pos", "cards"}
		}
		return []string{"id", "name", "closed"}

	case "card":
		if o.verbose {
			return []string{"id", "name", "desc", "url", "due", "labels", "closed", "checklists", "attachments"}
		}
		return []string{"id", "name", "desc", "due", "labels", "closed"}

	case "member":
		if o.verbose {
			return []string{"id", "username", "fullName", "url", "avatarHash"}
		}
		return []string{"id", "username", "fullName"}

	default:
		return []string{}
	}
}

// SummarizeCards provides a summary of cards for large lists
func (o *Optimizer) SummarizeCards(cards interface{}, maxItems int) string {
	// This would be implemented based on the actual card structure
	// For now, return a placeholder
	return "Card summary functionality would be implemented here"
}

// TruncateText truncates text to a maximum length
func TruncateText(text string, maxLen int) string {
	if maxLen <= 0 {
		return "..."
	}
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}

// FormatSummary creates a summary format for large datasets
func (o *Optimizer) FormatSummary(title string, count int, fields []string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s Summary\n\n", title))
	sb.WriteString(fmt.Sprintf("**Total Count:** %d\n\n", count))

	if len(fields) > 0 {
		sb.WriteString("**Included Fields:**\n")
		for _, field := range fields {
			sb.WriteString(fmt.Sprintf("- %s\n", field))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// GetRelevantFields filters fields based on context and verbosity
func (o *Optimizer) GetRelevantFields(entityType string) []string {
	if len(o.fields) > 0 {
		return o.fields
	}
	return o.GetDefaultFields(entityType)
}
