package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/example/vaultdiff/internal/vault"
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive a secret version to a local JSON file",
	RunE:  runArchive,
}

func init() {
	archiveCmd.Flags().String("path", "", "Vault secret path (required)")
	archiveCmd.Flags().Int("version", 0, "Secret version to archive (0 = latest)")
	archiveCmd.Flags().String("mount", "secret", "KV mount path")
	archiveCmd.Flags().String("output", "", "Output file path (default: stdout)")
	_ = archiveCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(archiveCmd)
}

func runArchive(cmd *cobra.Command, _ []string) error {
	path, _ := cmd.Flags().GetString("path")
	version, _ := cmd.Flags().GetInt("version")
	mount, _ := cmd.Flags().GetString("mount")
	output, _ := cmd.Flags().GetString("output")

	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("archive: loading config: %w", err)
	}

	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("archive: creating vault client: %w", err)
	}

	opts := vault.ArchiveOptions{
		Mount:   mount,
		Version: version,
	}

	result, err := vault.ArchiveSecret(client, path, opts)
	if err != nil {
		return fmt.Errorf("archive: %w", err)
	}

	encoded, err := json.MarshalIndent(result.Entry, "", "  ")
	if err != nil {
		return fmt.Errorf("archive: marshalling result: %w", err)
	}

	if output == "" {
		fmt.Println(string(encoded))
		return nil
	}

	if err := os.WriteFile(output, encoded, 0600); err != nil {
		return fmt.Errorf("archive: writing output file %q: %w", output, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "archived %q version %d to %s\n", path, version, output)
	return nil
}
