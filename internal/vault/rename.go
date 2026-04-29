package vault

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// RenameOptions configures a secret rename operation.
type RenameOptions struct {
	SourcePath string
	DestPath   string
	DeleteSource bool
	Mount      string
}

// RenameResult holds the outcome of a rename operation.
type RenameResult struct {
	SourcePath  string
	DestPath    string
	Version     int
	DeletedSource bool
	Timestamp   time.Time
}

// RenameSecret copies a secret from sourcePath to destPath and optionally
// deletes the source, effectively renaming it within Vault KV v2.
func RenameSecret(ctx context.Context, client *Client, opts RenameOptions) (*RenameResult, error) {
	if opts.SourcePath == "" {
		return nil, errors.New("rename: source path must not be empty")
	}
	if opts.DestPath == "" {
		return nil, errors.New("rename: destination path must not be empty")
	}
	if opts.SourcePath == opts.DestPath {
		return nil, errors.New("rename: source and destination paths must differ")
	}
	if opts.Mount == "" {
		opts.Mount = "secret"
	}

	// Read the latest version from source.
	secret, err := client.ReadSecretVersion(ctx, opts.SourcePath, 0)
	if err != nil {
		return nil, fmt.Errorf("rename: read source %q: %w", opts.SourcePath, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("rename: source secret %q not found", opts.SourcePath)
	}

	// Write to destination.
	writePath := opts.Mount + "/data/" + opts.DestPath
	_, err = client.vault.KVv2(opts.Mount).Put(ctx, opts.DestPath, secret)
	if err != nil {
		return nil, fmt.Errorf("rename: write dest %q (%s): %w", opts.DestPath, writePath, err)
	}

	result := &RenameResult{
		SourcePath: opts.SourcePath,
		DestPath:   opts.DestPath,
		Timestamp:  time.Now().UTC(),
	}

	if opts.DeleteSource {
		if err := client.vault.KVv2(opts.Mount).Delete(ctx, opts.SourcePath); err != nil {
			return nil, fmt.Errorf("rename: delete source %q: %w", opts.SourcePath, err)
		}
		result.DeletedSource = true
	}

	return result, nil
}
