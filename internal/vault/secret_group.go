package vault

import (
	"fmt"
	"strings"
)

// GroupRule maps a key pattern to a named group.
type GroupRule struct {
	Pattern string
	Group   string
}

// ParseGroupRule parses a rule in the form "pattern=group".
func ParseGroupRule(raw string) (GroupRule, error) {
	parts := strings.SplitN(raw, "=", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return GroupRule{}, fmt.Errorf("invalid group rule %q: expected pattern=group", raw)
	}
	return GroupRule{Pattern: parts[0], Group: parts[1]}, nil
}

// ParseGroupRules parses multiple group rules.
func ParseGroupRules(raws []string) ([]GroupRule, error) {
	rules := make([]GroupRule, 0, len(raws))
	for _, r := range raws {
		gr, err := ParseGroupRule(r)
		if err != nil {
			return nil, err
		}
		rules = append(rules, gr)
	}
	return rules, nil
}

// GroupSecrets partitions secrets into named groups based on rules.
// Keys not matching any rule are placed in the "default" group.
func GroupSecrets(secrets map[string]string, rules []GroupRule) map[string]map[string]string {
	result := map[string]map[string]string{}
	ensure := func(g string) {
		if _, ok := result[g]; !ok {
			result[g] = map[string]string{}
		}
	}
	for k, v := range secrets {
		matched := false
		for _, rule := range rules {
			if matchPattern(rule.Pattern, k) {
				ensure(rule.Group)
				result[rule.Group][k] = v
				matched = true
				break
			}
		}
		if !matched {
			ensure("default")
			result["default"][k] = v
		}
	}
	return result
}
