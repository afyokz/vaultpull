package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempImport(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "import-*")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParseImportFormat_Valid(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"dotenv", "dotenv"},
		{"json", "json"},
		{"shell", "shell"},
		{"DOTENV", "dotenv"},
	} {
		got, err := ParseImportFormat(tc.in)
		if err != nil {
			t.Errorf("ParseImportFormat(%q) unexpected error: %v", tc.in, err)
		}
		if string(got) != tc.want {
			t.Errorf("got %q, want %q", got, tc.want)
		}
	}
}

func TestParseImportFormat_Invalid(t *testing.T) {
	_, err := ParseImportFormat("xml")
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestImportSecrets_Dotenv(t *testing.T) {
	path := writeTempImport(t, "# comment\nDB_HOST=localhost\nDB_PORT=\"5432\"\n\nAPI_KEY=secret\n")
	got, err := ImportSecrets(path, ImportFormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expect := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432", "API_KEY": "secret"}
	for k, v := range expect {
		if got[k] != v {
			t.Errorf("key %q: got %q, want %q", k, got[k], v)
		}
	}
}

func TestImportSecrets_JSON(t *testing.T) {
	path := writeTempImport(t, `{"HOST":"db","PORT":"3306"}`)
	got, err := ImportSecrets(path, ImportFormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["HOST"] != "db" || got["PORT"] != "3306" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestImportSecrets_Shell(t *testing.T) {
	path := writeTempImport(t, "export FOO='bar'\nexport BAZ=\"qux\"\n")
	got, err := ImportSecrets(path, ImportFormatShell)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["FOO"] != "bar" {
		t.Errorf("FOO: got %q, want %q", got["FOO"], "bar")
	}
	if got["BAZ"] != "qux" {
		t.Errorf("BAZ: got %q, want %q", got["BAZ"], "qux")
	}
}

func TestImportSecrets_FileNotFound(t *testing.T) {
	_, err := ImportSecrets(filepath.Join(t.TempDir(), "missing.env"), ImportFormatDotenv)
	if err == nil {
		t.Error("expected error for missing file")
	}
}
