package diff

import (
	"testing"
)

func TestCompare_NoChanges(t *testing.T) {
	old := map[string]interface{}{"host": "localhost", "port": "5432"}
	new := map[string]interface{}{"host": "localhost", "port": "5432"}

	result := Compare(old, new)
	if result.HasDiff {
		t.Error("expected no diff, but HasDiff is true")
	}
	for _, c := range result.Changes {
		if c.Type != Unchanged {
			t.Errorf("expected all changes to be Unchanged, got %q for key %q", c.Type, c.Key)
		}
	}
}

func TestCompare_AddedKey(t *testing.T) {
	old := map[string]interface{}{"host": "localhost"}
	new := map[string]interface{}{"host": "localhost", "port": "5432"}

	result := Compare(old, new)
	if !result.HasDiff {
		t.Fatal("expected diff, but HasDiff is false")
	}
	found := false
	for _, c := range result.Changes {
		if c.Key == "port" && c.Type == Added {
			found = true
		}
	}
	if !found {
		t.Error("expected 'port' to be marked as Added")
	}
}

func TestCompare_RemovedKey(t *testing.T) {
	old := map[string]interface{}{"host": "localhost", "debug": "true"}
	new := map[string]interface{}{"host": "localhost"}

	result := Compare(old, new)
	if !result.HasDiff {
		t.Fatal("expected diff, but HasDiff is false")
	}
	for _, c := range result.Changes {
		if c.Key == "debug" && c.Type != Removed {
			t.Errorf("expected 'debug' to be Removed, got %q", c.Type)
		}
	}
}

func TestCompare_ModifiedKey(t *testing.T) {
	old := map[string]interface{}{"password": "old-pass"}
	new := map[string]interface{}{"password": "new-pass"}

	result := Compare(old, new)
	if !result.HasDiff {
		t.Fatal("expected diff, but HasDiff is false")
	}
	for _, c := range result.Changes {
		if c.Key == "password" {
			if c.Type != Modified {
				t.Errorf("expected Modified, got %q", c.Type)
			}
			if c.OldValue != "old-pass" || c.NewValue != "new-pass" {
				t.Errorf("unexpected values: old=%v new=%v", c.OldValue, c.NewValue)
			}
		}
	}
}
