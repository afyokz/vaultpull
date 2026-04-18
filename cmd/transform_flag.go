package cmd

import (
	"github.com/spf13/cobra"
	"github.com/user/vaultpull/internal/config"
	"github.com/user/vaultpull/internal/vault"
)

// transformFlags holds CLI flag values for key transformation.
type transformFlags struct {
	Prefix    string
	Suffix    string
	Uppercase bool
	Lowercase bool
}

// registerTransformFlags attaches transform flags to a command.
func registerTransformFlags(cmd *cobra.Command, f *transformFlags) {
	cmd.Flags().StringVar(&f.Prefix, "key-prefix", "", "prefix to add to every secret key")
	cmd.Flags().StringVar(&f.Suffix, "key-suffix", "", "suffix to add to every secret key")
	cmd.Flags().BoolVar(&f.Uppercase, "uppercase", false, "convert all keys to uppercase")
	cmd.Flags().BoolVar(&f.Lowercase, "lowercase", false, "convert all keys to lowercase")
}

// buildTransformRule merges CLI flags over config-level transform settings.
func buildTransformRule(f *transformFlags, base *config.TransformConfig) (*vault.TransformRule, error) {
	merged := &config.TransformConfig{}
	if base != nil {
		*merged = *base
	}
	if f.Prefix != "" {
		merged.Prefix = f.Prefix
	}
	if f.Suffix != "" {
		merged.Suffix = f.Suffix
	}
	if f.Uppercase {
		merged.Uppercase = true
		merged.Lowercase = false
	}
	if f.Lowercase {
		merged.Lowercase = true
		merged.Uppercase = false
	}
	return merged.BuildRule()
}
