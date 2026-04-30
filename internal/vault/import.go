package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
)

// ImportFormat represents the supported input formats for secret import.
type ImportFormat string

const (
	ImportFormatJSON ImportFormat = "json"
	ImportFormatEnv  ImportFormat = "env"
)

// ImportOptions configures how a secret is imported into Vault.
type ImportOptions struct {
	Client    *api.Client
	Path      string
	Mount     string
	FilePath  string
	Format    ImportFormat
	DryRun    bool
}

// ImportResult holds the outcome of an import operation.
type ImportResult struct {
	Path      string
	Mount     string
	KeyCount  int
	DryRun    bool
	Data      map[string]string
}

// DefaultImportOptions returns ImportOptions with sensible defaults.
func DefaultImportOptions() ImportOptions {
	return ImportOptions{
		Mount:  "secret",
		Format: ImportFormatJSON,
	}
}

// ImportSecret reads a file and writes its key-value pairs as a Vault secret.
func ImportSecret(opts ImportOptions) (*ImportResult, error) {
	if opts.Path == "" {
		return nil, fmt.Errorf("import: path must not be empty")
	}
	if opts.FilePath == "" {
		return nil, fmt.Errorf("import: file path must not be empty")
	}
	if opts.Client == nil {
		return nil, fmt.Errorf("import: vault client must not be nil")
	}

	raw, err := os.ReadFile(opts.FilePath)
	if err != nil {
		return nil, fmt.Errorf("import: reading file %q: %w", opts.FilePath, err)
	}

	var data map[string]string
	switch opts.Format {
	case ImportFormatJSON:
		if err := json.Unmarshal(raw, &data); err != nil {
			return nil, fmt.Errorf("import: parsing JSON: %w", err)
		}
	case ImportFormatEnv:
		data, err = parseEnvFile(string(raw))
		if err != nil {
			return nil, fmt.Errorf("import: parsing env file: %w", err)
		}
	default:
		return nil, fmt.Errorf("import: unsupported format %q", opts.Format)
	}

	result := &ImportResult{
		Path:     opts.Path,
		Mount:    opts.Mount,
		KeyCount: len(data),
		DryRun:   opts.DryRun,
		Data:     data,
	}

	if opts.DryRun {
		return result, nil
	}

	vaultPath := fmt.Sprintf("%s/data/%s", opts.Mount, opts.Path)
	payload := map[string]interface{}{"data": toAnyMap(data)}
	if _, err := opts.Client.Logical().Write(vaultPath, payload); err != nil {
		return nil, fmt.Errorf("import: writing to vault path %q: %w", vaultPath, err)
	}

	return result, nil
}

func parseEnvFile(content string) (map[string]string, error) {
	result := make(map[string]string)
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid env line: %q", line)
		}
		result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return result, nil
}

func toAnyMap(m map[string]string) map[string]interface{} {
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
