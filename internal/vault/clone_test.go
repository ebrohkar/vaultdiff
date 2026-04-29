package vault

import (
	"testing"
	"time"
)

func TestCloneOptions_DefaultVersion(t *testing.T) {
	opts := CloneOptions{
		SourcePath: "secret/data/src",
		DestPath:   "secret/data/dst",
	}
	if opts.Version != 0 {
		t.Errorf("expected default version 0, got %d", opts.Version)
	}
	if opts.Overwrite {
		t.Error("expected Overwrite to default to false")
	}
}

func TestCloneResult_Fields(t *testing.T) {
	now := time.Now().UTC()
	r := CloneResult{
		SourcePath:  "secret/data/src",
		DestPath:    "secret/data/dst",
		Version:     3,
		Overwritten: true,
		ClonedAt:    now,
		Data:        map[string]string{"key": "value"},
	}
	if r.SourcePath != "secret/data/src" {
		t.Errorf("unexpected SourcePath: %s", r.SourcePath)
	}
	if r.DestPath != "secret/data/dst" {
		t.Errorf("unexpected DestPath: %s", r.DestPath)
	}
	if r.Version != 3 {
		t.Errorf("unexpected Version: %d", r.Version)
	}
	if !r.Overwritten {
		t.Error("expected Overwritten to be true")
	}
	if r.Data["key"] != "value" {
		t.Errorf("unexpected Data value: %s", r.Data["key"])
	}
	if r.ClonedAt.Location() != time.UTC {
		t.Error("expected ClonedAt to be UTC")
	}
}

func TestCloneSecret_EmptySourcePath(t *testing.T) {
	_, err := CloneSecret(nil, nil, CloneOptions{DestPath: "secret/data/dst"})
	if err == nil {
		t.Fatal("expected error for empty source path")
	}
}

func TestCloneSecret_EmptyDestPath(t *testing.T) {
	_, err := CloneSecret(nil, nil, CloneOptions{SourcePath: "secret/data/src"})
	if err == nil {
		t.Fatal("expected error for empty destination path")
	}
}

func TestCloneSecret_SamePaths(t *testing.T) {
	_, err := CloneSecret(nil, nil, CloneOptions{
		SourcePath: "secret/data/same",
		DestPath:   "secret/data/same",
	})
	if err == nil {
		t.Fatal("expected error when source and destination are identical")
	}
}

func TestCloneSecret_InvalidVersion(t *testing.T) {
	_, err := CloneSecret(nil, nil, CloneOptions{
		SourcePath: "secret/data/src",
		DestPath:   "secret/data/dst",
		Version:    -1,
	})
	if err == nil {
		t.Fatal("expected error for negative version")
	}
}
