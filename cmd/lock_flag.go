package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpull/internal/vault"
)

var (
	lockEnabled  bool
	lockDir      string
	lockStaleTTL string
)

// registerLockFlags attaches lock-related flags to a command.
func registerLockFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&lockEnabled, "lock", false, "enable file-based advisory locking")
	cmd.Flags().StringVar(&lockDir, "lock-dir", "", "directory for lock files (default: system temp)")
	cmd.Flags().StringVar(&lockStaleTTL, "lock-stale-ttl", "5m", "age after which a lock is considered stale")
}

// buildLockManager constructs a LockManager from CLI flags.
// Returns nil, false if locking is not enabled.
func buildLockManager() (*vault.LockManager, bool, error) {
	if !lockEnabled {
		return nil, false, nil
	}

	opt := vault.DefaultLockOption()

	if lockDir != "" {
		opt.LockDir = lockDir
	}

	if lockStaleTTL != "" {
		import_d, err := parseDuration(lockStaleTTL)
		if err != nil {
			return nil, false, fmt.Errorf("--lock-stale-ttl: %w", err)
		}
		opt.StaleTTL = import_d
	}

	return vault.NewLockManager(opt), true, nil
}

// acquireOrDie acquires a lock or exits with a fatal message.
func acquireOrDie(m *vault.LockManager, key string) {
	if m == nil {
		return
	}
	if err := m.Acquire(key); err != nil {
		log.Fatalf("cannot acquire lock: %v", err)
	}
}
