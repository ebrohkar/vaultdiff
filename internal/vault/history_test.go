package vault

import (
	"testing"
	"time"
)

func TestSecretHistoryEntry_Fields(t *testing.T) {
	now := time.Now().UTC()
	entry := SecretHistoryEntry{
		Version:     3,
		CreatedTime: now,
		Destroyed:   false,
		Data:        map[string]interface{}{"key": "value"},
	}

	if entry.Version != 3 {
		t.Errorf("expected version 3, got %d", entry.Version)
	}
	if entry.Destroyed {
		t.Error("expected entry not to be destroyed")
	}
	if entry.DeletedTime != nil {
		t.Error("expected DeletedTime to be nil")
	}
	if entry.Data["key"] != "value" {
		t.Errorf("expected data key=value, got %v", entry.Data["key"])
	}
}

func TestSecretHistory_Fields(t *testing.T) {
	h := &SecretHistory{
		Path: "secret/myapp",
		Entries: []SecretHistoryEntry{
			{Version: 1},
			{Version: 2},
		},
	}

	if h.Path != "secret/myapp" {
		t.Errorf("unexpected path: %s", h.Path)
	}
	if len(h.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(h.Entries))
	}
}

func TestFetchHistory_EmptyPath(t *testing.T) {
	c := &Client{token: "tok", address: "http://127.0.0.1:8200"}
	_, err := FetchHistory(nil, c, "secret", "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestSecretHistoryEntry_WithDeletedTime(t *testing.T) {
	now := time.Now().UTC()
	entry := SecretHistoryEntry{
		Version:     5,
		CreatedTime: now,
		DeletedTime: &now,
		Destroyed:   false,
	}

	if entry.DeletedTime == nil {
		t.Error("expected DeletedTime to be set")
	}
	if entry.Destroyed {
		t.Error("expected entry not to be destroyed")
	}
}
