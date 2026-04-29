package main

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func executeRenameCmd(args []string) error {
	buf := &bytes.Buffer{}
	root := &cobra.Command{Use: "vaultdiff"}
	root.AddCommand(renameCmd)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	_, err := root.ExecuteC()
	return err
}

func TestRenameCmd_MissingSourceFlag(t *testing.T) {
	err := executeRenameCmd([]string{"rename", "--dest", "app/new"})
	if err == nil {
		t.Fatal("expected error for missing --source flag")
	}
}

func TestRenameCmd_MissingDestFlag(t *testing.T) {
	err := executeRenameCmd([]string{"rename", "--source", "app/old"})
	if err == nil {
		t.Fatal("expected error for missing --dest flag")
	}
}

func TestRenameCmd_DefaultFlags(t *testing.T) {
	mount, _ := renameCmd.Flags().GetString("mount")
	if mount != "secret" {
		t.Errorf("expected default mount 'secret', got %q", mount)
	}

	delSrc, _ := renameCmd.Flags().GetBool("delete-source")
	if delSrc {
		t.Error("expected delete-source to default to false")
	}
}

func TestRenameCmd_FlagRegistration(t *testing.T) {
	for _, flag := range []string{"source", "dest", "mount", "delete-source"} {
		if renameCmd.Flags().Lookup(flag) == nil {
			t.Errorf("flag %q not registered on rename command", flag)
		}
	}
}
