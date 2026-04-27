package config

import (
	"fmt"

	"github.com/your-org/vaultpull/internal/vault"
)

// RollbackConfig holds configuration for rollback behaviour.
type RollbackConfig struct {
	Enabled  bool   `yaml:"enabled"`
	MaxSize  int    `yaml:"max_size"`
	AutoSave bool   `yaml:"auto_save"`
	Label    string `yaml:"label"`
}

// BuildRollbackStore constructs a RollbackStore from config.
// Returns nil if cfg is nil or disabled.
func BuildRollbackStore(cfg *RollbackConfig) (*vault.RollbackStore, error) {
	if cfg == nil || !cfg.Enabled {
		return nil, nil
	}
	if cfg.MaxSize < 0 {
		return nil, fmt.Errorf("rollback max_size must be non-negative, got %d", cfg.MaxSize)
	}
	maxSize := cfg.MaxSize
	if maxSize == 0 {
		maxSize = vault.DefaultRollbackMaxSize
	}
	return vault.NewRollbackStore(maxSize), nil
}
