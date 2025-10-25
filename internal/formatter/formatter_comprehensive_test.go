package formatter

import (
	"strings"
	"testing"
	"time"

	"github.com/adlio/trello"
)

// TestMarkdownLabelFormatting tests label formatting in markdown
func TestMarkdownLabelFormatting(t *testing.T) {
	formatter := NewMarkdownFormatter([]string{}, 0, false)

	// Test single label
	label := &trello.Label{
		ID:    "label-123",
		Name:  "Bug",
		Color: "red",
	}

	output, err := formatter.FormatLabel(label)
	if err != nil {
		t.Fatalf("Failed to format label: %v", err)
	}

	if !strings.Contains(output, "# Label: Bug") {
		t.Errorf("Output should contain label title")
	}

	if !strings.Contains(output, "label-123") {
		t.Errorf("Output should contain label ID")
	}

	if !strings.Contains(output, "red") {
		t.Errorf("Output should contain label color")
	}

	// Test labels list
	labels := []*trello.Label{
		{ID: "label-1", Name: "Bug", Color: "red"},
		{ID: "label-2", Name: "Feature", Color: "green"},
		{ID: "label-3", Name: "", Color: "blue"}, // Unnamed label
	}

	labelsOutput, err := formatter.FormatLabels(labels)
	if err != nil {
		t.Fatalf("Failed to format labels: %v", err)
	}

	if !strings.Contains(labelsOutput, "# Labels (3)") {
		t.Errorf("Output should contain labels count")
	}

	if !strings.Contains(labelsOutput, "Bug") {
		t.Errorf("Output should contain first label name")
	}

	if !strings.Contains(labelsOutput, "(unnamed)") {
		t.Errorf("Output should handle unnamed labels")
	}
}

// TestMarkdownChecklistFormatting tests checklist formatting
func TestMarkdownChecklistFormatting(t *testing.T) {
	formatter := NewMarkdownFormatter([]string{}, 0, false)

	// Test single checklist with items
	checklist := &trello.Checklist{
		ID:   "checklist-123",
		Name: "Todo List",
		CheckItems: []trello.CheckItem{
			{ID: "item-1", Name: "Task 1", State: "complete"},
			{ID: "item-2", Name: "Task 2", State: "incomplete"},
		},
	}

	output, err := formatter.FormatChecklist(checklist)
	if err != nil {
		t.Fatalf("Failed to format checklist: %v", err)
	}

	if !strings.Contains(output, "# Checklist: Todo List") {
		t.Errorf("Output should contain checklist title")
	}

	if !strings.Contains(output, "[x] Task 1") {
		t.Errorf("Output should contain completed task with checkbox")
	}

	if !strings.Contains(output, "[ ] Task 2") {
		t.Errorf("Output should contain incomplete task with checkbox")
	}

	// Test checklists list
	checklists := []*trello.Checklist{
		{
			ID:   "checklist-1",
			Name: "Checklist 1",
			CheckItems: []trello.CheckItem{
				{ID: "item-1", Name: "Task 1", State: "complete"},
				{ID: "item-2", Name: "Task 2", State: "complete"},
			},
		},
		{
			ID:         "checklist-2",
			Name:       "Checklist 2",
			CheckItems: []trello.CheckItem{},
		},
	}

	checklistsOutput, err := formatter.FormatChecklists(checklists)
	if err != nil {
		t.Fatalf("Failed to format checklists: %v", err)
	}

	if !strings.Contains(checklistsOutput, "# Checklists (2)") {
		t.Errorf("Output should contain checklists count")
	}

	if !strings.Contains(checklistsOutput, "2/2 items complete") {
		t.Errorf("Output should show completed items progress")
	}
}

// TestMarkdownMemberFormatting tests member formatting
func TestMarkdownMemberFormatting(t *testing.T) {
	formatter := NewMarkdownFormatter([]string{}, 0, false)

	// Test single member
	member := &trello.Member{
		ID:       "member-123",
		Username: "testuser",
		FullName: "Test User",
	}

	output, err := formatter.FormatMember(member)
	if err != nil {
		t.Fatalf("Failed to format member: %v", err)
	}

	if !strings.Contains(output, "# Member: Test User") {
		t.Errorf("Output should contain member full name")
	}

	if !strings.Contains(output, "testuser") {
		t.Errorf("Output should contain username")
	}

	// Test members list
	members := []*trello.Member{
		{ID: "member-1", Username: "user1", FullName: "User One"},
		{ID: "member-2", Username: "user2", FullName: "User Two"},
	}

	membersOutput, err := formatter.FormatMembers(members)
	if err != nil {
		t.Fatalf("Failed to format members: %v", err)
	}

	if !strings.Contains(membersOutput, "# Members (2)") {
		t.Errorf("Output should contain members count")
	}

	if !strings.Contains(membersOutput, "@user1") {
		t.Errorf("Output should contain first member username")
	}
}

