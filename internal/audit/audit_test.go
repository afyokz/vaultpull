package audit

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLog_Success(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")

	l := NewLogger(logPath)
	keys := []string{"DB_HOST", "DB_PASS"}
	if err := l.Log("pull", "secret/myapp", keys, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	f, _ := os.Open(logPath)
	defer f.Close()

	var entry Entry
	if err := json.NewDecoder(f).Decode(&entry); err != nil {
		t.Fatalf("decode entry: %v", err)
	}

	if !entry.Success {
		t.Error("expected Success=true")
	}
	if entry.Operation != "pull" {
		t.Errorf("expected operation 'pull', got %q", entry.Operation)
	}
	if entry.Path != "secret/myapp" {
		t.Errorf("expected path 'secret/myapp', got %q", entry.Path)
	}
	if len(entry.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(entry.Keys))
	}
}

func TestLog_WithError(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")

	l := NewLogger(logPath)
	pullErr := errors.New("permission denied")
	if err := l.Log("pull", "secret/myapp", nil, pullErr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	f, _ := os.Open(logPath)
	defer f.Close()

	var entry Entry
	if err := json.NewDecoder(f).Decode(&entry); err != nil {
		t.Fatalf("decode entry: %v", err)
	}

	if entry.Success {
		t.Error("expected Success=false")
	}
	if entry.Error != "permission denied" {
		t.Errorf("unexpected error string: %q", entry.Error)
	}
}

func TestLog_AppendMultiple(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")

	l := NewLogger(logPath)
	for i := 0; i < 3; i++ {
		if err := l.Log("pull", "secret/app", []string{"KEY"}, nil); err != nil {
			t.Fatalf("log %d failed: %v", i, err)
		}
	}

	f, _ := os.Open(logPath)
	defer f.Close()

	count := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Text() != "" {
			count++
		}
	}
	if count != 3 {
		t.Errorf("expected 3 log lines, got %d", count)
	}
}
