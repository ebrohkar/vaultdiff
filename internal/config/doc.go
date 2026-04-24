// Package config provides loading and validation of runtime configuration
// for the vaultdiff CLI tool.
//
// Configuration is sourced exclusively from environment variables so that
// the tool integrates cleanly with CI/CD pipelines and secret-management
// workflows without requiring config files on disk.
//
// Supported environment variables:
//
//	VAULT_ADDR            – Vault server address (default: http://127.0.0.1:8200)
//	VAULT_TOKEN           – Vault authentication token (required)
//	VAULTDIFF_FORMAT      – Output format: text | json | markdown (default: text)
//	VAULTDIFF_AUDIT_LOG   – File path for audit log output (optional)
//	VAULTDIFF_ENVIRONMENTS – Comma-separated list of environment names (optional)
//
// Example:
//
//	cfg, err := config.Load()
//	if err != nil {
//	    log.Fatalf("configuration error: %v", err)
//	}
package config
