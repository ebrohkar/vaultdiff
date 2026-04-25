package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultdiff/internal/config"
	"github.com/yourorg/vaultdiff/internal/diff"
	"github.com/yourorg/vaultdiff/internal/vault"
)

var (
	diffPath    string
	diffVersionA int
	diffVersionB int
	diffFormat  string
	diffEnv     string
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compare two versions of a Vault secret",
	Example: `  vaultdiff diff --path secret/myapp --version-a 1 --version-b 2
  vaultdiff diff --path secret/myapp --version-a 1 --version-b 2 --format json`,
	RunE: runDiff,
}

func init() {
	diffCmd.Flags().StringVar(&diffPath, "path", "", "Vault secret path (required)")
	diffCmd.Flags().IntVar(&diffVersionA, "version-a", 0, "First version to compare (required)")
	diffCmd.Flags().IntVar(&diffVersionB, "version-b", 0, "Second version to compare (required)")
	diffCmd.Flags().StringVar(&diffFormat, "format", "text", "Output format: text, json, markdown")
	diffCmd.Flags().StringVar(&diffEnv, "env", "", "Environment label for audit logging")
	_ = diffCmd.MarkFlagRequired("path")
	_ = diffCmd.MarkFlagRequired("version-a")
	_ = diffCmd.MarkFlagRequired("version-b")
}

func runDiff(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	pair, err := vault.FetchVersionPair(cmd.Context(), client, diffPath, diffVersionA, diffVersionB)
	if err != nil {
		return fmt.Errorf("fetching secret versions: %w", err)
	}

	mapA, mapB := pair.ToStringMaps()
	changes := diff.Compare(mapA, mapB)

	output, err := diff.Render(changes, diffFormat)
	if err != nil {
		return fmt.Errorf("rendering diff: %w", err)
	}

	fmt.Fprint(os.Stdout, output)
	return nil
}
