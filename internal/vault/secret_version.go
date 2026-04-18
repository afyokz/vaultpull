package vault

import (
	"fmt"
	"strconv"
	"strings"
)

// VersionedPath represents a secret path with an optional version specifier.
type VersionedPath struct {
	SecretPath
	Version int // 0 means latest
}

// ParseVersionedPath parses a path like "secret/data/myapp@3" into a VersionedPath.
func ParseVersionedPath(raw string) (VersionedPath, error) {
	var vp VersionedPath

	parts := strings.SplitN(raw, "@", 2)
	base := parts[0]

	sp, err := ParseSecretPath(base)
	if err != nil {
		return vp, fmt.Errorf("invalid secret path: %w", err)
	}
	vp.SecretPath = sp

	if len(parts) == 2 {
		v, err := strconv.Atoi(parts[1])
		if err != nil || v < 1 {
			return vp, fmt.Errorf("invalid version %q: must be a positive integer", parts[1])
		}
		vp.Version = v
	}

	return vp, nil
}

// String returns the canonical string representation, e.g. "secret/data/myapp@3".
func (vp VersionedPath) String() string {
	base := vp.SecretPath.String()
	if vp.Version > 0 {
		return fmt.Sprintf("%s@%d", base, vp.Version)
	}
	return base
}

// VersionParams returns the query parameters map to pass to the Vault API.
// An empty map means "fetch latest".
func (vp VersionedPath) VersionParams() map[string]string {
	if vp.Version > 0 {
		return map[string]string{"version": strconv.Itoa(vp.Version)}
	}
	return map[string]string{}
}
