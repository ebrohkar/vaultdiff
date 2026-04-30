package main

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func executeMergeCmd(args []string) (string, error) {
	buf := &bytes.Buffer{}
	root := &cobra.Command{Use: "vaultdiff"}
	mergeCmd.ResetFlags()
	mergeCmd.Flags().String("source", "", "Source secret path")
	mergeCmd.Flags().String("dest", "", "Destination secret path")
	mergeCmd.Flags().String("mount", "secret", "KV mount name")
	mergeCmd.Flags().String("strategy", "source", "Conflict strategy")
	root.AddCommand(mergeCmd)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func TestMergeCmd_MissingSourceFlag(t *testing.T) {
	_, err := executeMergeCmd([]string{"merge", "--dest", "env/prod/app"})
	if err == nil {
		t.Error("expected error when --source flag is missing")
	}
}

func TestMergeCmd_MissingDestFlag(t *testing.T) {
	_, err := executeMergeCmd([]string{"merge", "--source", "env/dev/app"})
	if err == nil {
		t.Error("expected error when --dest flag is missing")
	}
}

func TestMergeCmd_DefaultFlags(t *testing.T) {
	cmd := mergeCmd
	cmd.ResetFlags()
	cmd.Flags().String("source", "", "")
	cmd.Flags().String("dest", "", "")
	cmd.Flags().String("mount", "secret", "")
	cmd.Flags().String("strategy", "source", "")

	mount, _ := cmd.Flags().GetString("mount")
	if mount != "secret" {
		t.Errorf("expected default mount=secret, got %q", mount)
	}

	strategy, _ := cmd.Flags().GetString("strategy")
	if strategy != "source" {
		t.Errorf("expected default strategy=source, got %q", strategy)
	}
}

func TestMergeCmd_FlagRegistration(t *testing.T) {
	cmd := mergeCmd
	cmd.ResetFlags()
	cmd.Flags().String("source", "", "")
	cmd.Flags().String("dest", "", "")
	cmd.Flags().String("mount", "secret", "")
	cmd.Flags().String("strategy", "source", "")

	for _, name := range []string{"source", "dest", "mount", "strategy"} {
		if cmd.Flags().Lookup(name) == nil {
			t.Errorf("expected flag --%s to be registered", name)
		}
	}
}
