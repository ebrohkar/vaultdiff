package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/example/vaultdiff/internal/vault"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a secret against a set of rules",
	RunE:  runValidate,
}

func init() {
	validateCmd.Flags().String("path", "", "Secret path to validate (required)")
	validateCmd.Flags().String("mount", "secret", "KV mount point")
	validateCmd.Flags().StringSlice("require", []string{}, "Keys that must be present (comma-separated)")
	validateCmd.Flags().StringSlice("pattern", []string{}, "key=regex rules e.g. password=^.{8,}$")
	_ = validateCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, _ []string) error {
	path, _ := cmd.Flags().GetString("path")
	mount, _ := cmd.Flags().GetString("mount")
	requireKeys, _ := cmd.Flags().GetStringSlice("require")
	patternFlags, _ := cmd.Flags().GetStringSlice("pattern")

	opts := vault.DefaultValidateOptions()
	opts.Mount = mount

	for _, k := range requireKeys {
		if k != "" {
			opts.Rules = append(opts.Rules, vault.ValidationRule{Key: k, Required: true})
		}
	}

	for _, p := range patternFlags {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid pattern flag %q: expected key=regex", p)
		}
		opts.Rules = append(opts.Rules, vault.ValidationRule{Key: parts[0], Pattern: parts[1]})
	}

	client, err := vault.NewClient("")
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	result, err := vault.ValidateSecret(client, path, opts)
	if err != nil {
		return fmt.Errorf("validating secret: %w", err)
	}

	if result.Passed {
		fmt.Fprintf(cmd.OutOrStdout(), "✓ %s passed all validation rules\n", result.Path)
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✗ %s failed validation:\n", result.Path)
	for _, v := range result.Violations {
		fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", v)
	}
	os.Exit(1)
	return nil
}
