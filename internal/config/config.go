// Package config handles loading and validation of vaultdiff configuration
// from environment variables and optional config files.
package config

import (
	"errors"
	"os"
	"strings"
)

// Config holds all runtime configuration for vaultdiff.
type Config struct {
	// VaultAddr is the address of the Vault server.
	VaultAddr string
	// VaultToken is the authentication token.
	VaultToken string
	// OutputFormat controls diff rendering: text, json, or markdown.
	OutputFormat string
	// AuditLogPath is the optional path to write audit log entries.
	AuditLogPath string
	// Environments lists the named environments to compare.
	Environments []string
}

// Load reads configuration from environment variables.
// VAULT_ADDR, VAULT_TOKEN, VAULTDIFF_FORMAT, VAULTDIFF_AUDIT_LOG,
// and VAULTDIFF_ENVIRONMENTS are supported.
func Load() (*Config, error) {
	cfg := &Config{
		VaultAddr:    getEnvOrDefault("VAULT_ADDR", "http://127.0.0.1:8200"),
		VaultToken:   os.Getenv("VAULT_TOKEN"),
		OutputFormat: getEnvOrDefault("VAULTDIFF_FORMAT", "text"),
		AuditLogPath: os.Getenv("VAULTDIFF_AUDIT_LOG"),
	}

	if raw := os.Getenv("VAULTDIFF_ENVIRONMENTS"); raw != "" {
		for _, e := range strings.Split(raw, ",") {
			e = strings.TrimSpace(e)
			if e != "" {
				cfg.Environments = append(cfg.Environments, e)
			}
		}
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Validate checks that required fields are present and values are acceptable.
func (c *Config) Validate() error {
	if c.VaultToken == "" {
		return errors.New("config: VAULT_TOKEN must be set")
	}
	if c.VaultAddr == "" {
		return errors.New("config: VAULT_ADDR must not be empty")
	}
	valid := map[string]bool{"text": true, "json": true, "markdown": true}
	if !valid[c.OutputFormat] {
		return errors.New("config: VAULTDIFF_FORMAT must be one of: text, json, markdown")
	}
	return nil
}

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
