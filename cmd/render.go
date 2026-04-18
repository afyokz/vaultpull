package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/template"
	"vaultpull/internal/vault"
)

var renderCmd = &cobra.Command{
	Use:   "render [template-file]",
	Short: "Render a template file using secrets from Vault",
	Args:  cobra.ExactArgs(1),
	RunE:  runRender,
}

var renderOutput string

func init() {
	renderCmd.Flags().StringVarP(&renderOutput, "output", "o", "", "Output file path (default: stdout)")
	rootCmd.AddCommand(renderCmd)
}

func runRender(cmd *cobra.Command, args []string) error {
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

	r := template.NewRenderer()
	result, missing, err := r.RenderFile(args[0], merged)
	if err != nil {
		return err
	}

	if len(missing) > 0 {
		fmt.Fprintf(os.Stderr, "warning: unresolved placeholders: %s\n", strings.Join(missing, ", "))
	}

	if renderOutput == "" {
		fmt.Print(result)
		return nil
	}

	if err := os.WriteFile(renderOutput, []byte(result), 0600); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	fmt.Fprintf(os.Stderr, "rendered output written to %s\n", renderOutput)
	return nil
}
