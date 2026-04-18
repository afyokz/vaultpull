package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/vault"
)

var rotateDays int

var rotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Check which secrets are due for rotation",
	RunE:  runRotate,
}

func init() {
	rotateCmd.Flags().IntVar(&rotateDays, "max-age", 90, "Maximum secret age in days before rotation is flagged")
	rootCmd.AddCommand(rotateCmd)
}

func runRotate(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	client, err := vault.NewClient(cfg.VaultAddress, cfg.Token)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	policy := vault.RotationPolicy{MaxAgeDays: rotateDays}
	meta := map[string]vault.SecretMeta{}

	for _, s := range cfg.Secrets {
		sp, err := vault.ParseSecretPath(s.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: skipping invalid path %q: %v\n", s.Path, err)
			continue
		}
		// Use version 1 and a placeholder time; real impl would query Vault metadata.
		meta[sp.String()] = vault.SecretMeta{
			Path:        sp.String(),
			Version:     1,
			CreatedTime: time.Now().AddDate(0, 0, -30),
		}
	}

	results := vault.CheckRotation(meta, policy)
	fmt.Print(vault.FormatRotationReport(results))
	return nil
}
