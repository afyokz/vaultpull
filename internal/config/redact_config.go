package config

import (
	"fmt"

	"github.com/yourusername/vaultpull/internal/vault"
)

// RedactConfig holds redaction settings from the config file.
type RedactConfig struct {
	Enabled         bool     `yaml:"enabled"`
	UseDefaults     bool     `yaml:"use_defaults"`
	Patterns        []string `yaml:"patterns"`
	Replacement     string   `yaml:"replacement"`
}

// BuildRedactRules converts RedactConfig into RedactRule slice.
func BuildRedactRules(cfg *RedactConfig) ([]*vault.RedactRule, error) {
	if cfg == nil || !cfg.Enabled {
		return nil, nil
	}
	var rules []*vault.RedactRule
	if cfg.UseDefaults {
		rules = append(rules, vault.DefaultSensitivePatterns()...)
	}
	for _, p := range cfg.Patterns {
		r, err := vault.ParseRedactRule(p, cfg.Replacement)
		if err != nil {
			return nil, fmt.Errorf("invalid redact pattern %q: %w", p, err)
		}
		rules = append(rules, r)
	}
	return rules, nil
}
