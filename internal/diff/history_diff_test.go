package diff

import (
	"strings"
	"testing"
)

func TestDiffHistory_TooFewVersions(t *testing.T) {
	_, err := DiffHistory([]map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty versions slice")
	}

	_, err = DiffHistory([]map[string]string{{"a": "1"}})
	if err == nil {
		t.Fatal("expected error for single version")
	}
}

func TestDiffHistory_NoChanges(t *testing.T) {
	v1 := map[string]string{"key": "value"}
	v2 := map[string]string{"key": "value"}

	diffs, err := DiffHistory([]map[string]string{v1, v2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if len(diffs[0].Changes) != 0 {
		t.Errorf("expected no changes, got %d", len(diffs[0].Changes))
	}
}

func TestDiffHistory_WithChanges(t *testing.T) {
	v1 := map[string]string{"a": "1", "b": "old"}
	v2 := map[string]string{"a": "1", "b": "new", "c": "added"}
	v3 := map[string]string{"a": "1"}

	diffs, err := DiffHistory([]map[string]string{v1, v2, v3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diffs) != 2 {
		t.Fatalf("expected 2 diffs, got %d", len(diffs))
	}
	if diffs[0].FromVersion != 1 || diffs[0].ToVersion != 2 {
		t.Errorf("unexpected version range: %d->%d", diffs[0].FromVersion, diffs[0].ToVersion)
	}
}

func TestSummariseHistory_Empty(t *testing.T) {
	out := SummariseHistory(nil)
	if out != "No history diffs available." {
		t.Errorf("unexpected summary: %q", out)
	}
}

func TestSummariseHistory_WithDiffs(t *testing.T) {
	diffs := []VersionedDiff{
		{
			FromVersion: 1,
			ToVersion:   2,
			Changes: []Change{
				{Key: "password", Type: ChangeModified, OldValue: "old", NewValue: "new"},
			},
		},
	}
	out := SummariseHistory(diffs)
	if !strings.Contains(out, "v1 -> v2") {
		t.Errorf("expected version range in summary, got: %s", out)
	}
	if !strings.Contains(out, "password") {
		t.Errorf("expected key name in summary, got: %s", out)
	}
}
