package vault

import (
	"fmt"
	"strings"
)

// SecretPath represents a parsed Vault KV secret path.
type SecretPath struct {
	Mount   string
	SubPath string
	Version int // 0 means latest
}

// ParseSecretPath parses a path like "secret/data/myapp" or "secret/myapp".
// It normalises KV v2 paths by injecting "data/" if missing.
func ParseSecretPath(raw string) (*SecretPath, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("secret path must not be empty")
	}

	parts := strings.SplitN(raw, "/", 2)
	if len(parts) < 2 || parts[1] == "" {
		return nil, fmt.Errorf("secret path %q must contain a mount and a sub-path (e.g. secret/myapp)", raw)
	}

	mount := parts[0]
	sub := parts[1]

	// Inject "data/" segment for KV v2 if not already present.
	if !strings.HasPrefix(sub, "data/") && !strings.HasPrefix(sub, "metadata/") {
		sub = "data/" + sub
	}

	return &SecretPath{
		Mount:   mount,
		SubPath: sub,
	}, nil
}

// FullPath returns the full Vault API path.
func (s *SecretPath) FullPath() string {
	return s.Mount + "/" + s.SubPath
}

// String implements Stringer.
func (s *SecretPath) String() string {
	return s.FullPath()
}
