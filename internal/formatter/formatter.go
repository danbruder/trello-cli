// Package formatter provides flexible output formatting for Trello API responses.
// It supports multiple output formats (JSON, Markdown) with features optimized for
// LLM integration including field filtering, token limiting, and context optimization.
package formatter

import (
	"fmt"
)

// Formatter defines the interface for output formatters
type Formatter interface {
	FormatBoard(board interface{}) (string, error)
	FormatBoards(boards interface{}) (string, error)
	FormatList(list interface{}) (string, error)
	FormatLists(lists interface{}) (string, error)
	FormatCard(card interface{}) (string, error)
	FormatCards(cards interface{}) (string, error)
	FormatLabel(label interface{}) (string, error)
	FormatLabels(labels interface{}) (string, error)
	FormatChecklist(checklist interface{}) (string, error)
	FormatChecklists(checklists interface{}) (string, error)
	FormatMember(member interface{}) (string, error)
	FormatMembers(members interface{}) (string, error)
	FormatAttachment(attachment interface{}) (string, error)
	FormatAttachments(attachments interface{}) (string, error)
	FormatError(err error) string
	FormatSuccess(message string) string
}

// NewFormatter creates a new formatter based on the format type
func NewFormatter(format string, fields []string, maxTokens int, verbose bool) (Formatter, error) {
	switch format {
	case "json":
		return NewJSONFormatter(fields, maxTokens, verbose), nil
	case "markdown", "md":
		return NewMarkdownFormatter(fields, maxTokens, verbose), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s (supported: json, markdown)", format)
	}
}
