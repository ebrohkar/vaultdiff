package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/example/vaultdiff/internal/diff"
	"github.com/example/vaultdiff/internal/vault"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show the full diff history of a secret across all versions",
	RunE:  runHistory,
}

func init() {
	historyCmd.Flags().String("mount", "secret", "KV mount path")
	historyCmd.Flags().String("path", "", "Secret path (required)")
	historyCmd.Flags().String("format", "text", "Output format: text, json, markdown")
	_ = historyCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(historyCmd)
}

func runHistory(cmd *cobra.Command, _ []string) error {
	mount, _ := cmd.Flags().GetString("mount")
	path, _ := cmd.Flags().GetString("path")
	format, _ := cmd.Flags().GetString("format")

	if format != "text" && format != "json" && format != "markdown" {
		return fmt.Errorf("invalid format %q: must be text, json, or markdown", format)
	}

	client, err := vault.NewClient("", "")
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	ctx := context.Background()
	history, err := vault.FetchHistory(ctx, client, mount, path)
	if err != nil {
		return fmt.Errorf("fetching history: %w", err)
	}

	var versionMaps []map[string]string
	for _, entry := range history.Entries {
		if entry.Data != nil {
			sm := make(map[string]string, len(entry.Data))
			for k, v := range entry.Data {
				sm[k] = fmt.Sprintf("%v", v)
			}
			versionMaps = append(versionMaps, sm)
		}
	}

	if len(versionMaps) < 2 {
		fmt.Fprintln(os.Stdout, "Not enough readable versions to compute history diff.")
		return nil
	}

	diffs, err := diff.DiffHistory(versionMaps)
	if err != nil {
		return fmt.Errorf("computing history diff: %w", err)
	}

	fmt.Fprint(os.Stdout, diff.SummariseHistory(diffs))
	return nil
}
