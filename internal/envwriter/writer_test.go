package envwriter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWrite_NewFile(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), ".env")
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}

	if err := Write(tmp, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := readEnvFile(tmp)
	if err != nil {
		t.Fatalf("failed to read back env file: %v", err)
	}
	for k, v := range secrets {
		if result[k] != v {
			t.Errorf("expected %s=%s, got %s", k, v, result[k])
		}
	}
}

func TestWrite_MergesExisting(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), ".env")
	_ = os.WriteFile(tmp, []byte("EXISTING_KEY=old_value\nOTHER=keep\n"), 0600)

	secrets := map[string]string{"EXISTING_KEY": "new_value", "NEW_KEY": "added"}
	if err := Write(tmp, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := readEnvFile(tmp)
	if err != nil {
		t.Fatalf("failed to read back env file: %v", err)
	}
	if result["EXISTING_KEY"] != "new_value" {
		t.Errorf("expected updated value, got %s", result["EXISTING_KEY"])
	}
	if result["OTHER"] != "keep" {
		t.Errorf("expected OTHER to be preserved, got %s", result["OTHER"])
	}
	if result["NEW_KEY"] != "added" {
		t.Errorf("expected NEW_KEY=added, got %s", result["NEW_KEY"])
	}
}

func TestWrite_FilePermissions(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), ".env")
	_ = Write(tmp, map[string]string{"KEY": "val"})

	info, err := os.Stat(tmp)
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected file mode 0600, got %v", info.Mode().Perm())
	}
}

func TestWrite_EmptySecrets(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), ".env")
	_ = os.WriteFile(tmp, []byte("KEEP=value\n"), 0600)

	if err := Write(tmp, map[string]string{}); err != nil {
		t.Fatalf("unexpected error writing empty secrets: %v", err)
	}

	result, err := readEnvFile(tmp)
	if err != nil {
		t.Fatalf("failed to read back env file: %v", err)
	}
	if result["KEEP"] != "value" {
		t.Errorf("expected KEEP=value to be preserved, got %s", result["KEEP"])
	}
}
