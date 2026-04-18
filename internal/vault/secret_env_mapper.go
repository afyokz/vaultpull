package vault

import (
	"fmt"
	"strings"
)

// MappingRule defines how a vault path maps to an env var prefix.
type MappingRule struct {
	Path   string
	Prefix string
}

// ParseMappingRule parses a rule in the format "vault/path:ENV_PREFIX".
func ParseMappingRule(raw string) (MappingRule, error) {
	parts := strings.SplitN(raw, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return MappingRule{}, fmt.Errorf("invalid mapping rule %q: expected path:PREFIX", raw)
	}
	return MappingRule{Path: parts[0], Prefix: parts[1]}, nil
}

// ParseMappingRules parses multiple raw mapping rules.
func ParseMappingRules(raws []string) ([]MappingRule, error) {
	rules := make([]MappingRule, 0, len(raws))
	for _, r := range raws {
		rule, err := ParseMappingRule(r)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

// ApplyMappings applies prefix mappings to secrets keyed by vault path.
// Secrets not matching any rule are included unchanged.
func ApplyMappings(byPath map[string]map[string]string, rules []MappingRule) map[string]string {
	result := make(map[string]string)
	matched := make(map[string]bool)

	for _, rule := range rules {
		if secrets, ok := byPath[rule.Path]; ok {
			for k, v := range secrets {
				result[rule.Prefix+"_"+k] = v
			}
			matched[rule.Path] = true
		}
	}

	for path, secrets := range byPath {
		if matched[path] {
			continue
		}
		for k, v := range secrets {
			result[k] = v
		}
	}

	return result
}
