package config_test

import (
	"testing"

	"github.com/your-org/vaultdiff/internal/config"
)

func setEnv(t *testing.T, pairs map[string]string) {
	t.Helper()
	for k, v := range pairs {
		t.Setenv(k, v)
	}
}

func TestLoad_MissingToken(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "")
	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error when VAULT_TOKEN is missing")
	}
}

func TestLoad_Defaults(t *testing.T) {
	setEnv(t, map[string]string{
		"VAULT_TOKEN":  "test-token",
		"VAULT_ADDR":   "",
		"VAULTDIFF_FORMAT": "",
	})
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VaultAddr != "http://127.0.0.1:8200" {
		t.Errorf("expected default addr, got %q", cfg.VaultAddr)
	}
	if cfg.OutputFormat != "text" {
		t.Errorf("expected default format 'text', got %q", cfg.OutputFormat)
	}
}

func TestLoad_InvalidFormat(t *testing.T) {
	setEnv(t, map[string]string{
		"VAULT_TOKEN":      "test-token",
		"VAULTDIFF_FORMAT": "yaml",
	})
	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestLoad_Environments(t *testing.T) {
	setEnv(t, map[string]string{
		"VAULT_TOKEN":           "test-token",
		"VAULTDIFF_ENVIRONMENTS": "staging, production, dev",
	})
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Environments) != 3 {
		t.Errorf("expected 3 environments, got %d", len(cfg.Environments))
	}
	if cfg.Environments[1] != "production" {
		t.Errorf("expected 'production', got %q", cfg.Environments[1])
	}
}

func TestValidate_AllFormats(t *testing.T) {
	formats := []string{"text", "json", "markdown"}
	for _, f := range formats {
		cfg := &config.Config{
			VaultAddr:    "http://vault:8200",
			VaultToken:   "tok",
			OutputFormat: f,
		}
		if err := cfg.Validate(); err != nil {
			t.Errorf("format %q should be valid, got error: %v", f, err)
		}
	}
}
