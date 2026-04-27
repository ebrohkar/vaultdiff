package vault

import (
	"testing"
	"time"
)

func TestRollbackResult_Fields(t *testing.T) {
	now := time.Now().UTC()
	r := &RollbackResult{
		Path:         "secret/app/config",
		FromVersion:  5,
		ToVersion:    3,
		DryRun:       false,
		RolledBackAt: now,
		Data:         map[string]interface{}{"key": "value"},
	}

	if r.Path != "secret/app/config" {
		t.Errorf("expected path 'secret/app/config', got %q", r.Path)
	}
	if r.FromVersion != 5 {
		t.Errorf("expected FromVersion 5, got %d", r.FromVersion)
	}
	if r.ToVersion != 3 {
		t.Errorf("expected ToVersion 3, got %d", r.ToVersion)
	}
	if r.DryRun {
		t.Error("expected DryRun to be false")
	}
	if r.RolledBackAt.IsZero() {
		t.Error("expected RolledBackAt to be set")
	}
}

func TestRollbackOptions_DefaultVersion(t *testing.T) {
	opts := RollbackOptions{
		Path:          "secret/app/config",
		TargetVersion: 0,
	}
	if opts.TargetVersion != 0 {
		t.Errorf("expected zero default, got %d", opts.TargetVersion)
	}
}

func TestRollbackSecret_EmptyPath(t *testing.T) {
	_, err := RollbackSecret(nil, nil, RollbackOptions{Path: "", TargetVersion: 1})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
	if err.Error() != "rollback: path must not be empty" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestRollbackSecret_InvalidVersion(t *testing.T) {
	_, err := RollbackSecret(nil, nil, RollbackOptions{Path: "secret/app", TargetVersion: 0})
	if err == nil {
		t.Fatal("expected error for version < 1")
	}
	if err.Error() != "rollback: target version must be >= 1" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestSplitMountAndPath(t *testing.T) {
	tests := []struct {
		input     string
		expMount  string
		expSecret string
	}{
		{"secret/app/config", "secret", "app/config"},
		{"kv/myapp", "kv", "myapp"},
		{"noseparator", "noseparator", ""},
	}
	for _, tt := range tests {
		m, s := splitMountAndPath(tt.input)
		if m != tt.expMount {
			t.Errorf("input %q: expected mount %q, got %q", tt.input, tt.expMount, m)
		}
		if s != tt.expSecret {
			t.Errorf("input %q: expected secret path %q, got %q", tt.input, tt.expSecret, s)
		}
	}
}
