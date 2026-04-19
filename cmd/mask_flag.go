package cmd

import (
	"github.com/spf13/cobra"
	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/vault"
)

func registerMaskFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("mask", false, "Mask sensitive secret values in output")
	cmd.Flags().StringSlice("mask-keys", nil, "Comma-separated list of keys to mask")
	cmd.Flags().Int("mask-prefix", 2, "Number of leading characters to reveal")
	cmd.Flags().Int("mask-suffix", 2, "Number of trailing characters to reveal")
	cmd.Flags().String("mask-replacement", "****", "Replacement string for masked portion")
}

func buildMaskRules(cmd *cobra.Command) (*vault.MaskOption, []string) {
	enabled, _ := cmd.Flags().GetBool("mask")
	keys, _ := cmd.Flags().GetStringSlice("mask-keys")
	prefix, _ := cmd.Flags().GetInt("mask-prefix")
	suffix, _ := cmd.Flags().GetInt("mask-suffix")
	replacement, _ := cmd.Flags().GetString("mask-replacement")

	cfg := &config.MaskConfig{
		Enabled:     enabled,
		Keys:        keys,
		ShowPrefix:  prefix,
		ShowSuffix:  suffix,
		Replacement: replacement,
	}
	opt := config.BuildMaskOption(cfg)
	return opt, keys
}
