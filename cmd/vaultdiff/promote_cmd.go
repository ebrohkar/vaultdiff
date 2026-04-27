package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/example/vaultdiff/internal/config"
	"github.com/example/vaultdiff/internal/vault"
)

var promoteCmd = &cobra.Command{
	Use:   "promote",
	Short: "Promote a secret from one environment to another",
	RunE:  runPromote,
}

func init() {
	promoteCmd.Flags().String("source-path", "", "Source secret path (required)")
	promoteCmd.Flags().String("dest-path", "", "Destination secret path (required)")
	promoteCmd.Flags().String("source-env", "staging", "Source environment label")
	promoteCmd.Flags().String("dest-env", "production", "Destination environment label")
	promoteCmd.Flags().Int("version", 0, "Source secret version (0 = latest)")
	promoteCmd.Flags().Bool("dry-run", false, "Preview promotion without writing")
	_ = promoteCmd.MarkFlagRequired("source-path")
	_ = promoteCmd.MarkFlagRequired("dest-path")
	rootCmd.AddCommand(promoteCmd)
}

func runPromote(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	srcPath, _ := cmd.Flags().GetString("source-path")
	dstPath, _ := cmd.Flags().GetString("dest-path")
	srcEnv, _ := cmd.Flags().GetString("source-env")
	dstEnv, _ := cmd.Flags().GetString("dest-env")
	version, _ := cmd.Flags().GetInt("version")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	opts := vault.PromoteOptions{
		SourcePath: srcPath,
		DestPath:   dstPath,
		SourceEnv:  srcEnv,
		DestEnv:    dstEnv,
		Version:    version,
		DryRun:     dryRun,
	}

	result, err := vault.PromoteSecret(cmd.Context(), client, client, opts)
	if err != nil {
		return err
	}

	if result.DryRun {
		fmt.Fprintf(os.Stdout, "[dry-run] would promote %d key(s) from %s (%s) to %s (%s)\n",
			result.KeysWritten, result.SourcePath, result.SourceEnv, result.DestPath, result.DestEnv)
	} else {
		fmt.Fprintf(os.Stdout, "promoted %d key(s) from %s (%s) to %s (%s)\n",
			result.KeysWritten, result.SourcePath, result.SourceEnv, result.DestPath, result.DestEnv)
	}
	return nil
}
