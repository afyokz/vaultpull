package vault

import (
	"regexp"
	"strings"
)

// RedactRule defines a pattern and replacement for redacting secret values.
type RedactRule struct {
	KeyPattern *regexp.Regexp
	Replacement string
}

// DefaultRedactReplacement is used when no replacement is specified.
const DefaultRedactReplacement = "***REDACTED***"

// ParseRedactRule compiles a key pattern into a RedactRule.
func ParseRedactRule(pattern, replacement string) (*RedactRule, error) {
	re, err := regexp.Compile("(?i)" + pattern)
	if err != nil {
		return nil, err
	}
	if replacement == "" {
		replacement = DefaultRedactReplacement
	}
	return &RedactRule{KeyPattern: re, Replacement: replacement}, nil
}

// RedactSecrets returns a copy of secrets with matching keys' values replaced.
func RedactSecrets(secrets map[string]string, rules []*RedactRule) map[string]string {
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[k] = v
	}
	if len(rules) == 0 {
		return result
	}
	for k := range result {
		for _, rule := range rules {
			if rule.KeyPattern.MatchString(k) {
				result[k] = rule.Replacement
				break
			}
		}
	}
	return result
}

// DefaultSensitivePatterns returns rules for commonly sensitive key names.
func DefaultSensitivePatterns() []*RedactRule {
	patterns := []string{"password", "secret", "token", "api_key", "private_key", "credential"}
	var rules []*RedactRule
	for _, p := range patterns {
		r, _ := ParseRedactRule(p, DefaultRedactReplacement)
		rules = append(rules, r)
	}
	return rules
}

// ContainsSensitive checks if a key name looks sensitive using default patterns.
func ContainsSensitive(key string) bool {
	sensitive := []string{"password", "secret", "token", "api_key", "private_key", "credential"}
	lower := strings.ToLower(key)
	for _, s := range sensitive {
		if strings.Contains(lower, s) {
			return true
		}
	}
	return false
}
