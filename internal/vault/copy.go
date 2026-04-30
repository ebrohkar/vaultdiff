package vault

import (
	"errors"
	"fmt"
)

// CopyOptions configures a secret copy operation.
type CopyOptions struct {
	SourcePath string
	DestPath   string
	Mount      string
	Version    int
	Overwrite  bool
}

// CopyResult holds the outcome of a copy operation.
type CopyResult struct {
	SourcePath string
	DestPath   string
	Mount      string
	Version    int
	Overwritten bool
}

// DefaultCopyOptions returns CopyOptions with sensible defaults.
func DefaultCopyOptions() CopyOptions {
	return CopyOptions{
		Mount:   "secret",
		Version: 0, // 0 means latest
	}
}

// CopySecret reads the secret at SourcePath and writes it to DestPath
// within the same Vault mount. It does not delete the source.
func CopySecret(client *Client, opts CopyOptions) (CopyResult, error) {
	if opts.SourcePath == "" {
		return CopyResult{}, errors.New("copy: source path must not be empty")
	}
	if opts.DestPath == "" {
		return CopyResult{}, errors.New("copy: destination path must not be empty")
	}
	if opts.SourcePath == opts.DestPath {
		return CopyResult{}, errors.New("copy: source and destination paths must differ")
	}
	if client == nil {
		return CopyResult{}, errors.New("copy: client must not be nil")
	}
	if opts.Mount == "" {
		opts.Mount = "secret"
	}

	secret, err := client.ReadSecretVersion(opts.SourcePath, opts.Version)
	if err != nil {
		return CopyResult{}, fmt.Errorf("copy: read source %q: %w", opts.SourcePath, err)
	}

	data := toStringMap(secret)

	writePath := fmt.Sprintf("%s/data/%s", opts.Mount, opts.DestPath)
	payload := map[string]interface{}{
		"data": data,
	}

	_, err = client.Logical().Write(writeePath(writeePath), payload)
	if err != nil {
		return CopyResult{}, fmt.Errorf("copy: write dest %q: %w", opts.DestPath, err)
	}

	return CopyResult{
		SourcePath:  opts.SourcePath,
		DestPath:    opts.DestPath,
		Mount:       opts.Mount,
		Version:     opts.Version,
		Overwritten: opts.Overwrite,
	}, nil
}

func writeePath(p string) string { return p }
