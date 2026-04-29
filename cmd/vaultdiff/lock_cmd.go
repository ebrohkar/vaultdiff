package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/example/vaultdiff/internal/config"
	"github.com/example/vaultdiff/internal/vault"
)

var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Lock or unlock a secret path to prevent modifications",
}

var lockSecretCmd = &cobra.Command{
	Use:   "set",
	Short: "Advisorily lock a secret path",
	RunE:  runLock,
}

var unlockSecretCmd = &cobra.Command{
	Use:   "unset",
	Short: "Remove an advisory lock from a secret path",
	RunE:  runUnlock,
}

func init() {
	lockSecretCmd.Flags().String("path", "", "Secret path to lock (required)")
	lockSecretCmd.Flags().String("mount", "secret", "KV mount point")
	lockSecretCmd.Flags().String("reason", "", "Reason for locking")
	lockSecretCmd.Flags().Duration("ttl", 24*time.Hour, "Lock TTL duration")
	_ = lockSecretCmd.MarkFlagRequired("path")

	unlockSecretCmd.Flags().String("path", "", "Secret path to unlock (required)")
	unlockSecretCmd.Flags().String("mount", "secret", "KV mount point")
	_ = unlockSecretCmd.MarkFlagRequired("path")

	lockCmd.AddCommand(lockSecretCmd)
	lockCmd.AddCommand(unlockSecretCmd)
	rootCmd.AddCommand(lockCmd)
}

func runLock(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	path, _ := cmd.Flags().GetString("path")
	mount, _ := cmd.Flags().GetString("mount")
	reason, _ := cmd.Flags().GetString("reason")
	ttl, _ := cmd.Flags().GetDuration("ttl")

	result, err := vault.LockSecret(client, vault.LockOptions{
		Path:   path,
		Mount:  mount,
		Reason: reason,
		TTL:    ttl,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Locked %s/%s until %s\n", result.Mount, result.Path, result.ExpiresAt.Format(time.RFC3339))
	return nil
}

func runUnlock(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	path, _ := cmd.Flags().GetString("path")
	mount, _ := cmd.Flags().GetString("mount")

	_, err = vault.UnlockSecret(client, path, mount)
	if err != nil {
		return err
	}

	fmt.Printf("Unlocked %s/%s\n", mount, path)
	return nil
}
