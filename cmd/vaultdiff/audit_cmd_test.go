package main

import (
	"bytes"
	"os"
	"testing"
)

func TestAuditCmd_FileNotFound(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"audit", "--log", "/nonexistent/audit.log"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when log file does not exist")
	}
}

func TestAuditCmd_InvalidSinceDuration(t *testing.T) {
	f, err := os.CreateTemp("", "audit-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	// Write empty JSON array so ReadEntries succeeds.
	_, _ = f.WriteString("[]\n")
	f.Close()

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{
		"audit",
		"--log", f.Name(),
		"--since", "notaduration",
	})

	err = rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid --since duration")
	}
}

func TestAuditCmd_DefaultFlags(t *testing.T) {
	if got := auditCmd.Flags().Lookup("env"); got == nil {
		t.Error("expected --env flag to be registered")
	}
	if got := auditCmd.Flags().Lookup("path-prefix"); got == nil {
		t.Error("expected --path-prefix flag to be registered")
	}
	if got := auditCmd.Flags().Lookup("min-changes"); got == nil {
		t.Error("expected --min-changes flag to be registered")
	}
}
