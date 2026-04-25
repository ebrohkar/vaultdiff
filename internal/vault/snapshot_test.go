package vault

import (
	"testing"
)

func TestDiffSnapshots_BothNil(t *testing.T) {
	_, _, err := DiffSnapshots(nil, nil)
	if err == nil {
		t.Fatal("expected error for nil snapshots, got nil")
	}
}

func TestDiffSnapshots_OneNil(t *testing.T) {
	a := &SecretSnapshot{Data: map[string]string{"k": "v"}}
	_, _, err := DiffSnapshots(a, nil)
	if err == nil {
		t.Fatal("expected error when second snapshot is nil")
	}
}

func TestDiffSnapshots_ReturnsMaps(t *testing.T) {
	a := &SecretSnapshot{
		Path:        "secret/data/app",
		Environment: "staging",
		Version:     1,
		Data:        map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"},
	}
	b := &SecretSnapshot{
		Path:        "secret/data/app",
		Environment: "production",
		Version:     2,
		Data:        map[string]string{"DB_HOST": "prod-db.internal", "DB_PORT": "5432"},
	}

	mapA, mapB, err := DiffSnapshots(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mapA["DB_HOST"] != "localhost" {
		t.Errorf("expected mapA DB_HOST=localhost, got %q", mapA["DB_HOST"])
	}
	if mapB["DB_HOST"] != "prod-db.internal" {
		t.Errorf("expected mapB DB_HOST=prod-db.internal, got %q", mapB["DB_HOST"])
	}
}

func TestTakeSnapshot_EmptyPath(t *testing.T) {
	c := &Client{} // zero-value client; path validation fires before any network call
	_, err := TakeSnapshot(nil, c, "", "staging", 1)
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestSecretSnapshot_Fields(t *testing.T) {
	s := &SecretSnapshot{
		Path:        "secret/data/myapp",
		Version:     3,
		Environment: "dev",
		Data:        map[string]string{"KEY": "value"},
	}
	if s.Path != "secret/data/myapp" {
		t.Errorf("unexpected path: %s", s.Path)
	}
	if s.Version != 3 {
		t.Errorf("unexpected version: %d", s.Version)
	}
	if s.Data["KEY"] != "value" {
		t.Errorf("unexpected data value: %s", s.Data["KEY"])
	}
}
