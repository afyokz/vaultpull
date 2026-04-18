package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func tempSnapshotPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snap.json")
}

func TestSaveAndLoadSnapshot(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	path := tempSnapshotPath(t)

	if err := SaveSnapshot(path, secrets); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	snap, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}
	if snap.Secrets["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %s", snap.Secrets["DB_HOST"])
	}
	if snap.CapturedAt.IsZero() {
		t.Error("expected non-zero CapturedAt")
	}
}

func TestLoadSnapshot_NotFound(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSaveSnapshot_FilePermissions(t *testing.T) {
	path := tempSnapshotPath(t)
	if err := SaveSnapshot(path, map[string]string{"KEY": "val"}); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %o", info.Mode().Perm())
	}
}

func TestDiffSnapshot(t *testing.T) {
	snap := &Snapshot{
		Secrets: map[string]string{
			"KEEP":    "same",
			"CHANGE":  "old",
			"REMOVED": "gone",
		},
	}
	current := map[string]string{
		"KEEP":   "same",
		"CHANGE": "new",
		"ADDED":  "fresh",
	}
	added, removed, changed := DiffSnapshot(snap, current)

	if len(added) != 1 || added[0] != "ADDED" {
		t.Errorf("added: expected [ADDED], got %v", added)
	}
	if len(removed) != 1 || removed[0] != "REMOVED" {
		t.Errorf("removed: expected [REMOVED], got %v", removed)
	}
	if len(changed) != 1 || changed[0] != "CHANGE" {
		t.Errorf("changed: expected [CHANGE], got %v", changed)
	}
}
