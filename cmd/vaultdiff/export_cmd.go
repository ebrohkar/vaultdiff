package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/example/vaultdiff/internal/config"
	"github.com/example/vaultdiff/internal/vault"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export a Vault secret version to a file or stdout",
	RunE:  runExport,
}

func init() {
	exportCmd.Flags().String("path", "", "Vault secret path (required)")
	exportCmd.Flags().Int("version", 0, "Secret version to export (0 = latest)")
	exportCmd.Flags().String("format", "json", "Output format: json or env")
	exportCmd.Flags().String("output", "", "Output file path (default: stdout)")
	exportCmd.Flags().Bool("redact", false, "Replace secret values with *** in output")
	_ = exportCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	path, _ := cmd.Flags().GetString("path")
	version, _ := cmd.Flags().GetInt("version")
	format, _ := cmd.Flags().GetString("format")
	output, _ := cmd.Flags().GetString("output")
	redact, _ := cmd.Flags().GetBool("redact")

	result, err := vault.ExportSecret(client, vault.ExportOptions{
		Path:       path,
		Version:    version,
		Format:     vault.ExportFormat(format),
		OutputFile: output,
		Redact:     redact,
	})
	if err != nil {
		return fmt.Errorf("export: %w", err)
	}

	if output != "" {
		fmt.Fprintf(os.Stderr, "exported %d key(s) from %s (v%d) to %s\n",
			result.KeyCount, result.Path, result.Version, result.OutputFile)
	}
	return nil
}
