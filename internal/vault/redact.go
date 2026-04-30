package vault

import (
	"regexp"
	"strings"
)

// RedactOptions controls how secret values are redacted.
type RedactOptions struct {
	// Patterns is a list of key name patterns (case-insensitive) whose values
	// should be redacted in output.
	Patterns []string
	// Placeholder is the string used in place of a redacted value.
	Placeholder string
}

// RedactResult holds the redacted data map and metadata.
type RedactResult struct {
	Data        map[string]string
	RedactedKeys []string
}

// DefaultRedactOptions returns sensible defaults for redaction.
func DefaultRedactOptions() RedactOptions {
	return RedactOptions{
		Patterns: []string{
			"password", "passwd", "secret", "token", "api_key",
			"apikey", "private_key", "credential", "auth",
		},
		Placeholder: "[REDACTED]",
	}
}

// RedactSecret applies redaction rules to a string map, replacing sensitive
// values with the configured placeholder.
func RedactSecret(data map[string]string, opts RedactOptions) RedactResult {
	if opts.Placeholder == "" {
		opts.Placeholder = "[REDACTED]"
	}

	result := RedactResult{
		Data: make(map[string]string, len(data)),
	}

	for k, v := range data {
		if matchesAnyPattern(k, opts.Patterns) {
			result.Data[k] = opts.Placeholder
			result.RedactedKeys = append(result.RedactedKeys, k)
		} else {
			result.Data[k] = v
		}
	}

	return result
}

// matchesAnyPattern returns true if the key matches any of the given patterns
// using case-insensitive substring or regex matching.
func matchesAnyPattern(key string, patterns []string) bool {
	lower := strings.ToLower(key)
	for _, p := range patterns {
		re, err := regexp.Compile("(?i)" + p)
		if err != nil {
			// fallback to substring match
			if strings.Contains(lower, strings.ToLower(p)) {
				return true
			}
			continue
		}
		if re.MatchString(key) {
			return true
		}
	}
	return false
}
