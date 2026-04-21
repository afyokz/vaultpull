package config

import (
	"fmt"

	"vaultpull/internal/vault"
)

// ImportConfig holds configuration for the import command sourced from the
// YAML config file.
type ImportConfig struct {
	File   string `yaml:"file"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

// BuildImportFormat validates and returns the ImportFormat for the given
// config block. If cfg is nil or Format is empty, the default dotenv format
// is returned.
func BuildImportFormat(cfg *ImportConfig) (vault.ImportFormat, error) {
	if cfg == nil || cfg.Format == "" {
		return vault.ImportFormatDotenv, nil
	}
	fmt_, err := vault.ParseImportFormat(cfg.Format)
	if err != nil {
		return "", fmt.Errorf("config: import format: %w", err)
	}
	return fmt_, nil
}

// ResolveImportOutput returns the output path from the config, falling back to
// the provided default when the config is nil or the field is empty.
func ResolveImportOutput(cfg *ImportConfig, defaultPath string) string {
	if cfg == nil || cfg.Output == "" {
		return defaultPath
	}
	return cfg.Output
}
