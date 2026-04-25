package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultdiff/internal/audit"
)

var (
	auditLogFile   string
	auditEnv       string
	auditPathPrefix string
	auditSince     string
	auditMinChanges int
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Query and filter the audit log",
	Example: `  vaultdiff audit --log audit.json --env production
  vaultdiff audit --log audit.json --path-prefix secret/myapp --since 24h`,
	RunE: runAudit,
}

func init() {
	auditCmd.Flags().StringVar(&auditLogFile, "log", "audit.log", "Path to audit log file")
	auditCmd.Flags().StringVar(&auditEnv, "env", "", "Filter by environment")
	auditCmd.Flags().StringVar(&auditPathPrefix, "path-prefix", "", "Filter by secret path prefix")
	auditCmd.Flags().StringVar(&auditSince, "since", "", "Filter entries newer than duration (e.g. 24h, 7d)")
	auditCmd.Flags().IntVar(&auditMinChanges, "min-changes", 0, "Filter entries with at least N changes")
}

func runAudit(cmd *cobra.Command, args []string) error {
	f, err := os.Open(auditLogFile)
	if err != nil {
		return fmt.Errorf("opening audit log: %w", err)
	}
	defer f.Close()

	entries, err := audit.ReadEntries(f)
	if err != nil {
		return fmt.Errorf("reading audit entries: %w", err)
	}

	filter := audit.Filter{
		Environment: auditEnv,
		PathPrefix:  auditPathPrefix,
		MinChanges:  auditMinChanges,
	}

	if auditSince != "" {
		d, err := time.ParseDuration(auditSince)
		if err != nil {
			return fmt.Errorf("parsing --since duration: %w", err)
		}
		filter.Since = time.Now().UTC().Add(-d)
	}

	filtered := audit.ApplyFilter(entries, filter)
	for _, e := range filtered {
		fmt.Fprintf(os.Stdout, "[%s] %s @ %s (v%d→v%d) changes=%d\n",
			e.Timestamp.Format(time.RFC3339),
			e.Environment,
			e.Path,
			e.VersionA,
			e.VersionB,
			e.ChangeCount,
		)
	}
	return nil
}
