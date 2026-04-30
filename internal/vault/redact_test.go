package vault

import (
	"testing"
)

func TestDefaultRedactOptions(t *testing.T) {
	opts := DefaultRedactOptions()
	if opts.Placeholder != "[REDACTED]" {
		t.Errorf("expected placeholder [REDACTED], got %q", opts.Placeholder)
	}
	if len(opts.Patterns) == 0 {
		t.Error("expected non-empty default patterns")
	}
}

func TestRedactSecret_NoSensitiveKeys(t *testing.T) {
	data := map[string]string{
		"host": "localhost",
		"port": "5432",
	}
	result := RedactSecret(data, DefaultRedactOptions())
	if result.Data["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %q", result.Data["host"])
	}
	if len(result.RedactedKeys) != 0 {
		t.Errorf("expected no redacted keys, got %v", result.RedactedKeys)
	}
}

func TestRedactSecret_SensitiveKeys(t *testing.T) {
	data := map[string]string{
		"db_password": "s3cr3t",
		"api_token":   "abc123",
		"username":    "admin",
	}
	result := RedactSecret(data, DefaultRedactOptions())
	if result.Data["db_password"] != "[REDACTED]" {
		t.Errorf("expected db_password to be redacted, got %q", result.Data["db_password"])
	}
	if result.Data["api_token"] != "[REDACTED]" {
		t.Errorf("expected api_token to be redacted, got %q", result.Data["api_token"])
	}
	if result.Data["username"] != "admin" {
		t.Errorf("expected username=admin, got %q", result.Data["username"])
	}
	if len(result.RedactedKeys) != 2 {
		t.Errorf("expected 2 redacted keys, got %d", len(result.RedactedKeys))
	}
}

func TestRedactSecret_CustomPlaceholder(t *testing.T) {
	data := map[string]string{"secret_key": "topsecret"}
	opts := RedactOptions{
		Patterns:    []string{"secret"},
		Placeholder: "***",
	}
	result := RedactSecret(data, opts)
	if result.Data["secret_key"] != "***" {
		t.Errorf("expected ***, got %q", result.Data["secret_key"])
	}
}

func TestRedactSecret_EmptyPlaceholderFallback(t *testing.T) {
	data := map[string]string{"password": "hunter2"}
	opts := RedactOptions{
		Patterns:    []string{"password"},
		Placeholder: "",
	}
	result := RedactSecret(data, opts)
	if result.Data["password"] != "[REDACTED]" {
		t.Errorf("expected fallback placeholder, got %q", result.Data["password"])
	}
}

func TestMatchesAnyPattern_CaseInsensitive(t *testing.T) {
	if !matchesAnyPattern("DB_PASSWORD", []string{"password"}) {
		t.Error("expected DB_PASSWORD to match pattern 'password'")
	}
	if matchesAnyPattern("hostname", []string{"password", "secret"}) {
		t.Error("expected hostname not to match any sensitive pattern")
	}
}

func TestRedactSecret_EmptyData(t *testing.T) {
	result := RedactSecret(map[string]string{}, DefaultRedactOptions())
	if len(result.Data) != 0 {
		t.Errorf("expected empty result, got %v", result.Data)
	}
	if len(result.RedactedKeys) != 0 {
		t.Errorf("expected no redacted keys, got %v", result.RedactedKeys)
	}
}
