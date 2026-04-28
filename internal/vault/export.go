package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ExportFormat defines the output format for exported secrets.
type ExportFormat string

const (
	ExportFormatJSON ExportFormat = "json"
	ExportFormatEnv  ExportFormat = "env"
)

// ExportOptions controls the behaviour of ExportSecret.
type ExportOptions struct {
	Path       string
	Version    int
	Format     ExportFormat
	OutputFile string
	Redact     bool
}

// ExportResult holds the outcome of an export operation.
type ExportResult struct {
	Path        string
	Version     int
	Format      ExportFormat
	OutputFile  string
	ExportedAt  time.Time
	KeyCount    int
}

// ExportSecret reads a secret version from Vault and writes it to a file or stdout.
func ExportSecret(client *Client, opts ExportOptions) (*ExportResult, error) {
	if opts.Path == "" {
		return nil, fmt.Errorf("export: path must not be empty")
	}
	if opts.Format == "" {
		opts.Format = ExportFormatJSON
	}
	if opts.Format != ExportFormatJSON && opts.Format != ExportFormatEnv {
		return nil, fmt.Errorf("export: unsupported format %q", opts.Format)
	}

	data, err := client.ReadSecretVersion(opts.Path, opts.Version)
	if err != nil {
		return nil, fmt.Errorf("export: read secret: %w", err)
	}

	if opts.Redact {
		for k := range data {
			data[k] = "***"
		}
	}

	var payload []byte
	switch opts.Format {
	case ExportFormatJSON:
		payload, err = json.MarshalIndent(data, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("export: marshal json: %w", err)
		}
	case ExportFormatEnv:
		payload = marshalEnv(data)
	}

	var out *os.File
	if opts.OutputFile != "" {
		if err := os.MkdirAll(filepath.Dir(opts.OutputFile), 0o750); err != nil {
			return nil, fmt.Errorf("export: create dirs: %w", err)
		}
		out, err = os.OpenFile(opts.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
		if err != nil {
			return nil, fmt.Errorf("export: open file: %w", err)
		}
		defer out.Close()
	} else {
		out = os.Stdout
	}

	if _, err := out.Write(payload); err != nil {
		return nil, fmt.Errorf("export: write output: %w", err)
	}

	return &ExportResult{
		Path:       opts.Path,
		Version:    opts.Version,
		Format:     opts.Format,
		OutputFile: opts.OutputFile,
		ExportedAt: time.Now().UTC(),
		KeyCount:   len(data),
	}, nil
}

func marshalEnv(data map[string]string) []byte {
	var buf []byte
	for k, v := range data {
		line := fmt.Sprintf("%s=%s\n", k, v)
		buf = append(buf, []byte(line)...)
	}
	return buf
}
