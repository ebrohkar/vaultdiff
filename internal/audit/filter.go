package audit

import (
	"strings"
	"time"
)

// Filter holds criteria for selecting audit entries.
type Filter struct {
	Environment string
	PathPrefix  string
	Since       time.Time
	Until       time.Time
	MinChanges  int
}

// Match reports whether the given Entry satisfies all non-zero filter criteria.
func (f Filter) Match(e Entry) bool {
	if f.Environment != "" && e.Environment != f.Environment {
		return false
	}

	if f.PathPrefix != "" && !strings.HasPrefix(e.Path, f.PathPrefix) {
		return false
	}

	if !f.Since.IsZero() && e.Timestamp.Before(f.Since) {
		return false
	}

	if !f.Until.IsZero() && e.Timestamp.After(f.Until) {
		return false
	}

	if f.MinChanges > 0 && e.ChangeCount < f.MinChanges {
		return false
	}

	return true
}

// ApplyFilter returns only entries that match the given filter.
func ApplyFilter(entries []Entry, f Filter) []Entry {
	result := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if f.Match(e) {
			result = append(result, e)
		}
	}
	return result
}
