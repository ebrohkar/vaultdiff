package vault

import (
	"errors"
	"fmt"
	"time"
)

// LockOptions configures a secret lock operation.
type LockOptions struct {
	Path   string
	Mount  string
	Reason string
	TTL    time.Duration
}

// LockResult holds the outcome of a lock or unlock operation.
type LockResult struct {
	Path      string
	Mount     string
	Locked    bool
	Reason    string
	LockedAt  time.Time
	ExpiresAt time.Time
}

// LockSecret marks a secret path as locked by writing a lock metadata tag.
// Locking is advisory: it records intent and blocks vaultdiff promote/rollback
// commands from modifying the path without an explicit --force flag.
func LockSecret(c *Client, opts LockOptions) (*LockResult, error) {
	if opts.Path == "" {
		return nil, errors.New("lock: path must not be empty")
	}
	if opts.Mount == "" {
		opts.Mount = "secret"
	}
	if opts.TTL <= 0 {
		opts.TTL = 24 * time.Hour
	}

	now := time.Now().UTC()
	expiresAt := now.Add(opts.TTL)

	tags := map[string]string{
		"vaultdiff:locked":     "true",
		"vaultdiff:lock_reason": opts.Reason,
		"vaultdiff:locked_at":  now.Format(time.RFC3339),
		"vaultdiff:expires_at": expiresAt.Format(time.RFC3339),
	}

	_, err := SetTags(c, TagOptions{
		Path:  opts.Path,
		Mount: opts.Mount,
		Tags:  tags,
	})
	if err != nil {
		return nil, fmt.Errorf("lock: failed to set lock tags on %q: %w", opts.Path, err)
	}

	return &LockResult{
		Path:      opts.Path,
		Mount:     opts.Mount,
		Locked:    true,
		Reason:    opts.Reason,
		LockedAt:  now,
		ExpiresAt: expiresAt,
	}, nil
}

// UnlockSecret removes the advisory lock from a secret path.
func UnlockSecret(c *Client, path, mount string) (*LockResult, error) {
	if path == "" {
		return nil, errors.New("unlock: path must not be empty")
	}
	if mount == "" {
		mount = "secret"
	}

	tags := map[string]string{
		"vaultdiff:locked":     "false",
		"vaultdiff:lock_reason": "",
		"vaultdiff:locked_at":  "",
		"vaultdiff:expires_at": "",
	}

	_, err := SetTags(c, TagOptions{
		Path:  path,
		Mount: mount,
		Tags:  tags,
	})
	if err != nil {
		return nil, fmt.Errorf("unlock: failed to clear lock tags on %q: %w", path, err)
	}

	return &LockResult{
		Path:   path,
		Mount:  mount,
		Locked: false,
	}, nil
}
