package vault

import (
	"context"
	"fmt"
	"time"
)

// WatchEvent represents a detected change in a secret between polls.
type WatchEvent struct {
	Path      string
	OldVersion int
	NewVersion int
	ChangedAt  time.Time
}

// WatchOptions configures the polling behaviour of WatchSecret.
type WatchOptions struct {
	// Interval between polls. Defaults to 30s if zero.
	Interval time.Duration
	// MountPath is the KV v2 mount (e.g. "secret").
	MountPath string
}

// WatchSecret polls a Vault KV v2 secret at the given path and emits a
// WatchEvent on the returned channel whenever the current version changes.
// The caller must cancel ctx to stop watching; the channel is closed on exit.
func WatchSecret(ctx context.Context, c *Client, path string, opts WatchOptions) (<-chan WatchEvent, error) {
	if path == "" {
		return nil, fmt.Errorf("watch: path must not be empty")
	}
	if opts.Interval <= 0 {
		opts.Interval = 30 * time.Second
	}
	if opts.MountPath == "" {
		opts.MountPath = "secret"
	}

	events := make(chan WatchEvent, 8)

	go func() {
		defer close(events)

		var lastVersion int

		ticker := time.NewTicker(opts.Interval)
		defer ticker.Stop()

		// poll immediately, then on each tick
		for {
			versions, err := ListSecretVersions(ctx, c, opts.MountPath, path)
			if err == nil && len(versions) > 0 {
				latest := versions[len(versions)-1]
				if lastVersion != 0 && latest.Version != lastVersion {
					events <- WatchEvent{
						Path:       path,
						OldVersion: lastVersion,
						NewVersion: latest.Version,
						ChangedAt:  time.Now().UTC(),
					}
				}
				lastVersion = latest.Version
			}

			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}()

	return events, nil
}
