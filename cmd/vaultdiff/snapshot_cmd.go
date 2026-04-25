package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/example/vaultdiff/internal/config"
	"github.com/example/vaultdiff/internal/vault"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Capture a snapshot of a secret version and write it to a JSON file",
	RunE:  runSnapshot,
}

func init() {
	snapshotCmd.Flags().String("path", "", "Vault secret path (required)")
	snapshotCmd.Flags().String("env", "default", "Environment label for the snapshot")
	snapshotCmd.Flags().Int("version", 0, "Secret version to snapshot (0 = latest)")
	snapshotCmd.Flags().String("out", "snapshot.json", "Output file path")
	_ = snapshotCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(snapshotCmd)
}

func runSnapshot(cmd *cobra.Command, _ []string) error {
	path, _ := cmd.Flags().GetString("path")
	env, _ := cmd.Flags().GetString("env")
	version, _ := cmd.Flags().GetInt("version")
	out, _ := cmd.Flags().GetString("out")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	snap, err := vault.TakeSnapshot(context.Background(), client, path, env, version)
	if err != nil {
		return fmt.Errorf("taking snapshot: %w", err)
	}

	f, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("creating output file %q: %w", out, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return fmt.Errorf("encoding snapshot: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Snapshot written to %s (path=%s, version=%d, env=%s)\n",
		out, snap.Path, snap.Version, snap.Environment)
	return nil
}
