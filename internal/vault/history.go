package vault

import (
	"context"
	"fmt"
	"sort"
	"time"
)

// SecretHistoryEntry represents a single version entry in a secret's history.
type SecretHistoryEntry struct {
	Version     int
	CreatedTime time.Time
	DeletedTime *time.Time
	Destroyed   bool
	Data        map[string]interface{}
}

// SecretHistory holds all version entries for a given secret path.
type SecretHistory struct {
	Path    string
	Entries []SecretHistoryEntry
}

// FetchHistory retrieves all available versions of a secret and returns them
// sorted from oldest to newest.
func FetchHistory(ctx context.Context, c *Client, mountPath, secretPath string) (*SecretHistory, error) {
	if secretPath == "" {
		return nil, fmt.Errorf("secret path must not be empty")
	}

	versions, err := ListSecretVersions(ctx, c, mountPath, secretPath)
	if err != nil {
		return nil, fmt.Errorf("listing versions for %q: %w", secretPath, err)
	}

	history := &SecretHistory{Path: secretPath}

	for _, v := range versions {
		entry := SecretHistoryEntry{
			Version:     v.Version,
			CreatedTime: v.CreatedTime,
			DeletedTime: v.DeletedTime,
			Destroyed:   v.Destroyed,
		}

		if !v.Destroyed && v.DeletedTime == nil {
			data, err := ReadSecretVersion(ctx, c, mountPath, secretPath, v.Version)
			if err == nil {
				entry.Data = data
			}
		}

		history.Entries = append(history.Entries, entry)
	}

	sort.Slice(history.Entries, func(i, j int) bool {
		return history.Entries[i].Version < history.Entries[j].Version
	})

	return history, nil
}
