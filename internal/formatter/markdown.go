package formatter

import (
	"fmt"
	"strings"
	"time"

	"github.com/adlio/trello"
)

// MarkdownFormatter formats output as Markdown
type MarkdownFormatter struct {
	fields    []string
	maxTokens int
	verbose   bool
}

// NewMarkdownFormatter creates a new Markdown formatter
func NewMarkdownFormatter(fields []string, maxTokens int, verbose bool) *MarkdownFormatter {
	return &MarkdownFormatter{
		fields:    fields,
		maxTokens: maxTokens,
		verbose:   verbose,
	}
}

func (f *MarkdownFormatter) FormatBoard(board interface{}) (string, error) {
	b, ok := board.(*trello.Board)
	if !ok {
		return "", fmt.Errorf("invalid board type")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Board: %s\n\n", b.Name))
	sb.WriteString(fmt.Sprintf("**ID:** `%s`\n\n", b.ID))

	if f.verbose || (len(f.fields) > 0 && f.shouldIncludeField("desc")) {
		if b.Desc != "" {
			sb.WriteString(fmt.Sprintf("**Description:** %s\n\n", b.Desc))
		}
	}

	if f.verbose || (len(f.fields) > 0 && f.shouldIncludeField("url")) {
		sb.WriteString(fmt.Sprintf("**URL:** %s\n\n", b.URL))
	}

	if f.verbose || (len(f.fields) > 0 && f.shouldIncludeField("closed")) {
		sb.WriteString(fmt.Sprintf("**Closed:** %t\n\n", b.Closed))
	}

	return f.applyTokenLimit(sb.String()), nil
}

func (f *MarkdownFormatter) FormatBoards(boards interface{}) (string, error) {
	boardList, ok := boards.([]*trello.Board)
	if !ok {
		return "", fmt.Errorf("invalid boards type")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Boards (%d)\n\n", len(boardList)))

	for _, board := range boardList {
		sb.WriteString(fmt.Sprintf("## %s\n", board.Name))
		sb.WriteString(fmt.Sprintf("- **ID:** `%s`\n", board.ID))

		if f.verbose || f.shouldIncludeField("desc") {
			if board.Desc != "" {
				sb.WriteString(fmt.Sprintf("- **Description:** %s\n", truncateText(board.Desc, 100)))
			}
		}

		if f.verbose || f.shouldIncludeField("url") {
			sb.WriteString(fmt.Sprintf("- **URL:** %s\n", board.URL))
		}

		sb.WriteString("\n")
	}

	return f.applyTokenLimit(sb.String()), nil
}

func (f *MarkdownFormatter) FormatList(list interface{}) (string, error) {
	l, ok := list.(*trello.List)
	if !ok {
		return "", fmt.Errorf("invalid list type")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# List: %s\n\n", l.Name))
	sb.WriteString(fmt.Sprintf("**ID:** `%s`\n\n", l.ID))

	if f.verbose || f.shouldIncludeField("closed") {
		sb.WriteString(fmt.Sprintf("**Closed:** %t\n\n", l.Closed))
	}

	if f.verbose || f.shouldIncludeField("pos") {
		sb.WriteString(fmt.Sprintf("**Position:** %.2f\n\n", l.Pos))
	}

	return f.applyTokenLimit(sb.String()), nil
}

func (f *MarkdownFormatter) FormatLists(lists interface{}) (string, error) {
	listList, ok := lists.([]*trello.List)
	if !ok {
		return "", fmt.Errorf("invalid lists type")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Lists (%d)\n\n", len(listList)))

	for _, list := range listList {
		status := "Active"
		if list.Closed {
			status = "Archived"
		}
		sb.WriteString(fmt.Sprintf("## %s [%s]\n", list.Name, status))
		sb.WriteString(fmt.Sprintf("- **ID:** `%s`\n", list.ID))

		if f.verbose || f.shouldIncludeField("pos") {
			sb.WriteString(fmt.Sprintf("- **Position:** %.2f\n", list.Pos))
		}

		sb.WriteString("\n")
	}

	return f.applyTokenLimit(sb.String()), nil
}

func (f *MarkdownFormatter) FormatCard(card interface{}) (string, error) {
	c, ok := card.(*trello.Card)
	if !ok {
		return "", fmt.Errorf("invalid card type")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Card: %s\n\n", c.Name))
	sb.WriteString(fmt.Sprintf("**ID:** `%s`\n\n", c.ID))

	if f.verbose || f.shouldIncludeField("desc") {
		if c.Desc != "" {
			sb.WriteString(fmt.Sprintf("**Description:**\n\n%s\n\n", c.Desc))
		}
	}

	if f.verbose || f.shouldIncludeField("url") {
		sb.WriteString(fmt.Sprintf("**URL:** %s\n\n", c.URL))
	}

	if f.verbose || f.shouldIncludeField("due") {
		if c.Due != nil {
			sb.WriteString(fmt.Sprintf("**Due:** %s\n\n", c.Due.Format(time.RFC3339)))
		}
	}

	if f.verbose || f.shouldIncludeField("labels") {
		if len(c.Labels) > 0 {
			sb.WriteString("**Labels:**\n")
			for _, label := range c.Labels {
				color := label.Color
				if color == "" {
					color = "no color"
				}
				sb.WriteString(fmt.Sprintf("- %s (%s)\n", label.Name, color))
			}
			sb.WriteString("\n")
		}
	}

	if f.verbose || f.shouldIncludeField("closed") {
		sb.WriteString(fmt.Sprintf("**Closed:** %t\n\n", c.Closed))
	}

	return f.applyTokenLimit(sb.String()), nil
}

func (f *MarkdownFormatter) FormatCards(cards interface{}) (string, error) {
	cardList, ok := cards.([]*trello.Card)
	if !ok {
		return "", fmt.Errorf("invalid cards type")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Cards (%d)\n\n", len(cardList)))

	for _, card := range cardList {
		status := "Active"
		if card.Closed {
			status = "Archived"
		}
		sb.WriteString(fmt.Sprintf("## %s [%s]\n", card.Name, status))
		sb.WriteString(fmt.Sprintf("- **ID:** `%s`\n", card.ID))

		if f.verbose || f.shouldIncludeField("desc") {
			if card.Desc != "" {
				sb.WriteString(fmt.Sprintf("- **Description:** %s\n", truncateText(card.Desc, 150)))
			}
		}

		if f.verbose || f.shouldIncludeField("labels") {
			if len(card.Labels) > 0 {
				labelNames := make([]string, len(card.Labels))
				for i, label := range card.Labels {
					if label.Name != "" {
						labelNames[i] = label.Name
					} else {
						labelNames[i] = label.Color
					}
				}
				sb.WriteString(fmt.Sprintf("- **Labels:** %s\n", strings.Join(labelNames, ", ")))
			}
		}

		if f.verbose || f.shouldIncludeField("due") {
			if card.Due != nil {
				sb.WriteString(fmt.Sprintf("- **Due:** %s\n", card.Due.Format("2006-01-02")))
			}
		}

		sb.WriteString("\n")
	}

	return f.applyTokenLimit(sb.String()), nil
}

func (f *MarkdownFormatter) FormatLabel(label interface{}) (string, error) {
	l, ok := label.(*trello.Label)
	if !ok {
		return "", fmt.Errorf("invalid label type")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Label: %s\n\n", l.Name))
	sb.WriteString(fmt.Sprintf("**ID:** `%s`\n\n", l.ID))
	sb.WriteString(fmt.Sprintf("**Color:** %s\n\n", l.Color))

	return sb.String(), nil
}

func (f *MarkdownFormatter) FormatLabels(labels interface{}) (string, error) {
	labelList, ok := labels.([]*trello.Label)
	if !ok {
		return "", fmt.Errorf("invalid labels type")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Labels (%d)\n\n", len(labelList)))

	for _, label := range labelList {
		name := label.Name
		if name == "" {
			name = "(unnamed)"
		}
		sb.WriteString(fmt.Sprintf("- **%s** - `%s` (ID: `%s`)\n", name, label.Color, label.ID))
	}

	return sb.String(), nil
}

func (f *MarkdownFormatter) FormatChecklist(checklist interface{}) (string, error) {
	cl, ok := checklist.(*trello.Checklist)
	if !ok {
		return "", fmt.Errorf("invalid checklist type")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Checklist: %s\n\n", cl.Name))
	sb.WriteString(fmt.Sprintf("**ID:** `%s`\n\n", cl.ID))

	if len(cl.CheckItems) > 0 {
		sb.WriteString("## Items:\n\n")
		for _, item := range cl.CheckItems {
			checkbox := "[ ]"
			if item.State == "complete" {
				checkbox = "[x]"
			}
			sb.WriteString(fmt.Sprintf("- %s %s\n", checkbox, item.Name))
		}
	}

	return sb.String(), nil
}

func (f *MarkdownFormatter) FormatChecklists(checklists interface{}) (string, error) {
	checklistList, ok := checklists.([]*trello.Checklist)
	if !ok {
		return "", fmt.Errorf("invalid checklists type")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Checklists (%d)\n\n", len(checklistList)))

	for _, checklist := range checklistList {
		completed := 0
		for _, item := range checklist.CheckItems {
			if item.State == "complete" {
				completed++
			}
		}
		sb.WriteString(fmt.Sprintf("## %s\n", checklist.Name))
		sb.WriteString(fmt.Sprintf("- **ID:** `%s`\n", checklist.ID))
		sb.WriteString(fmt.Sprintf("- **Progress:** %d/%d items complete\n\n", completed, len(checklist.CheckItems)))
	}

	return sb.String(), nil
}

func (f *MarkdownFormatter) FormatMember(member interface{}) (string, error) {
	m, ok := member.(*trello.Member)
	if !ok {
		return "", fmt.Errorf("invalid member type")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Member: %s\n\n", m.FullName))
	sb.WriteString(fmt.Sprintf("**ID:** `%s`\n\n", m.ID))
	sb.WriteString(fmt.Sprintf("**Username:** %s\n\n", m.Username))

	if f.verbose || f.shouldIncludeField("url") {
		sb.WriteString(fmt.Sprintf("**URL:** %s\n\n", m.Username))
	}

	return sb.String(), nil
}

func (f *MarkdownFormatter) FormatMembers(members interface{}) (string, error) {
	memberList, ok := members.([]*trello.Member)
	if !ok {
		return "", fmt.Errorf("invalid members type")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Members (%d)\n\n", len(memberList)))

	for _, member := range memberList {
		sb.WriteString(fmt.Sprintf("- **%s** (@%s) - ID: `%s`\n", member.FullName, member.Username, member.ID))
	}

	return sb.String(), nil
}

func (f *MarkdownFormatter) FormatAttachment(attachment interface{}) (string, error) {
	a, ok := attachment.(*trello.Attachment)
	if !ok {
		return "", fmt.Errorf("invalid attachment type")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Attachment: %s\n\n", a.Name))
	sb.WriteString(fmt.Sprintf("**ID:** `%s`\n\n", a.ID))
	sb.WriteString(fmt.Sprintf("**URL:** %s\n\n", a.URL))

	return sb.String(), nil
}

func (f *MarkdownFormatter) FormatAttachments(attachments interface{}) (string, error) {
	attachmentList, ok := attachments.([]*trello.Attachment)
	if !ok {
		return "", fmt.Errorf("invalid attachments type")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Attachments (%d)\n\n", len(attachmentList)))

	for _, attachment := range attachmentList {
		sb.WriteString(fmt.Sprintf("- **%s** - [Link](%s) (ID: `%s`)\n", attachment.Name, attachment.URL, attachment.ID))
	}

	return sb.String(), nil
}

func (f *MarkdownFormatter) FormatError(err error) string {
	return fmt.Sprintf("❌ **Error:** %s\n", err.Error())
}

func (f *MarkdownFormatter) FormatSuccess(message string) string {
	return fmt.Sprintf("✅ **Success:** %s\n", message)
}

// Helper functions

func (f *MarkdownFormatter) shouldIncludeField(field string) bool {
	if len(f.fields) == 0 {
		return true
	}
	for _, f := range f.fields {
		if f == field {
			return true
		}
	}
	return false
}

func (f *MarkdownFormatter) applyTokenLimit(text string) string {
	if f.maxTokens <= 0 {
		return text
	}

	// Rough estimation: 4 characters per token
	maxChars := f.maxTokens * 4
	if len(text) <= maxChars {
		return text
	}

	return text[:maxChars] + "\n\n... (output truncated to fit token limit)"
}

func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}
