package vault

import (
	"context"
	"testing"
)

func TestPromoteSecret_EmptySourcePath(t *testing.T) {
	_, err := PromoteSecret(context.Background(), nil, nil, PromoteOptions{
		SourcePath: "",
		DestPath:   "secret/data/prod/app",
	})
	if err == nil {
		t.Fatal("expected error for empty source path")
	}
}

func TestPromoteSecret_EmptyDestPath(t *testing.T) {
	_, err := PromoteSecret(context.Background(), nil, nil, PromoteOptions{
		SourcePath: "secret/data/staging/app",
		DestPath:   "",
	})
	if err == nil {
		t.Fatal("expected error for empty dest path")
	}
}

func TestPromoteResult_Fields(t *testing.T) {
	r := &PromoteResult{
		SourcePath:  "secret/data/staging/app",
		DestPath:    "secret/data/prod/app",
		SourceEnv:   "staging",
		DestEnv:     "prod",
		KeysWritten: 4,
		DryRun:      true,
	}

	if r.SourcePath != "secret/data/staging/app" {
		t.Errorf("unexpected SourcePath: %s", r.SourcePath)
	}
	if r.DestEnv != "prod" {
		t.Errorf("unexpected DestEnv: %s", r.DestEnv)
	}
	if r.KeysWritten != 4 {
		t.Errorf("unexpected KeysWritten: %d", r.KeysWritten)
	}
	if !r.DryRun {
		t.Error("expected DryRun to be true")
	}
}

func TestPromoteOptions_DefaultVersion(t *testing.T) {
	opts := PromoteOptions{
		SourcePath: "secret/data/staging/app",
		DestPath:   "secret/data/prod/app",
	}
	if opts.Version != 0 {
		t.Errorf("expected default version 0 (latest), got %d", opts.Version)
	}
}
