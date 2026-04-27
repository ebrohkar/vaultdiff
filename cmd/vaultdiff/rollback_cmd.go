package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/example/vaultdiff/internal/config"
	"github.com/example/vaultdiff/internal/vault"
)

var (
	rollbackPath    string
	rollbackVersion int
	rollbackDryRun  bool
)

func init() {
	rollbackCmd := &cobra.Command{
		Use:   "rollback",
		Short: "Roll back a secret to a previous version",
		RunE:  runRollback,
	}

	rollbackCmd.Flags().StringVar(&rollbackPath, "path", "", "Secret path to roll back (required)")
	rollbackCmd.Flags().IntVar(&rollbackVersion, "version", 0, "Target version to restore (required)")
	rollbackCmd.Flags().BoolVar(&rollbackDryRun, "dry-run", false, "Preview rollback without writing")

	_ = rollbackCmd.MarkFlagRequired("path")
	_ = rollbackCmd.MarkFlagRequired("version")

	rootCmd.AddCommand(rollbackCmd)
}

func runRollback(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	opts := vault.RollbackOptions{
		Path:          rollbackPath,
		TargetVersion: rollbackVersion,
		DryRun:        rollbackDryRun,
	}

	result, err := vault.RollbackSecret(cmd.Context(), client, opts)
	if err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	if result.DryRun {
		fmt.Fprintf(os.Stdout, "[dry-run] Would roll back %s to version %d\n", result.Path, result.ToVersion)
		fmt.Fprintf(os.Stdout, "[dry-run] Data keys: ")
		for k := range result.Data {
			fmt.Fprintf(os.Stdout, "%s ", k)
		}
		fmt.Fprintln(os.Stdout)
		return nil
	}

	fmt.Fprintf(os.Stdout, "Rolled back %s: version %d -> %d (at %s)\n",
		result.Path,
		result.ToVersion,
		result.FromVersion,
		result.RolledBackAt.Format("2006-01-02T15:04:05Z"),
	)
	return nil
}
