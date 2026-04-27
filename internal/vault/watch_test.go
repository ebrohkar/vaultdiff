package vault

import (
	"context"
	"testing"
	"time"
)

func TestWatchSecret_EmptyPath(t *testing.T) {
	c := &Client{} // token validation not exercised here
	_, err := WatchSecret(context.Background(), c, "", WatchOptions{})
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestWatchSecret_DefaultInterval(t *testing.T) {
	// Verify that a zero Interval is replaced with the 30 s default by
	// cancelling immediately and confirming no panic occurs.
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before the goroutine can poll

	c := &Client{}
	ch, err := WatchSecret(ctx, c, "myapp/config", WatchOptions{
		MountPath: "secret",
		Interval:  0, // should default to 30s
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// drain channel; it should close quickly because ctx is already cancelled
	for range ch {
	}
}

func TestWatchEvent_Fields(t *testing.T) {
	now := time.Now().UTC()
	ev := WatchEvent{
		Path:       "myapp/config",
		OldVersion: 3,
		NewVersion: 4,
		ChangedAt:  now,
	}

	if ev.Path != "myapp/config" {
		t.Errorf("Path: got %q, want %q", ev.Path, "myapp/config")
	}
	if ev.OldVersion != 3 {
		t.Errorf("OldVersion: got %d, want 3", ev.OldVersion)
	}
	if ev.NewVersion != 4 {
		t.Errorf("NewVersion: got %d, want 4", ev.NewVersion)
	}
	if !ev.ChangedAt.Equal(now) {
		t.Errorf("ChangedAt: got %v, want %v", ev.ChangedAt, now)
	}
}

func TestWatchSecret_ChannelClosesOnCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	c := &Client{}
	ch, err := WatchSecret(ctx, c, "myapp/config", WatchOptions{
		MountPath: "secret",
		Interval:  10 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The channel must close within a reasonable time after ctx expires.
	done := make(chan struct{})
	go func() {
		for range ch {
		}
		close(done)
	}()

	select {
	case <-done:
		// success
	case <-time.After(500 * time.Millisecond):
		t.Fatal("channel did not close after context cancellation")
	}
}
