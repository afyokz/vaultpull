package vault

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

var sampleSecrets = map[string]string{
	"DB_HOST": "localhost",
	"DB_PASS": "s3cr3t",
}

func TestExportSecrets_DotenvStdout(t *testing.T) {
	opts := ExportOptions{Format: FormatDotenv}
	if err := ExportSecrets(sampleSecrets, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExportSecrets_JSONToFile(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "export-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	opts := ExportOptions{Format: FormatJSON, OutFile: tmp.Name()}
	if err := ExportSecrets(sampleSecrets, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(tmp.Name())
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if m["DB_HOST"] != "localhost" {
		t.Errorf("expected localhost, got %s", m["DB_HOST"])
	}
}

func TestExportSecrets_ShellFormat(t *testing.T) {
	tmp := t.TempDir() + "/out.sh"
	opts := ExportOptions{Format: FormatShell, OutFile: tmp}
	if err := ExportSecrets(sampleSecrets, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	if !strings.Contains(string(data), "export DB_HOST=") {
		t.Error("expected shell export statement")
	}
}

func TestExportSecrets_Redacted(t *testing.T) {
	tmp := t.TempDir() + "/out.env"
	opts := ExportOptions{Format: FormatDotenv, OutFile: tmp, Redact: true}
	if err := ExportSecrets(sampleSecrets, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	if strings.Contains(string(data), "s3cr3t") {
		t.Error("expected redacted output, got plaintext")
	}
	if !strings.Contains(string(data), "***") {
		t.Error("expected *** in redacted output")
	}
}

func TestExportSecrets_InvalidFormat(t *testing.T) {
	err := ExportSecrets(sampleSecrets, ExportOptions{Format: "xml"})
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}
