package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/example/vaultpull/internal/config"
	"github.com/example/vaultpull/internal/vault"
)

var (
	checksumVerify string
	checksumQuiet  bool
)

func init() {
	checksumCmd := &cobra.Command{
		Use:   "checksum",
		Short: "Compute or verify a checksum of fetched secrets",
		Long: `Fetches secrets from Vault and prints a SHA-256 digest.
Pass --verify <digest> to assert the current secrets match a known digest.
Exits with code 1 when verification fails.`,
		RunE: runChecksum,
	}

	checksumCmd.Flags().StringVar(&checksumVerify, "verify", "", "expected digest to verify against (sha256:<hex> or bare hex)")
	checksumCmd.Flags().BoolVar(&checksumQuiet, "quiet", false, "only print the digest, no report")

	rootCmd.AddCommand(checksumCmd)
}

func runChecksum(cmd *cobra.Command, args []string) error {
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

	merged := vault.MergeSecrets(secrets)
	result := vault.ComputeChecksum(merged)

	if checksumVerify != "" {
		if !vault.VerifyChecksum(merged, checksumVerify) {
			fmt.Fprintf(os.Stderr, "checksum mismatch\n  got:      %s\n  expected: %s\n",
				result.Digest, checksumVerify)
			os.Exit(1)
		}
		if !checksumQuiet {
			fmt.Println("Checksum verified OK")
		}
		return nil
	}

	if checksumQuiet {
		fmt.Println(result.Digest)
		return nil
	}

	fmt.Print(vault.FormatChecksumReport(result))
	return nil
}
