package config

import (
	"fmt"
	"os"
	"time"

	"github.com/your-org/vaultpull/internal/vault"
)

// LockConfig holds lock-related configuration from the YAML config file.
type LockConfig struct {
	Enabled  bool   `yaml:"enabled"`
	LockDir  string `yaml:"lock_dir"`
	StaleTTL string `yaml:"stale_ttl"`
}

// BuildLockOption converts a LockConfig into a vault.LockOption.
// If cfg is nil, defaults are returned with locking disabled.
func BuildLockOption(cfg *LockConfig) (vault.LockOption, bool, error) {
	if cfg == nil || !cfg.Enabled {
		return vault.DefaultLockOption(), false, nil
	}

	opt := vault.DefaultLockOption()

	if cfg.LockDir != "" {
		if err := os.MkdirAll(cfg.LockDir, 0700); err != nil {
			return opt, false, fmt.Errorf("lock_dir: %w", err)
		}
		opt.LockDir = cfg.LockDir
	}

	if cfg.StaleTTL != "" {
		d, err := time.ParseDuration(cfg.StaleTTL)
		if err != nil {
			return opt, false, fmt.Errorf("stale_ttl: %w", err)
		}
		if d <= 0 {
			return opt, false, fmt.Errorf("stale_ttl must be positive")
		}
		opt.StaleTTL = d
	}

	return opt, true, nil
}
