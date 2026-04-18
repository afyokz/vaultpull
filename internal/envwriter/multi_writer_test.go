package envwriter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGroupedWrite_CreatesFiles(t *testing.T) {
	dir := t.TempDir()
	groups := map[string]map[string]string{
		"database": {"DB_HOST": "localhost"},
		"cache":    {"REDIS_URL": "redis://localhost"},
	}
	if err := GroupedWrite(dir, groups, 0600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, name := range []string{".env.database", ".env.cache"} {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s to exist", name)
		}
	}
}

func TestGroupedWrite_ContainsKeys(t *testing.T) {
	dir := t.TempDir()
	groups := map[string]map[string]string{
		"default": {"FOO": "bar"},
	}
	if err := GroupedWrite(dir, groups, 0600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, ".env.default"))
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if !strings.Contains(string(data), "FOO=bar") {
		t.Errorf("expected FOO=bar in file, got: %s", data)
	}
}

func TestGroupedWrite_EmptyGroups(t *testing.T) {
	dir := t.TempDir()
	if err := GroupedWrite(dir, map[string]map[string]string{}, 0600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Errorf("expected no files, got %d", len(entries))
	}
}
