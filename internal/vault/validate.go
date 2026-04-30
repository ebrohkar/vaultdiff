package vault

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ValidationRule defines a single rule applied to a secret's key-value data.
type ValidationRule struct {
	Key     string // empty means apply to all keys
	Pattern string // regex pattern the value must match
	Required bool   // if true, key must be present
}

// ValidationResult holds the outcome of validating a secret.
type ValidationResult struct {
	Path     string
	Passed   bool
	Violations []string
}

// ValidateOptions configures secret validation behaviour.
type ValidateOptions struct {
	Mount string
	Rules []ValidationRule
}

// DefaultValidateOptions returns sensible defaults.
func DefaultValidateOptions() ValidateOptions {
	return ValidateOptions{
		Mount: "secret",
		Rules: []ValidationRule{},
	}
}

// ValidateSecret checks a secret's data against the provided rules.
func ValidateSecret(client *Client, path string, opts ValidateOptions) (*ValidationResult, error) {
	if path == "" {
		return nil, errors.New("path must not be empty")
	}
	if client == nil {
		return nil, errors.New("client must not be nil")
	}

	data, err := client.ReadSecretVersion(path, 0)
	if err != nil {
		return nil, fmt.Errorf("reading secret %q: %w", path, err)
	}

	result := &ValidationResult{Path: path, Passed: true}

	for _, rule := range opts.Rules {
		if rule.Key != "" {
			val, ok := data[rule.Key]
			if !ok {
				if rule.Required {
					result.Violations = append(result.Violations,
						fmt.Sprintf("required key %q is missing", rule.Key))
				}
				continue
			}
			if rule.Pattern != "" {
				if err := matchPattern(rule.Key, fmt.Sprintf("%v", val), rule.Pattern); err != nil {
					result.Violations = append(result.Violations, err.Error())
				}
			}
		} else if rule.Pattern != "" {
			for k, v := range data {
				if err := matchPattern(k, fmt.Sprintf("%v", v), rule.Pattern); err != nil {
					result.Violations = append(result.Violations, err.Error())
				}
			}
		}
	}

	if len(result.Violations) > 0 {
		result.Passed = false
	}
	return result, nil
}

func matchPattern(key, value, pattern string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid pattern %q: %w", pattern, err)
	}
	if !re.MatchString(value) {
		return fmt.Errorf("key %q value %q does not match pattern %q",
			key, truncate(value, 32), pattern)
	}
	return nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return strings.TrimSpace(s[:n]) + "..."
}
