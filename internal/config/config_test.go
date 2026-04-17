package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "vaultpull-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	yaml := `
vault:
  address: https://vault.example.com
secrets:
  - path: myapp/prod
    env_file: .env
`
	path := writeTemp(t, yaml)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Address != "https://vault.example.com" {
		t.Errorf("expected address, got %q", cfg.Vault.Address)
	}
	if cfg.Vault.TokenEnv != defaultTokenEnv {
		t.Errorf("expected default token env, got %q", cfg.Vault.TokenEnv)
	}
	if cfg.Vault.MountPath != defaultMount {
		t.Errorf("expected default mount, got %q", cfg.Vault.MountPath)
	}
	if len(cfg.Secrets) != 1 || cfg.Secrets[0].EnvFile != ".env" {
		t.Errorf("unexpected secrets: %+v", cfg.Secrets)
	}
}

func TestLoad_MissingAddress(t *testing.T) {
	yaml := `
secrets:
  - path: myapp/prod
    env_file: .env
`
	path := writeTemp(t, yaml)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing vault address")
	}
}

func TestLoad_NoSecrets(t *testing.T) {
	yaml := `
vault:
  address: https://vault.example.com
`
	path := writeTemp(t, yaml)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
