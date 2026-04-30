package vault

import (
	"testing"
)

func TestLintSecret_NoData(t *testing.T) {
	violations := LintSecret("secret/app", map[string]string{}, DefaultLintOptions())
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Rule != RuleNoData {
		t.Errorf("expected rule %s, got %s", RuleNoData, violations[0].Rule)
	}
}

func TestLintSecret_Clean(t *testing.T) {
	data := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "s3cr3t",
	}
	violations := LintSecret("secret/app", data, DefaultLintOptions())
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d: %v", len(violations), violations)
	}
}

func TestLintSecret_EmptyValue(t *testing.T) {
	data := map[string]string{
		"API_KEY": "",
		"HOST":    "example.com",
	}
	violations := LintSecret("secret/app", data, DefaultLintOptions())
	found := false
	for _, v := range violations {
		if v.Rule == RuleEmptyValue && v.Key == "API_KEY" {
			found = true
		}
	}
	if !found {
		t.Error("expected empty-value violation for API_KEY")
	}
}

func TestLintSecret_KeyCasing(t *testing.T) {
	data := map[string]string{
		"db_host": "localhost",
	}
	violations := LintSecret("secret/app", data, DefaultLintOptions())
	if len(violations) != 1 || violations[0].Rule != RuleKeyCase {
		t.Errorf("expected key-case violation, got %v", violations)
	}
}

func TestLintSecret_KeySpaces(t *testing.T) {
	data := map[string]string{
		"MY KEY": "value",
	}
	violations := LintSecret("secret/app", data, DefaultLintOptions())
	found := false
	for _, v := range violations {
		if v.Rule == RuleKeySpaces {
			found = true
		}
	}
	if !found {
		t.Error("expected key-spaces violation")
	}
}

func TestLintSecret_DisabledRules(t *testing.T) {
	data := map[string]string{
		"lower_key": "",
	}
	opts := LintOptions{
		CheckEmptyValues: false,
		CheckKeyCasing:   false,
		CheckKeySpaces:   false,
	}
	violations := LintSecret("secret/app", data, opts)
	if len(violations) != 0 {
		t.Errorf("expected no violations with all rules disabled, got %d", len(violations))
	}
}

func TestLintViolation_String_WithKey(t *testing.T) {
	v := LintViolation{Path: "secret/app", Key: "DB_HOST", Rule: RuleEmptyValue, Message: "empty"}
	s := v.String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
}

func TestLintViolation_String_NoKey(t *testing.T) {
	v := LintViolation{Path: "secret/app", Rule: RuleNoData, Message: "no keys"}
	s := v.String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
}
