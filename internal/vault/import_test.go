package vault

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/vault/api"
)

func TestDefaultImportOptions(t *testing.T) {
	opts := DefaultImportOptions()
	if opts.Mount != "secret" {
		t.Errorf("expected mount %q, got %q", "secret", opts.Mount)
	}
	if opts.Format != ImportFormatJSON {
		t.Errorf("expected format %q, got %q", ImportFormatJSON, opts.Format)
	}
}

func TestImportSecret_EmptyPath(t *testing.T) {
	_, err := ImportSecret(ImportOptions{Client: &api.Client{}, FilePath: "x.json"})
	if err == nil || err.Error() != "import: path must not be empty" {
		t.Errorf("expected empty path error, got %v", err)
	}
}

func TestImportSecret_EmptyFilePath(t *testing.T) {
	_, err := ImportSecret(ImportOptions{Client: &api.Client{}, Path: "myapp/config"})
	if err == nil || err.Error() != "import: file path must not be empty" {
		t.Errorf("expected empty file path error, got %v", err)
	}
}

func TestImportSecret_NilClient(t *testing.T) {
	_, err := ImportSecret(ImportOptions{Path: "myapp/config", FilePath: "x.json"})
	if err == nil || err.Error() != "import: vault client must not be nil" {
		t.Errorf("expected nil client error, got %v", err)
	}
}

func TestImportSecret_InvalidFormat(t *testing.T) {
	f := writeTempFile(t, "test.txt", []byte("key=val"))
	_, err := ImportSecret(ImportOptions{
		Client:   &api.Client{},
		Path:     "myapp/config",
		FilePath: f,
		Format:   "toml",
	})
	if err == nil {
		t.Error("expected unsupported format error")
	}
}

func TestImportSecret_DryRunJSON(t *testing.T) {
	data := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	b, _ := json.Marshal(data)
	f := writeTempFile(t, "secrets.json", b)

	result, err := ImportSecret(ImportOptions{
		Client:   &api.Client{},
		Path:     "myapp/db",
		Mount:    "secret",
		FilePath: f,
		Format:   ImportFormatJSON,
		DryRun:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.KeyCount != 2 {
		t.Errorf("expected 2 keys, got %d", result.KeyCount)
	}
	if !result.DryRun {
		t.Error("expected DryRun to be true")
	}
	if result.Data["DB_HOST"] != "localhost" {
		t.Errorf("unexpected DB_HOST value: %q", result.Data["DB_HOST"])
	}
}

func TestImportSecret_DryRunEnv(t *testing.T) {
	content := "# comment\nAPP_ENV=production\nAPP_PORT=8080\n"
	f := writeTempFile(t, "secrets.env", []byte(content))

	result, err := ImportSecret(ImportOptions{
		Client:   &api.Client{},
		Path:     "myapp/env",
		Mount:    "secret",
		FilePath: f,
		Format:   ImportFormatEnv,
		DryRun:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.KeyCount != 2 {
		t.Errorf("expected 2 keys, got %d", result.KeyCount)
	}
	if result.Data["APP_ENV"] != "production" {
		t.Errorf("unexpected APP_ENV: %q", result.Data["APP_ENV"])
	}
}

func TestParseEnvFile_InvalidLine(t *testing.T) {
	_, err := parseEnvFile("BADLINE")
	if err == nil {
		t.Error("expected error for invalid env line")
	}
}

func writeTempFile(t *testing.T, name string, content []byte) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, content, 0600); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return p
}
