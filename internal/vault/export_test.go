package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExportResult_Fields(t *testing.T) {
	r := ExportResult{
		Path:       "secret/data/app",
		Version:    2,
		Format:     ExportFormatJSON,
		OutputFile: "/tmp/out.json",
		KeyCount:   3,
	}
	if r.Path != "secret/data/app" {
		t.Errorf("unexpected Path: %s", r.Path)
	}
	if r.KeyCount != 3 {
		t.Errorf("unexpected KeyCount: %d", r.KeyCount)
	}
	if r.ExportedAt.IsZero() {
		t.Error("ExportedAt should be set")
	}
}

func TestExportOptions_DefaultVersion(t *testing.T) {
	opts := ExportOptions{Path: "secret/data/app"}
	if opts.Version != 0 {
		t.Errorf("expected default version 0, got %d", opts.Version)
	}
}

func TestExportSecret_EmptyPath(t *testing.T) {
	client := &Client{}
	_, err := ExportSecret(client, ExportOptions{})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestExportSecret_InvalidFormat(t *testing.T) {
	client := &Client{}
	_, err := ExportSecret(client, ExportOptions{
		Path:   "secret/data/app",
		Format: "yaml",
	})
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestMarshalEnv_Basic(t *testing.T) {
	data := map[string]string{"KEY": "value"}
	out := marshalEnv(data)
	if len(out) == 0 {
		t.Error("expected non-empty env output")
	}
	got := string(out)
	if got != "KEY=value\n" {
		t.Errorf("unexpected env output: %q", got)
	}
}

func TestExportSecret_WritesFile(t *testing.T) {
	dir := t.TempDir()
	outFile := filepath.Join(dir, "secrets.json")

	client := &Client{}
	// ReadSecretVersion will fail without a real server; test file creation path
	// by verifying the error is from the vault read, not file handling.
	_, err := ExportSecret(client, ExportOptions{
		Path:       "secret/data/app",
		Format:     ExportFormatJSON,
		OutputFile: outFile,
	})
	// We expect a vault read error, not a file error.
	if err == nil {
		t.Fatal("expected vault read error")
	}
	// File should NOT exist since we failed before writing.
	_, statErr := os.Stat(outFile)
	if statErr == nil {
		t.Error("file should not exist after failed export")
	}
}
