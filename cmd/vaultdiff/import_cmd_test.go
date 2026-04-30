package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func executeImportCmd(t *testing.T, args []string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	cmd := &cobra.Command{Use: "vaultdiff"}
	cmd.AddCommand(importCmd)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

func TestImportCmd_MissingPathFlag(t *testing.T) {
	_, err := executeImportCmd(t, []string{"import", "--file", "secrets.json"})
	if err == nil {
		t.Error("expected error for missing --path flag")
	}
}

func TestImportCmd_MissingFileFlag(t *testing.T) {
	_, err := executeImportCmd(t, []string{"import", "--path", "myapp/config"})
	if err == nil {
		t.Error("expected error for missing --file flag")
	}
}

func TestImportCmd_DefaultFlags(t *testing.T) {
	f := importCmd.Flags().Lookup("format")
	if f == nil {
		t.Fatal("--format flag not registered")
	}
	if f.DefValue != "json" {
		t.Errorf("expected default format %q, got %q", "json", f.DefValue)
	}

	m := importCmd.Flags().Lookup("mount")
	if m == nil {
		t.Fatal("--mount flag not registered")
	}
	if m.DefValue != "secret" {
		t.Errorf("expected default mount %q, got %q", "secret", m.DefValue)
	}

	dr := importCmd.Flags().Lookup("dry-run")
	if dr == nil {
		t.Fatal("--dry-run flag not registered")
	}
	if dr.DefValue != "false" {
		t.Errorf("expected dry-run default false, got %q", dr.DefValue)
	}
}

func TestImportCmd_FlagRegistration(t *testing.T) {
	for _, name := range []string{"path", "file", "format", "mount", "dry-run"} {
		if importCmd.Flags().Lookup(name) == nil {
			t.Errorf("expected flag --%s to be registered", name)
		}
	}
}

func writeTempImportFile(t *testing.T, name string, data interface{}) string {
	t.Helper()
	b, _ := json.Marshal(data)
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, b, 0600); err != nil {
		t.Fatalf("writeTempImportFile: %v", err)
	}
	return p
}
