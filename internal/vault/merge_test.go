package vault

import "testing"

func TestDefaultMergeOptions(t *testing.T) {
	opts := DefaultMergeOptions()
	if opts.Mount != "secret" {
		t.Errorf("expected mount=secret, got %q", opts.Mount)
	}
	if opts.Strategy != "source" {
		t.Errorf("expected strategy=source, got %q", opts.Strategy)
	}
}

func TestMergeResult_Fields(t *testing.T) {
	r := MergeResult{
		SourcePath:  "env/dev/app",
		DestPath:    "env/prod/app",
		Merged:      map[string]string{"KEY": "val"},
		Conflicts:   []string{"KEY"},
		AddedKeys:   []string{"NEW"},
		SkippedKeys: []string{"OLD"},
	}
	if r.SourcePath != "env/dev/app" {
		t.Errorf("unexpected SourcePath: %s", r.SourcePath)
	}
	if len(r.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(r.Conflicts))
	}
}

func TestMergeSecret_EmptySourcePath(t *testing.T) {
	_, err := MergeSecret(&Client{}, MergeOptions{DestPath: "a", Strategy: "source"})
	if err == nil || err.Error() != "merge: source path must not be empty" {
		t.Errorf("expected source path error, got %v", err)
	}
}

func TestMergeSecret_EmptyDestPath(t *testing.T) {
	_, err := MergeSecret(&Client{}, MergeOptions{SourcePath: "a", Strategy: "source"})
	if err == nil || err.Error() != "merge: dest path must not be empty" {
		t.Errorf("expected dest path error, got %v", err)
	}
}

func TestMergeSecret_SamePaths(t *testing.T) {
	_, err := MergeSecret(&Client{}, MergeOptions{SourcePath: "a", DestPath: "a", Strategy: "source"})
	if err == nil {
		t.Error("expected error for same source and dest paths")
	}
}

func TestMergeSecret_NilClient(t *testing.T) {
	_, err := MergeSecret(nil, MergeOptions{SourcePath: "a", DestPath: "b", Strategy: "source"})
	if err == nil {
		t.Error("expected error for nil client")
	}
}

func TestMergeSecret_InvalidStrategy(t *testing.T) {
	_, err := MergeSecret(&Client{}, MergeOptions{SourcePath: "a", DestPath: "b", Strategy: "unknown"})
	if err == nil {
		t.Error("expected error for unknown strategy")
	}
}

func TestMergeOptions_DefaultVersion(t *testing.T) {
	opts := DefaultMergeOptions()
	if opts.SourcePath != "" {
		t.Errorf("expected empty SourcePath by default, got %q", opts.SourcePath)
	}
	if opts.DestPath != "" {
		t.Errorf("expected empty DestPath by default, got %q", opts.DestPath)
	}
}
