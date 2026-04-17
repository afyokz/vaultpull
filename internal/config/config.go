package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the vaultpull configuration file structure.
type Config struct {
	Vault   VaultConfig   `yaml:"vault"`
	Secrets []SecretEntry `yaml:"secrets"`
}

// VaultConfig holds Vault connection settings.
type VaultConfig struct {
	Address   string `yaml:"address"`
	TokenEnv  string `yaml:"token_env"`
	MountPath string `yaml:"mount_path"`
}

// SecretEntry maps a Vault secret path to a local .env file.
type SecretEntry struct {
	Path    string `yaml:"path"`
	EnvFile string `yaml:"env_file"`
}

const defaultTokenEnv = "VAULT_TOKEN"
const defaultMount = "secret"

// Load reads and parses a vaultpull config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	if cfg.Vault.TokenEnv == "" {
		cfg.Vault.TokenEnv = defaultTokenEnv
	}
	if cfg.Vault.MountPath == "" {
		cfg.Vault.MountPath = defaultMount
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Vault.Address == "" {
		return fmt.Errorf("vault.address is required")
	}
	if len(c.Secrets) == 0 {
		return fmt.Errorf("at least one entry under secrets is required")
	}
	for i, s := range c.Secrets {
		if s.Path == "" {
			return fmt.Errorf("secrets[%d].path is required", i)
		}
		if s.EnvFile == "" {
			return fmt.Errorf("secrets[%d].env_file is required", i)
		}
	}
	return nil
}
