package vault

import (
	"testing"
	"time"
)

func TestDefaultArchiveOptions(t *testing.T) {
	opts := DefaultArchiveOptions()
	if opts.Mount != "secret" {
		t.Errorf("expected default mount \"secret\", got %q", opts.Mount)
	}
	if opts.Version != 0 {
		t.Errorf("expected default version 0, got %d", opts.Version)
	}
}

func TestArchiveEntry_Fields(t *testing.T) {
	now := time.Now().UTC()
	entry := ArchiveEntry{
		Path:       "prod/db",
		Version:    3,
		Data:       map[string]string{"key": "value"},
		ArchivedAt: now,
	}
	if entry.Path != "prod/db" {
		t.Errorf("unexpected Path: %s", entry.Path)
	}
	if entry.Version != 3 {
		t.Errorf("unexpected Version: %d", entry.Version)
	}
	if entry.Data["key"] != "value" {
		t.Errorf("unexpected Data value: %s", entry.Data["key"])
	}
	if !entry.ArchivedAt.Equal(now) {
		t.Errorf("unexpected ArchivedAt: %v", entry.ArchivedAt)
	}
}

func TestArchiveResult_Fields(t *testing.T) {
	result := ArchiveResult{
		Entry: ArchiveEntry{
			Path:    "staging/api",
			Version: 1,
		},
		AlreadyArchived: true,
	}
	if result.Entry.Path != "staging/api" {
		t.Errorf("unexpected Entry.Path: %s", result.Entry.Path)
	}
	if !result.AlreadyArchived {
		t.Error("expected AlreadyArchived to be true")
	}
}

func TestArchiveSecret_EmptyPath(t *testing.T) {
	client := &Client{}
	_, err := ArchiveSecret(client, "", DefaultArchiveOptions())
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestArchiveSecret_NilClient(t *testing.T) {
	_, err := ArchiveSecret(nil, "prod/db", DefaultArchiveOptions())
	if err == nil {
		t.Fatal("expected error for nil client, got nil")
	}
}

func TestArchiveEntry_ArchivedAtIsUTC(t *testing.T) {
	entry := ArchiveEntry{
		ArchivedAt: time.Now().UTC(),
	}
	if entry.ArchivedAt.Location() != time.UTC {
		t.Errorf("expected UTC location, got %v", entry.ArchivedAt.Location())
	}
}
