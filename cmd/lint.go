package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/vault"
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint secrets fetched from Vault for common issues",
	RunE:  runLint,
}

func init() {
	lintCmd.Flags().StringP("config", "c", "vaultpull.yaml", "Path to config file")
	rootCmd.AddCommand(lintCmd)
}

func runLint(cmd *cobra.Command, args []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	client, err := vault.NewClient(cfg.Address, cfg.Token)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secrets, err := vault.FetchSecrets(client, cfg.Secrets)
	if err != nil {
		return fmt.Errorf("fetch secrets: %w", err)
	}

	merged := vault.MergeSecrets(secrets)
	issues := vault.LintSecrets(merged)
	report := vault.FormatLintReport(issues)
	fmt.Println(report)

	for _, issue := range issues {
		if issue.Severity == vault.SeverityError {
			os.Exit(1)
		}
	}
	return nil
}
