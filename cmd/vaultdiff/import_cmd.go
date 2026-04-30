package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/config"
	"github.com/yourusername/vaultdiff/internal/vault"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import secrets from a file into Vault",
	RunE:  runImport,
}

func init() {
	importCmd.Flags().String("path", "", "Vault secret path (required)")
	importCmd.Flags().String("file", "", "Path to the input file (required)")
	importCmd.Flags().String("format", "json", "Input format: json or env")
	importCmd.Flags().String("mount", "secret", "Vault KV mount point")
	importCmd.Flags().Bool("dry-run", false, "Preview import without writing to Vault")
	_ = importCmd.MarkFlagRequired("path")
	_ = importCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	path, _ := cmd.Flags().GetString("path")
	file, _ := cmd.Flags().GetString("file")
	format, _ := cmd.Flags().GetString("format")
	mount, _ := cmd.Flags().GetString("mount")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	result, err := vault.ImportSecret(vault.ImportOptions{
		Client:   client,
		Path:     path,
		Mount:    mount,
		FilePath: file,
		Format:   vault.ImportFormat(format),
		DryRun:   dryRun,
	})
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Fprintf(os.Stdout, "[dry-run] Would import %d key(s) to %s/%s\n", result.KeyCount, result.Mount, result.Path)
		for k, v := range result.Data {
			fmt.Fprintf(os.Stdout, "  %s = %s\n", k, v)
		}
		return nil
	}

	fmt.Fprintf(os.Stdout, "Imported %d key(s) to %s/%s\n", result.KeyCount, result.Mount, result.Path)
	return nil
}
