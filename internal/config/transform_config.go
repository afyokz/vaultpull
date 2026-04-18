package config

import (
	"fmt"

	"github.com/user/vaultpull/internal/vault"
)

// TransformConfig holds transform settings from the config file.
type TransformConfig struct {
	Prefix    string `yaml:"prefix"`
	Suffix    string `yaml:"suffix"`
	Uppercase bool   `yaml:"uppercase"`
	Lowercase bool   `yaml:"lowercase"`
}

// BuildRule converts TransformConfig into a vault.TransformRule.
func (tc *TransformConfig) BuildRule() (*vault.TransformRule, error) {
	if tc == nil {
		return vault.NewTransformRule()
	}
	var opts []vault.TransformOption
	if tc.Prefix != "" {
		opts = append(opts, vault.WithPrefix(tc.Prefix))
	}
	if tc.Suffix != "" {
		opts = append(opts, vault.WithSuffix(tc.Suffix))
	}
	if tc.Uppercase {
		opts = append(opts, vault.WithUppercase())
	}
	if tc.Lowercase {
		opts = append(opts, vault.WithLowercase())
	}
	rule, err := vault.NewTransformRule(opts...)
	if err != nil {
		return nil, fmt.Errorf("transform config: %w", err)
	}
	return rule, nil
}
