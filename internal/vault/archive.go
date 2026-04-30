package vault

import (
	"errors"
	"fmt"
	"time"
)

// ArchiveOptions configures the archive operation.
type ArchiveOptions struct {
	Mount   string
	Version int
}

// ArchiveEntry represents a single archived secret version.
type ArchiveEntry struct {
	Path      string
	Version   int
	Data      map[string]string
	ArchivedAt time.Time
}

// ArchiveResult holds the result of an archive operation.
type ArchiveResult struct {
	Entry     ArchiveEntry
	AlreadyArchived bool
}

// DefaultArchiveOptions returns ArchiveOptions with sensible defaults.
func DefaultArchiveOptions() ArchiveOptions {
	return ArchiveOptions{
		Mount:   "secret",
		Version: 0,
	}
}

// ArchiveSecret reads a specific version of a secret and returns it as an
// ArchiveEntry for offline storage or auditing purposes. It does not modify
// the secret in Vault.
func ArchiveSecret(client *Client, path string, opts ArchiveOptions) (*ArchiveResult, error) {
	if path == "" {
		return nil, errors.New("archive: path must not be empty")
	}
	if client == nil {
		return nil, errors.New("archive: client must not be nil")
	}
	if opts.Mount == "" {
		opts.Mount = "secret"
	}

	version := opts.Version
	secret, err := client.ReadSecretVersion(path, version)
	if err != nil {
		return nil, fmt.Errorf("archive: reading secret %q version %d: %w", path, version, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("archive: secret %q version %d not found", path, version)
	}

	data := toStringMap(secret)

	entry := ArchiveEntry{
		Path:       path,
		Version:    version,
		Data:       data,
		ArchivedAt: time.Now().UTC(),
	}

	return &ArchiveResult{
		Entry:           entry,
		AlreadyArchived: false,
	}, nil
}
