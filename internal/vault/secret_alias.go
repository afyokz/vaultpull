package vault

import (
	"fmt"
	"strings"
)

// AliasMap holds a mapping from original secret keys to alias names.
type AliasMap map[string]string

// ParseAliasRule parses a single alias rule in the form "original=alias".
func ParseAliasRule(rule string) (original, alias string, err error) {
	parts := strings.SplitN(rule, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid alias rule %q: expected format ORIGINAL=ALIAS", rule)
	}
	original = strings.TrimSpace(parts[0])
	alias = strings.TrimSpace(parts[1])
	if original == "" {
		return "", "", fmt.Errorf("invalid alias rule %q: original key must not be empty", rule)
	}
	if alias == "" {
		return "", "", fmt.Errorf("invalid alias rule %q: alias must not be empty", rule)
	}
	return original, alias, nil
}

// ParseAliasRules parses a slice of alias rule strings into an AliasMap.
func ParseAliasRules(rules []string) (AliasMap, error) {
	am := make(AliasMap, len(rules))
	for _, r := range rules {
		orig, alias, err := ParseAliasRule(r)
		if err != nil {
			return nil, err
		}
		am[orig] = alias
	}
	return am, nil
}

// ApplyAliases returns a new map where keys that appear in the AliasMap are
// renamed to their alias. Keys not present in the map are kept as-is.
func ApplyAliases(secrets map[string]string, am AliasMap) map[string]string {
	if len(am) == 0 {
		return secrets
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if alias, ok := am[k]; ok {
			out[alias] = v
		} else {
			out[k] = v
		}
	}
	return out
}
