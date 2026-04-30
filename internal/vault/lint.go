package vault

import (
	"fmt"
	"strings"
)

// LintRule represents a single validation rule applied to a secret's data.
type LintRule string

const (
	RuleEmptyValue   LintRule = "empty-value"
	RuleKeyCase      LintRule = "key-case"
	RuleKeySpaces    LintRule = "key-spaces"
	RuleNoData       LintRule = "no-data"
)

// LintViolation describes a rule violation found in a secret.
type LintViolation struct {
	Path    string
	Key     string
	Rule    LintRule
	Message string
}

func (v LintViolation) String() string {
	if v.Key == "" {
		return fmt.Sprintf("[%s] %s: %s", v.Rule, v.Path, v.Message)
	}
	return fmt.Sprintf("[%s] %s#%s: %s", v.Rule, v.Path, v.Key, v.Message)
}

// LintOptions configures which rules are applied during linting.
type LintOptions struct {
	CheckEmptyValues bool
	CheckKeyCasing   bool
	CheckKeySpaces   bool
}

// DefaultLintOptions returns a LintOptions with all checks enabled.
func DefaultLintOptions() LintOptions {
	return LintOptions{
		CheckEmptyValues: true,
		CheckKeyCasing:   true,
		CheckKeySpaces:   true,
	}
}

// LintSecret validates the data map at the given path against the configured rules.
// It returns a slice of violations (empty if the secret is clean).
func LintSecret(path string, data map[string]string, opts LintOptions) []LintViolation {
	var violations []LintViolation

	if len(data) == 0 {
		violations = append(violations, LintViolation{
			Path:    path,
			Rule:    RuleNoData,
			Message: "secret has no keys",
		})
		return violations
	}

	for k, v := range data {
		if opts.CheckEmptyValues && strings.TrimSpace(v) == "" {
			violations = append(violations, LintViolation{
				Path:    path,
				Key:     k,
				Rule:    RuleEmptyValue,
				Message: "value is empty or whitespace-only",
			})
		}
		if opts.CheckKeySpaces && strings.Contains(k, " ") {
			violations = append(violations, LintViolation{
				Path:    path,
				Key:     k,
				Rule:    RuleKeySpaces,
				Message: "key contains spaces",
			})
		}
		if opts.CheckKeyCasing && k != strings.ToUpper(k) {
			violations = append(violations, LintViolation{
				Path:    path,
				Key:     k,
				Rule:    RuleKeyCase,
				Message: fmt.Sprintf("key %q is not upper-case", k),
			})
		}
	}
	return violations
}
