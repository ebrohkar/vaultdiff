package vault

import (
	"testing"
	"time"
)

func TestRenameOptions_DefaultMount(t *testing.T) {
	opts := RenameOptions{
		SourcePath: "app/config",
		DestPath:   "app/settings",
	}
	if opts.Mount != "" {
		t.Errorf("expected empty mount, got %q", opts.Mount)
	}
}

func TestRenameResult_Fields(t *testing.T) {
	now := time.Now().UTC()
	r := RenameResult{
		SourcePath:    "old/path",
		DestPath:      "new/path",
		Version:       3,
		DeletedSource: true,
		Timestamp:     now,
	}
	if r.SourcePath != "old/path" {
		t.Errorf("unexpected SourcePath: %s", r.SourcePath)
	}
	if r.DestPath != "new/path" {
		t.Errorf("unexpected DestPath: %s", r.DestPath)
	}
	if r.Version != 3 {
		t.Errorf("unexpected Version: %d", r.Version)
	}
	if !r.DeletedSource {
		t.Error("expected DeletedSource to be true")
	}
	if r.Timestamp != now {
		t.Error("unexpected Timestamp")
	}
}

func TestRenameSecret_EmptySourcePath(t *testing.T) {
	_, err := RenameSecret(nil, nil, RenameOptions{DestPath: "new/path"})
	if err == nil {
		t.Fatal("expected error for empty source path")
	}
	if err.Error() != "rename: source path must not be empty" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRenameSecret_EmptyDestPath(t *testing.T) {
	_, err := RenameSecret(nil, nil, RenameOptions{SourcePath: "old/path"})
	if err == nil {
		t.Fatal("expected error for empty dest path")
	}
	if err.Error() != "rename: destination path must not be empty" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRenameSecret_SamePaths(t *testing.T) {
	_, err := RenameSecret(nil, nil, RenameOptions{
		SourcePath: "app/config",
		DestPath:   "app/config",
	})
	if err == nil {
		t.Fatal("expected error for identical paths")
	}
	if err.Error() != "rename: source and destination paths must differ" {
		t.Errorf("unexpected error: %v", err)
	}
}
