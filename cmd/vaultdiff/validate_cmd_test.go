package main

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func executeValidateCmd(args []string) (string, error) {
	buf := &bytes.Buffer{}
	root := &cobra.Command{Use: "vaultdiff"}
	root.AddCommand(validateCmd)
	validateCmd.SetOut(buf)
	root.SetOut(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func TestValidateCmd_MissingPathFlag(t *testing.T) {
	_, err := executeValidateCmd([]string{"validate"})
	if err == nil {
		t.Fatal("expected error when --path is missing")
	}
}

func TestValidateCmd_InvalidPatternFlag(t *testing.T) {
	// pattern without '=' separator should return an error
	_, err := executeValidateCmd([]string{
		"validate",
		"--path", "secret/myapp",
		"--pattern", "no-equals-sign",
	})
	if err == nil {
		t.Fatal("expected error for malformed pattern flag")
	}
}

func TestValidateCmd_DefaultFlags(t *testing.T) {
	mount, err := validateCmd.Flags().GetString("mount")
	if err != nil {
		t.Fatalf("unexpected error reading mount flag: %v", err)
	}
	if mount != "secret" {
		t.Errorf("expected default mount=secret, got %q", mount)
	}
}

func TestValidateCmd_FlagRegistration(t *testing.T) {
	expected := []string{"path", "mount", "require", "pattern"}
	for _, name := range expected {
		if validateCmd.Flags().Lookup(name) == nil {
			t.Errorf("expected flag --%s to be registered", name)
		}
	}
}

func TestValidateCmd_RequireFlag(t *testing.T) {
	f := validateCmd.Flags().Lookup("require")
	if f == nil {
		t.Fatal("expected --require flag to exist")
	}
	if f.DefValue != "[]" {
		t.Errorf("unexpected default for --require: %q", f.DefValue)
	}
}
