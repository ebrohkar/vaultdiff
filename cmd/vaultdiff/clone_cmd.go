package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/example/vaultdiff/internal/config"
	"github.com/example/vaultdiff/internal/vault"
)

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone a secret from one path to another",
	RunE:  runClone,
}

func init() {
	cloneCmd.Flags().String("src", "", "Source secret path (required)")
	cloneCmd.Flags().String("dst", "", "Destination secret path (required)")
	cloneCmd.Flags().Int("version", 0, "Source version to clone (0 = latest)")
	cloneCmd.Flags().Bool("overwrite", false, "Overwrite destination if it already exists")
	_ = cloneCmd.MarkFlagRequired("src")
	_ = cloneCmd.MarkFlagRequired("dst")
	rootCmd.AddCommand(cloneCmd)
}

func runClone(cmd *cobra.Command, _ []string) error {
	src, _ := cmd.Flags().GetString("src")
	dst, _ := cmd.Flags().GetString("dst")
	version, _ := cmd.Flags().GetInt("version")
	overwrite, _ := cmd.Flags().GetBool("overwrite")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	result, err := vault.CloneSecret(context.Background(), client, vault.CloneOptions{
		SourcePath: src,
		DestPath:   dst,
		Version:    version,
		Overwrite:  overwrite,
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Cloned %q → %q (%d keys) at %s\n",
		result.SourcePath,
		result.DestPath,
		len(result.Data),
		result.ClonedAt.Format("2006-01-02T15:04:05Z"),
	)
	return nil
}
