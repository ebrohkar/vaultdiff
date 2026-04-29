package vault

import (
	"testing"
	"time"
)

func TestLockOptions_DefaultMount(t *testing.T) {
	opts := LockOptions{
		Path:   "myapp/config",
		Reason: "release freeze",
	}
	if opts.Mount != "" {
		t.Fatalf("expected empty mount, got %q", opts.Mount)
	}
}

func TestLockOptions_DefaultTTL(t *testing.T) {
	opts := LockOptions{Path: "myapp/config"}
	if opts.TTL != 0 {
		t.Fatalf("expected zero TTL before defaulting, got %v", opts.TTL)
	}
}

func TestLockResult_Fields(t *testing.T) {
	now := time.Now().UTC()
	r := LockResult{
		Path:      "myapp/config",
		Mount:     "secret",
		Locked:    true,
		Reason:    "deployment freeze",
		LockedAt:  now,
		ExpiresAt: now.Add(24 * time.Hour),
	}

	if r.Path != "myapp/config" {
		t.Errorf("unexpected Path: %q", r.Path)
	}
	if !r.Locked {
		t.Error("expected Locked to be true")
	}
	if r.Reason != "deployment freeze" {
		t.Errorf("unexpected Reason: %q", r.Reason)
	}
	if r.ExpiresAt.Before(r.LockedAt) {
		t.Error("ExpiresAt should be after LockedAt")
	}
}

func TestLockSecret_EmptyPath(t *testing.T) {
	_, err := LockSecret(nil, LockOptions{Path: ""})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
	if err.Error() != "lock: path must not be empty" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestUnlockSecret_EmptyPath(t *testing.T) {
	_, err := UnlockSecret(nil, "", "secret")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
	if err.Error() != "unlock: path must not be empty" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestLockResult_UnlockedState(t *testing.T) {
	r := LockResult{
		Path:   "myapp/config",
		Mount:  "secret",
		Locked: false,
	}
	if r.Locked {
		t.Error("expected Locked to be false")
	}
	if r.Reason != "" {
		t.Errorf("expected empty Reason, got %q", r.Reason)
	}
}
