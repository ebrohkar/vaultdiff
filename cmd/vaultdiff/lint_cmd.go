package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/youorg/vaultdiff/internal/vault"
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint a Vault secret for common key/value issues",
	RunE:  runLint,
}

func init() {
	lintCmd.Flags().String("path", "", "Secret path to lint (required)")
	lintCmd.Flags().Int("version", 0, "Secret version (0 = latest)")
	lintCmd.Flags().Bool("no-empty", true, "Check for empty values")
	lintCmd.Flags().Bool("no-spaces", true, "Check for spaces in keys")
	lintCmd.Flags().Bool("upper-case", true, "Enforce upper-case keys")
	_ = lintCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(lintCmd)
}

func runLint(cmd *cobra.Command, _ []string) error {
	path, _ := cmd.Flags().GetString("path")
	version, _ := cmd.Flags().GetInt("version")
	checkEmpty, _ := cmd.Flags().GetBool("no-empty")
	checkSpaces, _ := cmd.Flags().GetBool("no-spaces")
	checkCase, _ := cmd.Flags().GetBool("upper-case")

	client, err := vault.NewClient("")
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secret, err := vault.ReadSecretVersion(cmd.Context(), client, path, version)
	if err != nil {
		return fmt.Errorf("read secret: %w", err)
	}

	data := vault.ToStringMap(secret)

	opts := vault.LintOptions{
		CheckEmptyValues: checkEmpty,
		CheckKeySpaces:   checkSpaces,
		CheckKeyCasing:   checkCase,
	}

	violations := vault.LintSecret(path, data, opts)
	if len(violations) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "✔ No lint violations found.")
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✖ %d violation(s) found:\n", len(violations))
	for _, v := range violations {
		fmt.Fprintln(cmd.OutOrStdout(), " ", v.String())
	}

	os.Exit(1)
	return nil
}
