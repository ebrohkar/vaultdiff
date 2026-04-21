package diff

import "fmt"

// ChangeType represents the kind of change detected for a key.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// Change represents a single key-level difference between two secret versions.
type Change struct {
	Key      string
	Type     ChangeType
	OldValue interface{}
	NewValue interface{}
}

// Result holds the full diff output between two secret snapshots.
type Result struct {
	Changes []Change
	HasDiff bool
}

// Compare computes the diff between two secret data maps.
// Values are compared by their string representation to handle mixed types.
func Compare(oldData, newData map[string]interface{}) *Result {
	result := &Result{}
	seen := map[string]bool{}

	for key, oldVal := range oldData {
		seen[key] = true
		newVal, exists := newData[key]
		if !exists {
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Type:     Removed,
				OldValue: oldVal,
				NewValue: nil,
			})
			result.HasDiff = true
			continue
		}
		if fmt.Sprintf("%v", oldVal) != fmt.Sprintf("%v", newVal) {
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Type:     Modified,
				OldValue: oldVal,
				NewValue: newVal,
			})
			result.HasDiff = true
		} else {
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Type:     Unchanged,
				OldValue: oldVal,
				NewValue: newVal,
			})
		}
	}

	for key, newVal := range newData {
		if !seen[key] {
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Type:     Added,
				OldValue: nil,
				NewValue: newVal,
			})
			result.HasDiff = true
		}
	}

	return result
}
