package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpull/internal/audit"
	"github.com/yourusername/vaultpull/internal/backup"
	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/diff"
	"github.com/yourusername/vaultpull/internal/envwriter"
	"github.com/yourusername/vaultpull/internal/prompt"
	"github.com/yourusername/vaultpull/internal/vault"
)

var (
	cfgFile   string
	autoYes   bool
	dryRun    bool
)

func init() {
	pullCmd.Flags().StringVarP(&cfgFile, "config", "c", "vaultpull.yaml", "config file path")
	pullCmd.Flags().BoolVarP(&autoYes, "yes", "y", false, "auto-confirm changes without prompting")
	pullCmd.Flags().BoolVar(&dryRun, "dry-run", false, "show changes without writing")
	rootCmd.AddCommand(pullCmd)
}

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Fetch secrets from Vault and write to .env file",
	RunE:  runPull,
}

func runPull(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	client, err := vault.NewClient(cfg.Vault.Address, cfg.Vault.Token)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secrets, err := vault.FetchSecrets(client, cfg.Secrets)
	if err != nil {
		return fmt.Errorf("fetch secrets: %w", err)
	}

	result, err := diff.Compute(cfg.Output, secrets)
	if err != nil {
		return fmt.Errorf("compute diff: %w", err)
	}

	if dryRun {
		fmt.Println(result.Summary())
		return nil
	}

	if !autoYes {
		confirmer := prompt.NewConfirmer()
		ok, err := confirmer.ConfirmDiff(result)
		if err != nil {
			return fmt.Errorf("prompt: %w", err)
		}
		if !ok {
			fmt.Println("Aborted.")
			return nil
		}
	}

	bm := backup.NewManager(cfg.Output)
	if err := bm.Create(); err != nil {
		return fmt.Errorf("backup: %w", err)
	}

	if err := envwriter.Write(cfg.Output, secrets); err != nil {
		return fmt.Errorf("write env: %w", err)
	}

	auditLog := audit.NewLogger(cfg.AuditLog)
	if err := auditLog.Log(cfg.Output, result, nil); err != nil {
		fmt.Fprintf(os.Stderr, "warning: audit log failed: %v\n", err)
	}

	fmt.Println("Done.", result.Summary())
	return nil
}
