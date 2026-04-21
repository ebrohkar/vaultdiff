package audit_test

import (
	"testing"
	"time"

	"github.com/vaultdiff/internal/audit"
)

func makeEntry(env, path string, ts time.Time, changeCount int) audit.Entry {
	return audit.Entry{
		Timestamp:   ts,
		Environment: env,
		Path:        path,
		ChangeCount: changeCount,
	}
}

func TestFilter_ByEnvironment(t *testing.T) {
	now := time.Now().UTC()
	entries := []audit.Entry{
		makeEntry("production", "secret/app", now, 1),
		makeEntry("staging", "secret/app", now, 1),
	}

	result := audit.ApplyFilter(entries, audit.Filter{Environment: "production"})
	if len(result) != 1 || result[0].Environment != "production" {
		t.Errorf("expected 1 production entry, got %d", len(result))
	}
}

func TestFilter_ByPathPrefix(t *testing.T) {
	now := time.Now().UTC()
	entries := []audit.Entry{
		makeEntry("prod", "secret/app/config", now, 1),
		makeEntry("prod", "secret/db/creds", now, 1),
	}

	result := audit.ApplyFilter(entries, audit.Filter{PathPrefix: "secret/app"})
	if len(result) != 1 {
		t.Errorf("expected 1 entry with path prefix 'secret/app', got %d", len(result))
	}
}

func TestFilter_BySince(t *testing.T) {
	now := time.Now().UTC()
	old := now.Add(-48 * time.Hour)

	entries := []audit.Entry{
		makeEntry("prod", "secret/app", old, 1),
		makeEntry("prod", "secret/app", now, 1),
	}

	result := audit.ApplyFilter(entries, audit.Filter{Since: now.Add(-1 * time.Hour)})
	if len(result) != 1 {
		t.Errorf("expected 1 recent entry, got %d", len(result))
	}
}

func TestFilter_ByMinChanges(t *testing.T) {
	now := time.Now().UTC()
	entries := []audit.Entry{
		makeEntry("prod", "secret/a", now, 1),
		makeEntry("prod", "secret/b", now, 5),
	}

	result := audit.ApplyFilter(entries, audit.Filter{MinChanges: 3})
	if len(result) != 1 || result[0].Path != "secret/b" {
		t.Errorf("expected 1 entry with >= 3 changes, got %d", len(result))
	}
}

func TestFilter_NoFilter_ReturnsAll(t *testing.T) {
	now := time.Now().UTC()
	entries := []audit.Entry{
		makeEntry("prod", "secret/a", now, 1),
		makeEntry("staging", "secret/b", now, 2),
	}

	result := audit.ApplyFilter(entries, audit.Filter{})
	if len(result) != 2 {
		t.Errorf("expected all 2 entries, got %d", len(result))
	}
}
