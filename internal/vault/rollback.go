package vault

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// RollbackOptions configures a secret rollback operation.
type RollbackOptions struct {
	Path        string
	TargetVersion int
	DryRun      bool
}

// RollbackResult holds the outcome of a rollback operation.
type RollbackResult struct {
	Path          string
	FromVersion   int
	ToVersion     int
	DryRun        bool
	RolledBackAt  time.Time
	Data          map[string]interface{}
}

// RollbackSecret restores a KV v2 secret to a previous version.
// If DryRun is true, it fetches the target version data without writing.
func RollbackSecret(ctx context.Context, client *Client, opts RollbackOptions) (*RollbackResult, error) {
	if opts.Path == "" {
		return nil, errors.New("rollback: path must not be empty")
	}
	if opts.TargetVersion < 1 {
		return nil, errors.New("rollback: target version must be >= 1")
	}

	// Read the target version to confirm it exists and is not destroyed.
	secret, err := client.ReadSecretVersion(ctx, opts.Path, opts.TargetVersion)
	if err != nil {
		return nil, fmt.Errorf("rollback: failed to read version %d of %s: %w", opts.TargetVersion, opts.Path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("rollback: version %d of %s not found or destroyed", opts.TargetVersion, opts.Path)
	}

	data, _ := secret["data"].(map[string]interface{})

	result := &RollbackResult{
		Path:         opts.Path,
		ToVersion:    opts.TargetVersion,
		DryRun:       opts.DryRun,
		RolledBackAt: time.Now().UTC(),
		Data:         data,
	}

	if opts.DryRun {
		return result, nil
	}

	// Write the target version's data back as a new version.
	mountPath, secretPath := splitMountAndPath(opts.Path)
	writePath := fmt.Sprintf("%s/data/%s", mountPath, secretPath)

	payload := map[string]interface{}{"data": data}
	written, err := client.Logical().WriteWithContext(ctx, writePath, payload)
	if err != nil {
		return nil, fmt.Errorf("rollback: failed to write rollback data to %s: %w", opts.Path, err)
	}

	if written != nil && written.Data != nil {
		if v, ok := written.Data["version"].(float64); ok {
			result.FromVersion = int(v)
		}
	}

	return result, nil
}

// splitMountAndPath splits a full secret path into mount and sub-path.
// e.g. "secret/myapp/config" -> "secret", "myapp/config"
func splitMountAndPath(full string) (string, string) {
	for i, c := range full {
		if c == '/' {
			return full[:i], full[i+1:]
		}
	}
	return full, ""
}
