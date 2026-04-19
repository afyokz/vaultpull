package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/vault"
)

func registerRedactFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("redact", false, "Redact sensitive values before writing")
	cmd.Flags().Bool("redact-defaults", true, "Include default sensitive patterns (password, token, etc.)")
	cmd.Flags().StringSlice("redact-pattern", nil, "Additional key patterns to redact (regex)")
	cmd.Flags().String("redact-replacement", vault.DefaultRedactReplacement, "Replacement string for redacted values")
}

func buildRedactRules(cmd *cobra.Command) ([]*vault.RedactRule, error) {
	enabled, _ := cmd.Flags().GetBool("redact")
	useDefaults, _ := cmd.Flags().GetBool("redact-defaults")
	patterns, _ := cmd.Flags().GetStringSlice("redact-pattern")
	replacement, _ := cmd.Flags().GetString("redact-replacement")

	cfg := &config.RedactConfig{
		Enabled:     enabled,
		UseDefaults: useDefaults,
		Patterns:    patterns,
		Replacement: replacement,
	}
	return config.BuildRedactRules(cfg)
}
