package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "vaultdiff",
	Short: "Diff and audit HashiCorp Vault secret versions across environments",
	Long: `vaultdiff is a CLI tool for comparing Vault secret versions
and auditing changes across different environments.`,
}

func init() {
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(auditCmd)
}
