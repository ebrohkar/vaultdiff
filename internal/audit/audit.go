package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/vaultdiff/internal/diff"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp   time.Time        `json:"timestamp"`
	Environment string           `json:"environment"`
	Path        string           `json:"path"`
	FromVersion int              `json:"from_version"`
	ToVersion   int              `json:"to_version"`
	Changes     []diff.Change    `json:"changes"`
	ChangeCount int              `json:"change_count"`
}

// Logger writes audit entries to an io.Writer.
type Logger struct {
	w       io.Writer
	format  Format
}

// Format controls the output format of audit entries.
type Format string

const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

// NewLogger creates a new audit Logger writing to w.
func NewLogger(w io.Writer, format Format) *Logger {
	return &Logger{w: w, format: format}
}

// Record writes an audit entry for the given diff result.
func (l *Logger) Record(env, path string, fromVersion, toVersion int, changes []diff.Change) error {
	entry := Entry{
		Timestamp:   time.Now().UTC(),
		Environment: env,
		Path:        path,
		FromVersion: fromVersion,
		ToVersion:   toVersion,
		Changes:     changes,
		ChangeCount: len(changes),
	}

	switch l.format {
	case FormatJSON:
		return l.writeJSON(entry)
	case FormatText:
		return l.writeText(entry)
	default:
		return fmt.Errorf("unsupported audit format: %s", l.format)
	}
}

func (l *Logger) writeJSON(entry Entry) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: failed to marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", data)
	return err
}

func (l *Logger) writeText(entry Entry) error {
	_, err := fmt.Fprintf(l.w,
		"[%s] env=%s path=%s versions=%d->%d changes=%d\n",
		entry.Timestamp.Format(time.RFC3339),
		entry.Environment,
		entry.Path,
		entry.FromVersion,
		entry.ToVersion,
		entry.ChangeCount,
	)
	return err
}
