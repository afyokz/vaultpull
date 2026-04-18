package backup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreate_NoSourceFile(t *testing.T) {
	tmp := t.TempDir()
	m := NewManager(3)
	path, err := m.Create(filepath.Join(tmp, ".env"))
	if err != nil {
		t.Fatalf("expected no error for missing source, got %v", err)
	}
	if path != "" {
		t.Errorf("expected empty path, got %q", path)
	}
}

func TestCreate_CreatesBackup(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, ".env")
	if err := os.WriteFile(src, []byte("KEY=value\n"), 0600); err != nil {
		t.Fatal(err)
	}
	m := NewManager(3)
	backupPath, err := m.Create(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if backupPath == "" {
		t.Fatal("expected a backup path")
	}
	data, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("cannot read backup: %v", err)
	}
	if string(data) != "KEY=value\n" {
		t.Errorf("backup content mismatch: %q", string(data))
	}
}

func TestCreate_PrunesOldBackups(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, ".env")
	m := NewManager(2)

	for i := 0; i < 4; i++ {
		if err := os.WriteFile(src, []byte("X=1"), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := m.Create(src); err != nil {
			t.Fatalf("create backup %d: %v", i, err)
		}
	}

	matches, _ := filepath.Glob(filepath.Join(tmp, ".env.*.bak"))
	if len(matches) > 2 {
		t.Errorf("expected at most 2 backups, got %d", len(matches))
	}
}

func TestCreate_BackupPermissions(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, ".env")
	if err := os.WriteFile(src, []byte("A=1"), 0600); err != nil {
		t.Fatal(err)
	}
	m := NewManager(3)
	backupPath, err := m.Create(src)
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(backupPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}
