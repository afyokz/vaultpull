package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/audit"
	"vaultpull/internal/backup"
	"vaultpull/internal/config"
	"vaultpull/internal/diff"
	"vaultpull/internal/envwriter"
	"vaultpull/internal/vault"
)

var (
	cfgFile  string
	dryRun   bool
	showDiff bool
)

func init() {
	pullCmd := &cobra.Command{
		Use:   "pull",
		Short: "Sync secrets from Vault into a local .env file",
		RunE:  runPull,
	}
	pullCmd.Flags().StringVarP(&cfgFile, "config", "c", "vaultpull.yaml", "config file path")
	pullCmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview changes without writing")
	pullCmd.Flags().BoolVar(&showDiff, "diff", false, "show diff of changes")
	rootCmd.AddCommand(pullCmd)
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

	auditLog, err := audit.NewLogger(cfg.AuditLog)
	if err != nil {
		return fmt.Errorf("audit logger: %w", err)
	}

	if showDiff || dryRun {
		existing := map[string]string{}
		_ = envwriter.ReadEnvFile(cfg.OutputFile, existing)
		r := diff.Compute(existing, secrets)
		fmt.Fprintln(os.Stdout, r.Summary())
		for _, c := range r.Changes {
			if c.Type != diff.Unchanged {
				fmt.Fprintf(os.Stdout, "  [%s] %s\n", c.Type, c.Key)
			}
		}
		if dryRun {
			return nil
		}
	}

	bm := backup.NewManager(cfg.BackupDir, cfg.MaxBackups)
	if err := bm.Create(cfg.OutputFile); err != nil {
		log.Printf("warn: backup failed: %v", err)
	}

	if err := envwriter.Write(cfg.OutputFile, secrets); err != nil {
		_ = auditLog.Log("pull", cfg.OutputFile, err)
		return fmt.Errorf("write env: %w", err)
	}

	_ = auditLog.Log("pull", cfg.OutputFile, nil)
	fmt.Printf("Secrets written to %s\n", cfg.OutputFile)
	return nil
}
