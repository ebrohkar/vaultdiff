package vault

import (
	"context"
	"fmt"
	"time"
)

// SecretSnapshot captures the state of a secret at a point in time.
type SecretSnapshot struct {
	Path        string            `json:"path"`
	Version     int               `json:"version"`
	Data        map[string]string `json:"data"`
	CapturedAt  time.Time         `json:"captured_at"`
	Environment string            `json:"environment"`
}

// TakeSnapshot reads the current (or specified) version of a secret and
// returns a SecretSnapshot for later comparison or audit logging.
func TakeSnapshot(ctx context.Context, c *Client, path, environment string, version int) (*SecretSnapshot, error) {
	if path == "" {
		return nil, fmt.Errorf("snapshot: path must not be empty")
	}

	raw, err := c.ReadSecretVersion(ctx, path, version)
	if err != nil {
		return nil, fmt.Errorf("snapshot: reading secret %q version %d: %w", path, version, err)
	}

	data := toStringMap(raw)

	return &SecretSnapshot{
		Path:        path,
		Version:     version,
		Data:        data,
		CapturedAt:  time.Now().UTC(),
		Environment: environment,
	}, nil
}

// DiffSnapshots returns the data maps from two snapshots ready for diff.Compare.
func DiffSnapshots(a, b *SecretSnapshot) (map[string]string, map[string]string, error) {
	if a == nil || b == nil {
		return nil, nil, fmt.Errorf("snapshot: both snapshots must be non-nil")
	}
	return a.Data, b.Data, nil
}
