package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpull/internal/vault"
)

var (
	ttlWarnHours  int
	ttlErrorHours int
)

var ttlCmd = &cobra.Command{
	Use:   "ttl",
	Short: "Check TTL status of secrets from Vault metadata",
	RunE:  runTTL,
}

func init() {
	ttlCmd.Flags().IntVar(&ttlWarnHours, "warn-hours", 24, "Hours remaining before a secret is considered a warning")
	ttlCmd.Flags().IntVar(&ttlErrorHours, "error-hours", 0, "Hours remaining before a secret is considered expired")
	rootCmd.AddCommand(ttlCmd)
}

func runTTL(cmd *cobra.Command, args []string) error {
	client, err := vault.NewClient()
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	policy := vault.TTLPolicy{
		WarnThreshold:  time.Duration(ttlWarnHours) * time.Hour,
		ErrorThreshold: time.Duration(ttlErrorHours) * time.Hour,
	}

	meta, err := vault.FetchSecretMeta(client, args)
	if err != nil {
		return fmt.Errorf("fetching secret metadata: %w", err)
	}

	expiries := make(map[string]time.Time, len(meta))
	for k, m := range meta {
		expiries[k] = m.ExpiresAt
	}

	results := vault.CheckTTL(expiries, policy)
	fmt.Print(vault.FormatTTLReport(results))

	for _, r := range results {
		if r.Status == "expired" {
			os.Exit(2)
		}
	}
	return nil
}
