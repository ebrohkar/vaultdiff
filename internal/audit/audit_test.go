package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/vaultdiff/internal/audit"
	"github.com/vaultdiff/internal/diff"
)

var sampleChanges = []diff.Change{
	{Key: "DB_PASSWORD", Type: diff.Modified, OldValue: "old", NewValue: "new"},
	{Key: "API_KEY", Type: diff.Added, OldValue: "", NewValue: "abc123"},
}

func TestLogger_RecordJSON(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewLogger(&buf, audit.FormatJSON)

	err := logger.Record("production", "secret/app/config", 1, 2, sampleChanges)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry audit.Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}

	if entry.Environment != "production" {
		t.Errorf("expected environment 'production', got %q", entry.Environment)
	}
	if entry.ChangeCount != 2 {
		t.Errorf("expected change_count 2, got %d", entry.ChangeCount)
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLogger_RecordText(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewLogger(&buf, audit.FormatText)

	err := logger.Record("staging", "secret/app/config", 3, 4, sampleChanges)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for _, substr := range []string{"staging", "secret/app/config", "3->4", "changes=2"} {
		if !strings.Contains(output, substr) {
			t.Errorf("expected output to contain %q, got: %s", substr, output)
		}
	}
}

func TestLogger_RecordUnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewLogger(&buf, audit.Format("xml"))

	err := logger.Record("dev", "secret/app", 1, 2, nil)
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}

func TestEntry_TimestampIsUTC(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewLogger(&buf, audit.FormatJSON)
	_ = logger.Record("dev", "secret/app", 1, 2, nil)

	var entry audit.Entry
	_ = json.Unmarshal(buf.Bytes(), &entry)

	if entry.Timestamp.Location() != time.UTC {
		t.Errorf("expected UTC timestamp, got %v", entry.Timestamp.Location())
	}
}
