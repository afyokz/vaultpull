package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/audit"
	"vaultpull/internal/envwriter"
	"vaultpull/internal/vault"
)

var (
	importFile   string
	importFormat string
	importOutput string
	importDryRun bool
)

func init() {
	importCmd := &cobra.Command{
		Use:   "import",
		Short: "Import secrets from a local file into a .env file",
		Long: `Read secrets from a dotenv, JSON, or shell-export file and write them
into a local .env file using the same merge and backup logic as 'pull'.`,
		RunE: runImport,
	}

	importCmd.Flags().StringVarP(&importFile, "file", "f", "-", "Source file path (use '-' for stdin)")
	importCmd.Flags().StringVar(&importFormat, "format", "dotenv", "Source format: dotenv, json, shell")
	importCmd.Flags().StringVarP(&importOutput, "output", "o", ".env", "Destination .env file")
	importCmd.Flags().BoolVar(&importDryRun, "dry-run", false, "Print secrets without writing to disk")

	rootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, _ []string) error {
	format, err := vault.ParseImportFormat(importFormat)
	if err != nil {
		return err
	}

	secrets, err := vault.ImportSecrets(importFile, format)
	if err != nil {
		return fmt.Errorf("import failed: %w", err)
	}

	if len(secrets) == 0 {
		cmd.Println("No secrets found in source.")
		return nil
	}

	if importDryRun {
		cmd.Printf("Dry-run: %d secret(s) would be written to %s\n", len(secrets), importOutput)
		for k, v := range secrets {
			cmd.Printf("  %s=%s\n", k, v)
		}
		return nil
	}

	if err := envwriter.Write(importOutput, secrets); err != nil {
		return fmt.Errorf("write %q: %w", importOutput, err)
	}

	logger, _ := audit.NewLogger(".vaultpull_audit.log")
	if logger != nil {
		_ = logger.Log("import", importOutput, nil)
	}

	cmd.Printf("Imported %d secret(s) into %s\n", len(secrets), importOutput)
	return nil
}

// ensure os is imported for potential future stdin handling
var _ = os.Stdin