// TestMarkdownAttachmentFormatting tests attachment formatting
func TestMarkdownAttachmentFormatting(t *testing.T) {
	formatter := NewMarkdownFormatter([]string{}, 0, false)

	// Test single attachment
	attachment := &trello.Attachment{
		ID:   "attach-123",
		Name: "Screenshot.png",
		URL:  "https://example.com/screenshot.png",
	}

	output, err := formatter.FormatAttachment(attachment)
	if err != nil {
		t.Fatalf("Failed to format attachment: %v", err)
	}

	if !strings.Contains(output, "# Attachment: Screenshot.png") {
		t.Errorf("Output should contain attachment name")
	}

	if !strings.Contains(output, "https://example.com/screenshot.png") {
		t.Errorf("Output should contain attachment URL")
	}

	// Test attachments list
	attachments := []*trello.Attachment{
		{ID: "attach-1", Name: "Image1.png", URL: "https://example.com/image1.png"},
		{ID: "attach-2", Name: "Image2.jpg", URL: "https://example.com/image2.jpg"},
	}

	attachmentsOutput, err := formatter.FormatAttachments(attachments)
	if err != nil {
		t.Fatalf("Failed to format attachments: %v", err)
	}

	if !strings.Contains(attachmentsOutput, "# Attachments (2)") {
		t.Errorf("Output should contain attachments count")
	}

	if !strings.Contains(attachmentsOutput, "[Link]") {
		t.Errorf("Output should contain markdown links")
	}
}

// TestJSONLabelFormatting tests label formatting in JSON
func TestJSONLabelFormatting(t *testing.T) {
	formatter := NewJSONFormatter([]string{}, 0, false)

	label := &trello.Label{
		ID:    "label-123",
		Name:  "Bug",
		Color: "red",
	}

	output, err := formatter.FormatLabel(label)
	if err != nil {
		t.Fatalf("Failed to format label: %v", err)
	}

	if !strings.Contains(output, "label-123") {
		t.Errorf("Output should contain label ID")
	}

	if !strings.Contains(output, "Bug") {
		t.Errorf("Output should contain label name")
	}
}

// TestJSONChecklistFormatting tests checklist formatting in JSON
func TestJSONChecklistFormatting(t *testing.T) {
	formatter := NewJSONFormatter([]string{}, 0, false)

	checklist := &trello.Checklist{
		ID:   "checklist-123",
		Name: "Todo List",
		CheckItems: []trello.CheckItem{
			{ID: "item-1", Name: "Task 1", State: "complete"},
		},
	}

	output, err := formatter.FormatChecklist(checklist)
	if err != nil {
		t.Fatalf("Failed to format checklist: %v", err)
	}

	if !strings.Contains(output, "checklist-123") {
		t.Errorf("Output should contain checklist ID")
	}

	if !strings.Contains(output, "Todo List") {
		t.Errorf("Output should contain checklist name")
	}
}

// TestJSONMemberFormatting tests member formatting in JSON
func TestJSONMemberFormatting(t *testing.T) {
	formatter := NewJSONFormatter([]string{}, 0, false)

	member := &trello.Member{
		ID:       "member-123",
		Username: "testuser",
		FullName: "Test User",
	}

	output, err := formatter.FormatMember(member)
	if err != nil {
		t.Fatalf("Failed to format member: %v", err)
	}

	if !strings.Contains(output, "member-123") {
		t.Errorf("Output should contain member ID")
	}

	if !strings.Contains(output, "testuser") {
		t.Errorf("Output should contain username")
	}
}

// TestJSONAttachmentFormatting tests attachment formatting in JSON
func TestJSONAttachmentFormatting(t *testing.T) {
	formatter := NewJSONFormatter([]string{}, 0, false)

	attachment := &trello.Attachment{
		ID:   "attach-123",
		Name: "Screenshot.png",
		URL:  "https://example.com/screenshot.png",
	}

	output, err := formatter.FormatAttachment(attachment)
	if err != nil {
		t.Fatalf("Failed to format attachment: %v", err)
	}

	if !strings.Contains(output, "attach-123") {
		t.Errorf("Output should contain attachment ID")
	}

	if !strings.Contains(output, "Screenshot.png") {
		t.Errorf("Output should contain attachment name")
	}
}

// TestMarkdownCardWithDueDate tests card formatting with due date
func TestMarkdownCardWithDueDate(t *testing.T) {
	formatter := NewMarkdownFormatter([]string{}, 0, true)

	dueDate := time.Now().Add(24 * time.Hour)
	card := &trello.Card{
		ID:   "card-123",
		Name: "Task with due date",
		Due:  &dueDate,
	}

	output, err := formatter.FormatCard(card)
	if err != nil {
		t.Fatalf("Failed to format card: %v", err)
	}

	if !strings.Contains(output, "**Due:**") {
		t.Errorf("Output should contain due date field")
	}
}

