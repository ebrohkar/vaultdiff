package diff

import (
	"fmt"
	"strings"
)

// VersionedDiff represents the diff between two consecutive secret versions.
type VersionedDiff struct {
	FromVersion int
	ToVersion   int
	Changes     []Change
}

// DiffHistory computes the sequential diffs between each consecutive pair
// of secret version data maps.
func DiffHistory(versions []map[string]string) ([]VersionedDiff, error) {
	if len(versions) < 2 {
		return nil, fmt.Errorf("at least two versions are required to compute history diff")
	}

	var results []VersionedDiff

	for i := 1; i < len(versions); i++ {
		changes := Compare(versions[i-1], versions[i])
		results = append(results, VersionedDiff{
			FromVersion: i,
			ToVersion:   i + 1,
			Changes:     changes,
		})
	}

	return results, nil
}

// SummariseHistory returns a human-readable summary of all versioned diffs.
func SummariseHistory(diffs []VersionedDiff) string {
	if len(diffs) == 0 {
		return "No history diffs available."
	}

	var sb strings.Builder
	for _, d := range diffs {
		sb.WriteString(fmt.Sprintf("v%d -> v%d: %d change(s)\n", d.FromVersion, d.ToVersion, len(d.Changes)))
		for _, c := range d.Changes {
			sb.WriteString(fmt.Sprintf("  [%s] %s\n", c.Type, c.Key))
		}
	}
	return sb.String()
}
