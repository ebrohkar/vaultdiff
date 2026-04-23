package diff

import (
	"fmt"
	"io"
	"strings"
)

// OutputFormat defines how diff results are rendered.
type OutputFormat string

const (
	FormatText OutputFormat = "text"
	FormatJSON OutputFormat = "json"
	FormatMarkdown OutputFormat = "markdown"
)

// Change represents a single key-level diff result.
type Change struct {
	Key    string
	Type   ChangeType
	OldVal string
	NewVal string
}

// ChangeType categorizes the nature of a diff change.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
)

// Render writes formatted diff output to w using the specified format.
func Render(w io.Writer, changes []Change, format OutputFormat) error {
	switch format {
	case FormatText:
		return renderText(w, changes)
	case FormatJSON:
		return renderJSON(w, changes)
	case FormatMarkdown:
		return renderMarkdown(w, changes)
	default:
		return fmt.Errorf("unsupported output format: %q", format)
	}
}

func renderText(w io.Writer, changes []Change) error {
	if len(changes) == 0 {
		_, err := fmt.Fprintln(w, "No changes detected.")
		return err
	}
	for _, c := range changes {
		switch c.Type {
		case Added:
			fmt.Fprintf(w, "+ [%s] %s\n", c.Type, c.Key)
		case Removed:
			fmt.Fprintf(w, "- [%s] %s\n", c.Type, c.Key)
		case Modified:
			fmt.Fprintf(w, "~ [%s] %s\n", c.Type, c.Key)
		}
	}
	return nil
}

func renderJSON(w io.Writer, changes []Change) error {
	var sb strings.Builder
	sb.WriteString("[\n")
	for i, c := range changes {
		sb.WriteString(fmt.Sprintf(
			"  {\"key\": %q, \"type\": %q, \"old\": %q, \"new\": %q}",
			c.Key, c.Type, c.OldVal, c.NewVal,
		))
		if i < len(changes)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("]\n")
	_, err := fmt.Fprint(w, sb.String())
	return err
}

func renderMarkdown(w io.Writer, changes []Change) error {
	if len(changes) == 0 {
		_, err := fmt.Fprintln(w, "_No changes detected._")
		return err
	}
	fmt.Fprintln(w, "| Key | Change | Old | New |")
	fmt.Fprintln(w, "|-----|--------|-----|-----|")
	for _, c := range changes {
		fmt.Fprintf(w, "| %s | %s | %s | %s |\n", c.Key, c.Type, c.OldVal, c.NewVal)
	}
	return nil
}
