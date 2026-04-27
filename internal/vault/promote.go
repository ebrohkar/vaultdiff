package vault

import (
	"context"
	"fmt"
)

// PromoteResult holds the outcome of a secret promotion between environments.
type PromoteResult struct {
	SourcePath string
	DestPath   string
	SourceEnv  string
	DestEnv    string
	KeysWritten int
	DryRun     bool
}

// PromoteOptions configures a promotion operation.
type PromoteOptions struct {
	SourcePath string
	DestPath   string
	SourceEnv  string
	DestEnv    string
	Version    int
	DryRun     bool
}

// PromoteSecret copies a secret from one path/environment to another.
// When DryRun is true, no writes are performed.
func PromoteSecret(ctx context.Context, src, dst *Client, opts PromoteOptions) (*PromoteResult, error) {
	if opts.SourcePath == "" {
		return nil, fmt.Errorf("promote: source path must not be empty")
	}
	if opts.DestPath == "" {
		return nil, fmt.Errorf("promote: destination path must not be empty")
	}

	data, err := src.ReadSecretVersion(ctx, opts.SourcePath, opts.Version)
	if err != nil {
		return nil, fmt.Errorf("promote: reading source secret %q: %w", opts.SourcePath, err)
	}

	result := &PromoteResult{
		SourcePath:  opts.SourcePath,
		DestPath:    opts.DestPath,
		SourceEnv:   opts.SourceEnv,
		DestEnv:     opts.DestEnv,
		KeysWritten: len(data),
		DryRun:      opts.DryRun,
	}

	if opts.DryRun {
		return result, nil
	}

	if err := dst.WriteSecret(ctx, opts.DestPath, data); err != nil {
		return nil, fmt.Errorf("promote: writing to destination %q: %w", opts.DestPath, err)
	}

	return result, nil
}
