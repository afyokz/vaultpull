package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/vault"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export secrets to a file or stdout in various formats",
	RunE:  runExport,
}

var (
	exportFormat string
	exportOutput string
	exportRedact bool
)

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "dotenv", "Output format: dotenv, json, shell")
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Output file path (default: stdout)")
	exportCmd.Flags().BoolVar(&exportRedact, "redact", false, "Redact secret values in output")
}

func runExport(cmd *cobra.Command, args []string) error {
	cfgFile, _ := cmd.Root().PersistentFlags().GetString("config")
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	client, err := vault.NewClient(cfg.Address, cfg.Token)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secrets, err := vault.FetchSecrets(cmd.Context(), client, cfg.Secrets)
	if err != nil {
		return fmt.Errorf("fetching secrets: %w", err)
	}

	merged := vault.MergeSecrets(secrets)

	opts := vault.ExportOptions{
		Format:  vault.ExportFormat(exportFormat),
		OutFile: exportOutput,
		Redact:  exportRedact,
	}

	if err := vault.ExportSecrets(merged, opts); err != nil {
		fmt.Fprintf(os.Stderr, "export error: %v\n", err)
		return err
	}
	return nil
}
