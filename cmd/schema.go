package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpull/internal/config"
	"github.com/yourorg/vaultpull/internal/vault"
)

var schemaRules []string

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Validate fetched secrets against schema rules",
	RunE:  runSchema,
}

func init() {
	schemaCmd.Flags().StringArrayVar(&schemaRules, "rule", nil,
		`Schema rule in the form KEY:PATTERN or KEY:PATTERN:required (repeatable)`)
	schemaCmd.Flags().String("config", "vaultpull.yaml", "Path to config file")
	rootCmd.AddCommand(schemaCmd)
}

func runSchema(cmd *cobra.Command, _ []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	client, err := vault.NewClient(cfg.Address, cfg.Token)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	secrets, err := vault.FetchSecrets(client, cfg.Secrets)
	if err != nil {
		return fmt.Errorf("fetching secrets: %w", err)
	}
	merged := vault.MergeSecrets(secrets)

	var rules []vault.SchemaRule
	for _, raw := range schemaRules {
		r, err := vault.ParseSchemaRule(raw)
		if err != nil {
			return fmt.Errorf("invalid schema rule %q: %w", raw, err)
		}
		rules = append(rules, r)
	}

	violations := vault.ValidateSchema(merged, rules)
	fmt.Println(vault.FormatSchemaReport(violations))
	if len(violations) > 0 {
		os.Exit(1)
	}
	return nil
}
