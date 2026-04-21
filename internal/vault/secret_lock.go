package vault

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// LockOption configures lock behavior.
type LockOption struct {
	LockDir  string
	TTL      time.Duration
	StaleTTL time.Duration
}

// DefaultLockOption returns sensible defaults.
func DefaultLockOption() LockOption {
	return LockOption{
		LockDir:  os.TempDir(),
		TTL:      30 * time.Second,
		StaleTTL: 5 * time.Minute,
	}
}

// LockManager manages file-based advisory locks for vault pull operations.
type LockManager struct {
	opt LockOption
}

// NewLockManager creates a new LockManager.
func NewLockManager(opt LockOption) *LockManager {
	return &LockManager{opt: opt}
}

func (m *LockManager) lockPath(key string) string {
	return filepath.Join(m.opt.LockDir, fmt.Sprintf("vaultpull-%s.lock", key))
}

// Acquire attempts to create a lock file. Returns an error if already locked.
func (m *LockManager) Acquire(key string) error {
	path := m.lockPath(key)

	if info, err := os.Stat(path); err == nil {
		age := time.Since(info.ModTime())
		if age < m.opt.StaleTTL {
			return fmt.Errorf("lock already held for %q (age: %s)", key, age.Round(time.Second))
		}
		// stale lock — remove it
		_ = os.Remove(path)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("lock already held for %q", key)
		}
		return fmt.Errorf("acquire lock: %w", err)
	}
	_ = f.Close()
	return nil
}

// Release removes the lock file.
func (m *LockManager) Release(key string) error {
	path := m.lockPath(key)
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("release lock: %w", err)
	}
	return nil
}

// IsLocked reports whether a lock is currently held.
func (m *LockManager) IsLocked(key string) bool {
	path := m.lockPath(key)
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return time.Since(info.ModTime()) < m.opt.StaleTTL
}
