package vault

import (
	"testing"
)

func TestDefaultCopyOptions(t *testing.T) {
	opts := DefaultCopyOptions()
	if opts.Mount != "secret" {
		t.Errorf("expected default mount \"secret\", got %q", opts.Mount)
	}
	if opts.Version != 0 {
		t.Errorf("expected default version 0, got %d", opts.Version)
	}
}

func TestCopyResult_Fields(t *testing.T) {
	r := CopyResult{
		SourcePath:  "src/app",
		DestPath:    "dst/app",
		Mount:       "secret",
		Version:     3,
		Overwritten: true,
	}
	if r.SourcePath != "src/app" {
		t.Errorf("unexpected SourcePath: %s", r.SourcePath)
	}
	if r.DestPath != "dst/app" {
		t.Errorf("unexpected DestPath: %s", r.DestPath)
	}
	if r.Mount != "secret" {
		t.Errorf("unexpected Mount: %s", r.Mount)
	}
	if r.Version != 3 {
		t.Errorf("unexpected Version: %d", r.Version)
	}
	if !r.Overwritten {
		t.Error("expected Overwritten to be true")
	}
}

func TestCopySecret_EmptySourcePath(t *testing.T) {
	opts := DefaultCopyOptions()
	opts.DestPath = "dst/app"
	_, err := CopySecret(&Client{}, opts)
	if err == nil {
		t.Fatal("expected error for empty source path")
	}
}

func TestCopySecret_EmptyDestPath(t *testing.T) {
	opts := DefaultCopyOptions()
	opts.SourcePath = "src/app"
	_, err := CopySecret(&Client{}, opts)
	if err == nil {
		t.Fatal("expected error for empty dest path")
	}
}

func TestCopySecret_SamePaths(t *testing.T) {
	opts := DefaultCopyOptions()
	opts.SourcePath = "app/config"
	opts.DestPath = "app/config"
	_, err := CopySecret(&Client{}, opts)
	if err == nil {
		t.Fatal("expected error when source and dest are identical")
	}
}

func TestCopySecret_NilClient(t *testing.T) {
	opts := DefaultCopyOptions()
	opts.SourcePath = "src/app"
	opts.DestPath = "dst/app"
	_, err := CopySecret(nil, opts)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}
