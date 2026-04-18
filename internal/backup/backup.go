package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Manager handles creating and cleaning up .env backup files.
type Manager struct {
	MaxBackups int
}

// NewManager returns a Manager with a default max backup count.
func NewManager(maxBackups int) *Manager {
	if maxBackups <= 0 {
		maxBackups = 5
	}
	return &Manager{MaxBackups: maxBackups}
}

// Create copies src to a timestamped backup file alongside it.
// Returns the backup path or an error.
func (m *Manager) Create(src string) (string, error) {
	data, err := os.ReadFile(src)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // nothing to back up
		}
		return "", fmt.Errorf("backup: read source: %w", err)
	}

	timestamp := time.Now().UTC().Format("20060102T150405Z")
	dir := filepath.Dir(src)
	base := filepath.Base(src)
	backupPath := filepath.Join(dir, fmt.Sprintf("%s.%s.bak", base, timestamp))

	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return "", fmt.Errorf("backup: write backup: %w", err)
	}

	if err := m.pruneOld(dir, base); err != nil {
		return backupPath, fmt.Errorf("backup: prune: %w", err)
	}

	return backupPath, nil
}

// pruneOld removes oldest backup files when count exceeds MaxBackups.
func (m *Manager) pruneOld(dir, base string) error {
	pattern := filepath.Join(dir, base+".*.bak")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	if len(matches) <= m.MaxBackups {
		return nil
	}
	// Glob returns sorted order; oldest timestamps sort first.
	for _, f := range matches[:len(matches)-m.MaxBackups] {
		if err := os.Remove(f); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}
