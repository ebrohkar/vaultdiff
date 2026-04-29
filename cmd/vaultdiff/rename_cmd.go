package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/example/vaultdiff/internal/config"
	"github.com/example/vaultdiff/internal/vault"
)

var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename a secret by copying it to a new path (and optionally deleting the source)",
	RunE:  runRename,
}

func init() {
	renameCmd.Flags().String("source", "", "Source secret path (required)")
	renameCmd.Flags().String("dest", "", "Destination secret path (required)")
	renameCmd.Flags().String("mount", "secret", "KV v2 mount point")
	renameCmd.Flags().Bool("delete-source", false, "Delete the source secret after copying")
	_ = renameCmd.MarkFlagRequired("source")
	_ = renameCmd.MarkFlagRequired("dest")
	rootCmd.AddCommand(renameCmd)
}

func runRename(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	source, _ := cmd.Flags().GetString("source")
	dest, _ := cmd.Flags().GetString("dest")
	mount, _ := cmd.Flags().GetString("mount")
	delSrc, _ := cmd.Flags().GetBool("delete-source")

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	result, err := vault.RenameSecret(cmd.Context(), client, vault.RenameOptions{
		SourcePath:   source,
		DestPath:     dest,
		Mount:        mount,
		DeleteSource: delSrc,
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Renamed %q → %q\n", result.SourcePath, result.DestPath)
	if result.DeletedSource {
		fmt.Fprintf(os.Stdout, "Source %q deleted.\n", result.SourcePath)
	}
	return nil
}
