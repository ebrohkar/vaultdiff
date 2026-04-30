package vault

import (
	"testing"
)

func TestDefaultValidateOptions(t *testing.T) {
	opts := DefaultValidateOptions()
	if opts.Mount != "secret" {
		t.Errorf("expected mount=secret, got %q", opts.Mount)
	}
	if len(opts.Rules) != 0 {
		t.Errorf("expected empty rules, got %d", len(opts.Rules))
	}
}

func TestValidationResult_Fields(t *testing.T) {
	r := ValidationResult{
		Path:       "secret/myapp",
		Passed:     false,
		Violations: []string{"key \"foo\" missing"},
	}
	if r.Path != "secret/myapp" {
		t.Errorf("unexpected path: %q", r.Path)
	}
	if r.Passed {
		t.Error("expected Passed=false")
	}
	if len(r.Violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(r.Violations))
	}
}

func TestValidateSecret_EmptyPath(t *testing.T) {
	_, err := ValidateSecret(&Client{}, "", DefaultValidateOptions())
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestValidateSecret_NilClient(t *testing.T) {
	_, err := ValidateSecret(nil, "secret/myapp", DefaultValidateOptions())
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestMatchPattern_ValidMatch(t *testing.T) {
	if err := matchPattern("key", "hello123", `^[a-z0-9]+$`); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMatchPattern_NoMatch(t *testing.T) {
	if err := matchPattern("key", "HELLO", `^[a-z]+$`); err == nil {
		t.Error("expected error for non-matching value")
	}
}

func TestMatchPattern_InvalidRegex(t *testing.T) {
	if err := matchPattern("key", "value", `[invalid`); err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestTruncate_Short(t *testing.T) {
	out := truncate("hello", 10)
	if out != "hello" {
		t.Errorf("unexpected truncation: %q", out)
	}
}

func TestTruncate_Long(t *testing.T) {
	out := truncate("abcdefghijklmnopqrstuvwxyz", 10)
	if len(out) <= 10 {
		// truncated with ellipsis
		return
	}
	// if longer than original truncation limit + ellipsis it's a bug
	if len(out) > 14 {
		t.Errorf("truncated value too long: %q", out)
	}
}

func TestValidationRule_RequiredField(t *testing.T) {
	rule := ValidationRule{
		Key:      "api_key",
		Required: true,
		Pattern:  `^[A-Za-z0-9]+$`,
	}
	if rule.Key != "api_key" {
		t.Errorf("unexpected key: %q", rule.Key)
	}
	if !rule.Required {
		t.Error("expected Required=true")
	}
}
