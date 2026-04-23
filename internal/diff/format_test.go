package diff

import (
	"strings"
	"testing"
)

var sampleChanges = []Change{
	{Key: "DB_HOST", Type: Added, OldVal: "", NewVal: "localhost"},
	{Key: "API_KEY", Type: Removed, OldVal: "secret123", NewVal: ""},
	{Key: "TIMEOUT", Type: Modified, OldVal: "30", NewVal: "60"},
}

func TestRender_TextNoChanges(t *testing.T) {
	var buf strings.Builder
	if err := Render(&buf, []Change{}, FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No changes detected.") {
		t.Errorf("expected no-changes message, got: %q", buf.String())
	}
}

func TestRender_TextWithChanges(t *testing.T) {
	var buf strings.Builder
	if err := Render(&buf, sampleChanges, FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "+ [added] DB_HOST") {
		t.Errorf("missing added line, got: %q", out)
	}
	if !strings.Contains(out, "- [removed] API_KEY") {
		t.Errorf("missing removed line, got: %q", out)
	}
	if !strings.Contains(out, "~ [modified] TIMEOUT") {
		t.Errorf("missing modified line, got: %q", out)
	}
}

func TestRender_JSON(t *testing.T) {
	var buf strings.Builder
	if err := Render(&buf, sampleChanges, FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.HasPrefix(out, "[") || !strings.Contains(out, "]") {
		t.Errorf("expected JSON array, got: %q", out)
	}
	if !strings.Contains(out, `"key": "DB_HOST"`) {
		t.Errorf("expected DB_HOST in JSON output, got: %q", out)
	}
}

func TestRender_Markdown(t *testing.T) {
	var buf strings.Builder
	if err := Render(&buf, sampleChanges, FormatMarkdown); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "| Key | Change | Old | New |") {
		t.Errorf("expected markdown header, got: %q", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in markdown output, got: %q", out)
	}
}

func TestRender_MarkdownNoChanges(t *testing.T) {
	var buf strings.Builder
	if err := Render(&buf, []Change{}, FormatMarkdown); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "_No changes detected._") {
		t.Errorf("expected markdown no-changes message")
	}
}

func TestRender_UnsupportedFormat(t *testing.T) {
	var buf strings.Builder
	err := Render(&buf, sampleChanges, OutputFormat("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported output format") {
		t.Errorf("unexpected error message: %v", err)
	}
}
