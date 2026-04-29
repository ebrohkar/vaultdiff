// Package vault provides functionality for interacting with HashiCorp Vault's
// KV v2 secrets engine.
//
// # Rename
//
// RenameSecret copies a secret from one path to another within the same KV v2
// mount. It reads the latest version of the source secret and writes it to the
// destination path. Optionally, the source secret can be deleted after the copy
// to complete the rename.
//
// Example usage:
//
//	result, err := vault.RenameSecret(ctx, client, vault.RenameOptions{
//		SourcePath:   "app/old-config",
//		DestPath:     "app/new-config",
//		Mount:        "secret",
//		DeleteSource: true,
//	})
//
// The RenameResult contains the source and destination paths, the version
// written to the destination, whether the source was deleted, and a UTC
// timestamp of when the operation completed.
//
// Note: Rename is not an atomic operation in Vault. If the write succeeds but
// the delete fails, both source and destination will exist. Callers should
// handle this scenario appropriately.
package vault
