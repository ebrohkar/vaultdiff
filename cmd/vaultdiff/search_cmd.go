package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/config"
	"github.com/yourusername/vaultdiff/internal/vault"
)

var (
	searchMount         string
	searchPathPrefix    string
	searchCaseSensitive bool
)

func init() {
	searchCmd := &cobra.Command{
		Use:   "search <keyword>",
		Short: "Search for a keyword across secret keys and values",
		Args:  cobra.ExactArgs(1),
		RunE:  runSearch,
	}

	searchCmd.Flags().StringVar(&searchMount, "mount", "secret", "KV mount name")
	searchCmd.Flags().StringVar(&searchPathPrefix, "prefix", "", "Path prefix to restrict search")
	searchCmd.Flags().BoolVar(&searchCaseSensitive, "case-sensitive", false, "Enable case-sensitive matching")

	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	keyword := args[0]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	opts := vault.SearchOptions{
		Mount:         searchMount,
		PathPrefix:    searchPathPrefix,
		Keyword:       keyword,
		CaseSensitive: searchCaseSensitive,
	}

	results, err := vault.SearchSecrets(client, opts)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if len(results) == 0 {
		fmt.Fprintln(os.Stdout, "No matches found.")
		return nil
	}

	for _, r := range results {
		fmt.Fprintf(os.Stdout, "Path: %s\n", r.Path)
		fmt.Fprintf(os.Stdout, "  Matched keys: %s\n", strings.Join(r.MatchedKeys, ", "))
	}
	return nil
}
