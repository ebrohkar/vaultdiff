// Package audit provides structured logging and filtering for vaultdiff
// secret comparison operations.
//
// Each time a diff is performed between two Vault secret versions, an audit
// Entry is recorded via a Logger. Entries capture the environment, secret
// path, version range, timestamp, and the full list of changes detected.
//
// Entries can be emitted in JSON (machine-readable) or text (human-readable)
// format, making the audit trail suitable for both log aggregation pipelines
// and interactive terminal review.
//
// The Filter type allows callers to select a subset of entries by environment,
// path prefix, time range, or minimum change count — useful when replaying or
// reviewing a stored audit log.
//
// Basic usage:
//
//	logger := audit.NewLogger(os.Stdout, audit.FormatJSON)
//	err := logger.Record("production", "secret/app/config", 3, 4, changes)
//
// Filtering:
//
//	matched := audit.ApplyFilter(entries, audit.Filter{
//		Environment: "production",
//		MinChanges:  1,
//	})
package audit
