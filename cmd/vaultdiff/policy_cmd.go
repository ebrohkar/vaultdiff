package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultdiff/internal/config"
	"github.com/yourorg/vaultdiff/internal/vault"
)

var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Check Vault token access for a secret path",
	RunE:  runPolicy,
}

func init() {
	policyCmd.Flags().StringP("path", "p", "", "Secret path to check (required)")
	_ = policyCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(policyCmd)
}

func runPolicy(cmd *cobra.Command, _ []string) error {
	path, _ := cmd.Flags().GetString("path")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := vault.NewClient(nil, cfg.Token)
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	access, err := client.CheckPathAccess(cmd.Context(), path)
	if err != nil {
		return fmt.Errorf("policy check failed: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Path:     %s\n", access.Path)
	fmt.Fprintf(os.Stdout, "CanRead:  %v\n", access.CanRead)
	fmt.Fprintf(os.Stdout, "CanList:  %v\n", access.CanList)
	fmt.Fprintf(os.Stdout, "Denied:   %v\n", access.Denied)

	if access.Denied {
		fmt.Fprintln(os.Stderr, "warning: token does not have access to this path")
		os.Exit(2)
	}

	return nil
}
