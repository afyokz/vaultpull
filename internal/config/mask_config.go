package config

import "github.com/your-org/vaultpull/internal/vault"

// MaskConfig holds configuration for secret masking.
type MaskConfig struct {
	Enabled     bool     `yaml:"enabled"`
	Keys        []string `yaml:"keys"`
	ShowPrefix  int      `yaml:"show_prefix"`
	ShowSuffix  int      `yaml:"show_suffix"`
	Replacement string   `yaml:"replacement"`
}

// BuildMaskOption converts MaskConfig into a vault.MaskOption.
// Returns nil if masking is disabled or config is nil.
func BuildMaskOption(cfg *MaskConfig) *vault.MaskOption {
	if cfg == nil || !cfg.Enabled {
		return nil
	}
	opt := vault.DefaultMaskOption
	if cfg.ShowPrefix > 0 {
		opt.ShowPrefix = cfg.ShowPrefix
	}
	if cfg.ShowSuffix >= 0 && cfg.ShowSuffix != opt.ShowSuffix {
		opt.ShowSuffix = cfg.ShowSuffix
	}
	if cfg.Replacement != "" {
		opt.Replacement = cfg.Replacement
	}
	return &opt
}
