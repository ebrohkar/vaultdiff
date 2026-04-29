package vault

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// CloneOptions configures a secret clone operation.
type CloneOptions struct {
	SourcePath string
	DestPath   string
	Version    int
	Overwrite  bool
}

// CloneResult holds the outcome of a clone operation.
type CloneResult struct {
	SourcePath  string
	DestPath    string
	Version     int
	Overwritten bool
	ClonedAt    time.Time
	Data        map[string]string
}

// CloneSecret copies a secret from one path to another within the same Vault
// instance. If opts.Version is 0, the latest version is used. If the
// destination already exists and opts.Overwrite is false, an error is returned.
func CloneSecret(ctx context.Context, client *Client, opts CloneOptions) (*CloneResult, error) {
	if opts.SourcePath == "" {
		return nil, errors.New("clone: source path must not be empty")
	}
	if opts.DestPath == "" {
		return nil, errors.New("clone: destination path must not be empty")
	}
	if opts.SourcePath == opts.DestPath {
		return nil, errors.New("clone: source and destination paths must differ")
	}

	version := opts.Version
	if version < 0 {
		return nil, fmt.Errorf("clone: invalid version %d", version)
	}

	secret, err := client.ReadSecretVersion(ctx, opts.SourcePath, version)
	if err != nil {
		return nil, fmt.Errorf("clone: read source %q: %w", opts.SourcePath, err)
	}

	data := toStringMap(secret)

	if !opts.Overwrite {
		existing, _ := client.ReadSecretVersion(ctx, opts.DestPath, 0)
		if existing != nil {
			return nil, fmt.Errorf("clone: destination %q already exists; use overwrite to replace", opts.DestPath)
		}
	}

	if err := client.WriteSecret(ctx, opts.DestPath, data); err != nil {
		return nil, fmt.Errorf("clone: write destination %q: %w", opts.DestPath, err)
	}

	return &CloneResult{
		SourcePath:  opts.SourcePath,
		DestPath:    opts.DestPath,
		Version:     version,
		Overwritten: opts.Overwrite,
		ClonedAt:    time.Now().UTC(),
		Data:        data,
	}, nil
}