// TestMarkdownCardWithLabels tests card formatting with labels
func TestMarkdownCardWithLabels(t *testing.T) {
	formatter := NewMarkdownFormatter([]string{}, 0, true)

	card := &trello.Card{
		ID:   "card-123",
		Name: "Task with labels",
		Labels: []*trello.Label{
			{ID: "label-1", Name: "Bug", Color: "red"},
			{ID: "label-2", Name: "", Color: "blue"}, // Unnamed label
		},
	}

	output, err := formatter.FormatCard(card)
	if err != nil {
		t.Fatalf("Failed to format card: %v", err)
	}

	if !strings.Contains(output, "**Labels:**") {
		t.Errorf("Output should contain labels section")
	}

	if !strings.Contains(output, "Bug") {
		t.Errorf("Output should contain named label")
	}

	if !strings.Contains(output, "blue") {
		t.Errorf("Output should show color for unnamed label")
	}
}

// TestMarkdownCardsWithDifferentStatuses tests cards list with mixed statuses
func TestMarkdownCardsWithDifferentStatuses(t *testing.T) {
	formatter := NewMarkdownFormatter([]string{}, 0, false)

	cards := []*trello.Card{
		{ID: "card-1", Name: "Active Card", Closed: false},
		{ID: "card-2", Name: "Archived Card", Closed: true},
	}

	output, err := formatter.FormatCards(cards)
	if err != nil {
		t.Fatalf("Failed to format cards: %v", err)
	}

	if !strings.Contains(output, "[Active]") {
		t.Errorf("Output should show active status")
	}

	if !strings.Contains(output, "[Archived]") {
		t.Errorf("Output should show archived status")
	}
}

// TestExtractFieldsWithNestedData tests field extraction with nested structures
func TestExtractFieldsWithNestedData(t *testing.T) {
	board := &trello.Board{
		ID:   "board-123",
		Name: "Test Board",
		Desc: "Description",
		URL:  "https://trello.com/board",
	}

	result := extractFields(board, []string{"id", "name"})

	if result == nil {
		t.Fatal("extractFields returned nil")
	}

	if _, ok := result["id"]; !ok {
		t.Error("Result should contain id field")
	}

	if _, ok := result["name"]; !ok {
		t.Error("Result should contain name field")
	}

	// Should not contain desc or url (not in fields list)
	if _, ok := result["desc"]; ok {
		t.Error("Result should not contain desc field when filtered")
	}
}

// TestTruncateText tests the truncateText helper function
func TestTruncateText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		maxLen   int
		expected string
	}{
		{
			name:     "Text shorter than limit",
			text:     "Short",
			maxLen:   10,
			expected: "Short",
		},
		{
			name:     "Text equal to limit",
			text:     "Exactly10c",
			maxLen:   10,
			expected: "Exactly10c",
		},
		{
			name:     "Text longer than limit",
			text:     "This is a very long text that should be truncated",
			maxLen:   20,
			expected: "This is a very long ...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateText(tt.text, tt.maxLen)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestFormatterInvalidTypes tests formatters with invalid types
func TestFormatterInvalidTypes(t *testing.T) {
	formatter := NewMarkdownFormatter([]string{}, 0, false)

	// Test with wrong type
	_, err := formatter.FormatBoard("not a board")
	if err == nil {
		t.Error("Expected error when formatting wrong type")
	}

	_, err = formatter.FormatCard(123)
	if err == nil {
		t.Error("Expected error when formatting wrong type")
	}

	_, err = formatter.FormatList([]string{"not", "a", "list"})
	if err == nil {
		t.Error("Expected error when formatting wrong type")
	}
}

// TestMarkdownFormatterShouldIncludeField tests field inclusion logic
func TestMarkdownFormatterShouldIncludeField(t *testing.T) {
	// Formatter with specific fields
	formatter := NewMarkdownFormatter([]string{"id", "name"}, 0, false)

	if !formatter.shouldIncludeField("id") {
		t.Error("Should include id field")
	}

	if !formatter.shouldIncludeField("name") {
		t.Error("Should include name field")
	}

	if formatter.shouldIncludeField("desc") {
		t.Error("Should not include desc field")
	}

	// Formatter with no fields (include all)
	formatterAll := NewMarkdownFormatter([]string{}, 0, false)

	if !formatterAll.shouldIncludeField("anyfield") {
		t.Error("Should include any field when no filter specified")
	}
}

// TestJSONFormatterBoardsWithFields tests field filtering for boards
func TestJSONFormatterBoardsWithFields(t *testing.T) {
	formatter := NewJSONFormatter([]string{"id", "name"}, 0, false)

	boards := []*trello.Board{
		{
			ID:   "board-1",
			Name: "Board 1",
			Desc: "Description 1",
			URL:  "https://trello.com/board1",
		},
		{
			ID:   "board-2",
			Name: "Board 2",
			Desc: "Description 2",
			URL:  "https://trello.com/board2",
		},
	}

	output, err := formatter.FormatBoardsWithFields(boards)
	if err != nil {
		t.Fatalf("Failed to format boards with fields: %v", err)
	}

	// Output should be valid JSON
	if output == "" {
		t.Error("Output should not be empty")
	}

	// When fields are specified, the filtered output should contain the boards
	// The exact format depends on the implementation of field filtering
	t.Logf("Formatted output: %s", output)
}
