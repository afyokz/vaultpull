package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempLockDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "vaultpull-lock-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func newTestLockManager(t *testing.T) *LockManager {
	t.Helper()
	opt := DefaultLockOption()
	opt.LockDir = tempLockDir(t)
	opt.StaleTTL = 2 * time.Second
	return NewLockManager(opt)
}

func TestAcquire_Success(t *testing.T) {
	m := newTestLockManager(t)
	if err := m.Acquire("myenv"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	_ = m.Release("myenv")
}

func TestAcquire_DoubleAcquire(t *testing.T) {
	m := newTestLockManager(t)
	if err := m.Acquire("myenv"); err != nil {
		t.Fatal(err)
	}
	defer m.Release("myenv")

	if err := m.Acquire("myenv"); err == nil {
		t.Fatal("expected error on double acquire")
	}
}

func TestRelease_NotHeld(t *testing.T) {
	m := newTestLockManager(t)
	if err := m.Release("ghost"); err != nil {
		t.Fatalf("expected no error releasing non-existent lock, got %v", err)
	}
}

func TestIsLocked_False(t *testing.T) {
	m := newTestLockManager(t)
	if m.IsLocked("absent") {
		t.Fatal("expected false for absent lock")
	}
}

func TestIsLocked_True(t *testing.T) {
	m := newTestLockManager(t)
	_ = m.Acquire("env")
	defer m.Release("env")
	if !m.IsLocked("env") {
		t.Fatal("expected true for held lock")
	}
}

func TestAcquire_StaleRemoved(t *testing.T) {
	m := newTestLockManager(t)
	m.opt.StaleTTL = 1 * time.Millisecond

	// manually create a stale lock file
	lockPath := filepath.Join(m.opt.LockDir, "vaultpull-stale.lock")
	_ = os.WriteFile(lockPath, nil, 0600)
	time.Sleep(5 * time.Millisecond)

	if err := m.Acquire("stale"); err != nil {
		t.Fatalf("expected stale lock to be cleared, got %v", err)
	}
	_ = m.Release("stale")
}
