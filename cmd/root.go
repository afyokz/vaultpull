package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	vaultAddr  string
	vaultToken string
	envFile    string
)

var rootCmd = &cobra.Command{
	Use:   "vaultpull",
	Short: "Sync secrets from HashiCorp Vault into local .env files",
	Long: `vaultpull is a CLI tool that fetches secrets from HashiCorp Vault
and writes them safely into local .env files for development use.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&vaultAddr, "vault-addr", "", "Vault server address (overrides VAULT_ADDR env)")
	rootCmd.PersistentFlags().StringVar(&vaultToken, "vault-token", "", "Vault token (overrides VAULT_TOKEN env)")
	rootCmd.PersistentFlags().StringVar(&envFile, "env-file", ".env", "Target .env file path")
}
