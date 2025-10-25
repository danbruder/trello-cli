package formatter

import (
	"encoding/json"

	"github.com/adlio/trello"
)

// JSONFormatter formats output as JSON
type JSONFormatter struct {
	fields    []string
	maxTokens int
	verbose   bool
}

// NewJSONFormatter creates a new JSON formatter
func NewJSONFormatter(fields []string, maxTokens int, verbose bool) *JSONFormatter {
	return &JSONFormatter{
		fields:    fields,
		maxTokens: maxTokens,
		verbose:   verbose,
	}
}

func (f *JSONFormatter) format(data interface{}) (string, error) {
	// Apply field filtering if specified
	if len(f.fields) > 0 {
		data = extractFields(data, f.fields)
	}

	var output []byte
	var err error

	if f.verbose {
		output, err = json.MarshalIndent(data, "", "  ")
	} else {
		output, err = json.MarshalIndent(data, "", "  ")
	}

	if err != nil {
		return "", err
	}

	result := string(output)

	// Apply token limit if specified
	if f.maxTokens > 0 {
		result = truncateToTokenLimit(result, f.maxTokens)
	}

	return result, nil
}

func (f *JSONFormatter) FormatBoard(board interface{}) (string, error) {
	return f.format(board)
}

func (f *JSONFormatter) FormatBoards(boards interface{}) (string, error) {
	return f.format(boards)
}

func (f *JSONFormatter) FormatList(list interface{}) (string, error) {
	return f.format(list)
}

func (f *JSONFormatter) FormatLists(lists interface{}) (string, error) {
	return f.format(lists)
}

func (f *JSONFormatter) FormatCard(card interface{}) (string, error) {
	return f.format(card)
}

func (f *JSONFormatter) FormatCards(cards interface{}) (string, error) {
	return f.format(cards)
}

func (f *JSONFormatter) FormatLabel(label interface{}) (string, error) {
	return f.format(label)
}

func (f *JSONFormatter) FormatLabels(labels interface{}) (string, error) {
	return f.format(labels)
}

func (f *JSONFormatter) FormatChecklist(checklist interface{}) (string, error) {
	return f.format(checklist)
}

func (f *JSONFormatter) FormatChecklists(checklists interface{}) (string, error) {
	return f.format(checklists)
}

func (f *JSONFormatter) FormatMember(member interface{}) (string, error) {
	return f.format(member)
}

func (f *JSONFormatter) FormatMembers(members interface{}) (string, error) {
	return f.format(members)
}

func (f *JSONFormatter) FormatAttachment(attachment interface{}) (string, error) {
	return f.format(attachment)
}

func (f *JSONFormatter) FormatAttachments(attachments interface{}) (string, error) {
	return f.format(attachments)
}

func (f *JSONFormatter) FormatError(err error) string {
	output, _ := json.MarshalIndent(map[string]string{
		"error": err.Error(),
	}, "", "  ")
	return string(output)
}

func (f *JSONFormatter) FormatSuccess(message string) string {
	output, _ := json.MarshalIndent(map[string]string{
		"status":  "success",
		"message": message,
	}, "", "  ")
	return string(output)
}

// truncateToTokenLimit truncates the output to fit within the token limit
func truncateToTokenLimit(text string, maxTokens int) string {
	// Rough estimation: 4 characters per token
	maxChars := maxTokens * 4
	if len(text) <= maxChars {
		return text
	}
	return text[:maxChars] + "\n... (truncated)"
}

// Helper function to extract specific fields from Trello objects
func extractFields(obj interface{}, fields []string) map[string]interface{} {
	result := make(map[string]interface{})

	if len(fields) == 0 {
		// No filtering, return all
		data, _ := json.Marshal(obj)
		if err := json.Unmarshal(data, &result); err != nil {
			return nil
		}
		return result
	}

	// Convert to map for field extraction
	data, _ := json.Marshal(obj)
	var fullData map[string]interface{}
	if err := json.Unmarshal(data, &fullData); err != nil {
		return nil
	}

	for _, field := range fields {
		if val, ok := fullData[field]; ok {
			result[field] = val
		}
	}

	return result
}

// FormatBoardsWithFields formats boards with field filtering
func (f *JSONFormatter) FormatBoardsWithFields(boards []*trello.Board) (string, error) {
	if len(f.fields) == 0 {
		return f.format(boards)
	}

	filtered := make([]map[string]interface{}, len(boards))
	for i, board := range boards {
		filtered[i] = extractFields(board, f.fields)
	}

	return f.format(filtered)
}
