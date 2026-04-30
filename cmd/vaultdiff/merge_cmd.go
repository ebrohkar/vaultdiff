package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/example/vaultdiff/internal/vault"
)

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge secrets from a source path into a destination path",
	RunE:  runMerge,
}

func init() {
	mergeCmd.Flags().String("source", "", "Source secret path (required)")
	mergeCmd.Flags().String("dest", "", "Destination secret path (required)")
	mergeCmd.Flags().String("mount", "secret", "KV mount name")
	mergeCmd.Flags().String("strategy", "source", "Conflict strategy: source, dest, or error")
	_ = mergeCmd.MarkFlagRequired("source")
	_ = mergeCmd.MarkFlagRequired("dest")
	rootCmd.AddCommand(mergeCmd)
}

func runMerge(cmd *cobra.Command, _ []string) error {
	source, _ := cmd.Flags().GetString("source")
	dest, _ := cmd.Flags().GetString("dest")
	mount, _ := cmd.Flags().GetString("mount")
	strategy, _ := cmd.Flags().GetString("strategy")

	client, err := vault.NewClient(os.Getenv("VAULT_ADDR"), os.Getenv("VAULT_TOKEN"))
	if err != nil {
		return fmt.Errorf("merge: failed to create vault client: %w", err)
	}

	opts := vault.MergeOptions{
		Mount:      mount,
		SourcePath: source,
		DestPath:   dest,
		Strategy:   strategy,
	}

	result, err := vault.MergeSecret(client, opts)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Merged %s → %s\n", result.SourcePath, result.DestPath)
	fmt.Fprintf(cmd.OutOrStdout(), "  Added keys:    %s\n", joinOrNone(result.AddedKeys))
	fmt.Fprintf(cmd.OutOrStdout(), "  Conflicts:     %s\n", joinOrNone(result.Conflicts))
	fmt.Fprintf(cmd.OutOrStdout(), "  Skipped keys:  %s\n", joinOrNone(result.SkippedKeys))
	fmt.Fprintf(cmd.OutOrStdout(), "  Total keys:    %d\n", len(result.Merged))
	return nil
}

func joinOrNone(ss []string) string {
	if len(ss) == 0 {
		return "(none)"
	}
	return strings.Join(ss, ", ")
}
