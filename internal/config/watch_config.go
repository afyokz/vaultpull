package config

import (
	"fmt"
	"time"

	"github.com/your-org/vaultpull/internal/vault"
)

// WatchConfig holds configuration for the watch command.
type WatchConfig struct {
	Interval  string `yaml:"interval"`
	MaxErrors int    `yaml:"max_errors"`
}

// BuildWatchOption converts WatchConfig into a vault.WatchOption.
// It falls back to defaults for zero values.
func BuildWatchOption(cfg *WatchConfig) (vault.WatchOption, error) {
	opt := vault.DefaultWatchOption()
	if cfg == nil {
		return opt, nil
	}
	if cfg.Interval != "" {
		d, err := time.ParseDuration(cfg.Interval)
		if err != nil {
			return opt, fmt.Errorf("invalid watch interval %q: %w", cfg.Interval, err)
		}
		if d <= 0 {
			return opt, fmt.Errorf("watch interval must be positive, got %s", cfg.Interval)
		}
		opt.Interval = d
	}
	if cfg.MaxErrors > 0 {
		opt.MaxErrors = cfg.MaxErrors
	}
	return opt, nil
}
