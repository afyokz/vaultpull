package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/vaultpull/internal/vault"
	"github.com/user/vaultpull/internal/envwriter"
)

var secretPath string

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull secrets from Vault and write to .env file",
	RunE:  runPull,
}

func init() {
	pullCmd.Flags().StringVar(&secretPath, "path", "", "Vault secret path (required)")
	_ = pullCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(pullCmd)
}

func runPull(cmd *cobra.Command, args []string) error {
	addr := vaultAddr
	if addr == "" {
		addr = os.Getenv("VAULT_ADDR")
	}
	token := vaultToken
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	if addr == "" {
		return fmt.Errorf("vault address is required: set --vault-addr or VAULT_ADDR")
	}
	if token == "" {
		return fmt.Errorf("vault token is required: set --vault-token or VAULT_TOKEN")
	}

	client, err := vault.NewClient(addr, token)
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	secrets, err := client.GetSecrets(secretPath)
	if err != nil {
		return fmt.Errorf("failed to fetch secrets: %w", err)
	}

	if err := envwriter.Write(envFile, secrets); err != nil {
		return fmt.Errorf("failed to write env file: %w", err)
	}

	fmt.Printf("✓ Wrote %d secrets to %s\n", len(secrets), envFile)
	return nil
}
