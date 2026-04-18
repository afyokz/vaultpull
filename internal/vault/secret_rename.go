package vault

import (
	"fmt"
	"strings"
)

// RenameRule maps a source key prefix or exact key to a target name.
type RenameRule struct {
	From string
	To   string
	Exact bool
}

// ParseRenameRule parses a rename rule string in the format "FROM=TO".
func ParseRenameRule(s string) (RenameRule, error) {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return RenameRule{}, fmt.Errorf("invalid rename rule %q: expected FROM=TO", s)
	}
	return RenameRule{
		From:  parts[0],
		To:    parts[1],
		Exact: !strings.HasSuffix(parts[0], "*"),
	}, nil
}

// ApplyRenames applies a slice of RenameRules to a secrets map, returning a new map.
func ApplyRenames(secrets map[string]string, rules []RenameRule) map[string]string {
	if len(rules) == 0 {
		return secrets
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[applyRenameKey(k, rules)] = v
	}
	return out
}

func applyRenameKey(key string, rules []RenameRule) string {
	for _, r := range rules {
		if r.Exact {
			if key == r.From {
				return r.To
			}
		} else {
			prefix := strings.TrimSuffix(r.From, "*")
			if strings.HasPrefix(key, prefix) {
				return r.To + strings.TrimPrefix(key, prefix)
			}
		}
	}
	return key
}
