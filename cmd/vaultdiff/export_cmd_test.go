package main

import (
	"testing"

	"github.com/spf13/cobra"
)

func executeExportCmd(args []string) error {
	root := &cobra.Command{Use: "vaultdiff"}
	cmd := &cobra.Command{
		Use:  "export",
		RunE: runExport,
	}
	cmd.Flags().String("path", "", "Vault secret path")
	cmd.Flags().Int("version", 0, "Secret version")
	cmd.Flags().String("format", "json", "Output format")
	cmd.Flags().String("output", "", "Output file")
	cmd.Flags().Bool("redact", false, "Redact values")
	root.AddCommand(cmd)
	root.SetArgs(args)
	return root.Execute()
}

func TestExportCmd_MissingPath(t *testing.T) {
	err := executeExportCmd([]string{"export"})
	if err == nil {
		t.Fatal("expected error when --path is missing")
	}
}

func TestExportCmd_InvalidFormat(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "test-token")
	t.Setenv("VAULT_ADDR", "http://127.0.0.1:8200")
	err := executeExportCmd([]string{"export", "--path", "secret/data/app", "--format", "yaml"})
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestExportCmd_DefaultFlags(t *testing.T) {
	root := &cobra.Command{Use: "vaultdiff"}
	cmd := &cobra.Command{
		Use:  "export",
		RunE: runExport,
	}
	cmd.Flags().String("path", "", "Vault secret path")
	cmd.Flags().Int("version", 0, "Secret version")
	cmd.Flags().String("format", "json", "Output format")
	cmd.Flags().String("output", "", "Output file")
	cmd.Flags().Bool("redact", false, "Redact values")
	root.AddCommand(cmd)

	format, err := cmd.Flags().GetString("format")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if format != "json" {
		t.Errorf("expected default format json, got %s", format)
	}

	version, err := cmd.Flags().GetInt("version")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version != 0 {
		t.Errorf("expected default version 0, got %d", version)
	}
}
