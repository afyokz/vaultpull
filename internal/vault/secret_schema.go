package vault

import (
	"fmt"
	"regexp"
	"strings"
)

// SchemaRule defines an expected shape for a secret key.
type SchemaRule struct {
	Key     string
	Pattern *regexp.Regexp
	Required bool
}

// SchemaViolation describes a single schema mismatch.
type SchemaViolation struct {
	Key     string
	Message string
}

// ParseSchemaRule parses a rule string of the form "KEY:PATTERN" or "KEY:PATTERN:required".
func ParseSchemaRule(raw string) (SchemaRule, error) {
	parts := strings.SplitN(raw, ":", 3)
	if len(parts) < 2 {
		return SchemaRule{}, fmt.Errorf("invalid schema rule %q: expected KEY:PATTERN", raw)
	}
	key := strings.TrimSpace(parts[0])
	if key == "" {
		return SchemaRule{}, fmt.Errorf("schema rule key must not be empty")
	}
	re, err := regexp.Compile(parts[1])
	if err != nil {
		return SchemaRule{}, fmt.Errorf("invalid pattern in schema rule %q: %w", raw, err)
	}
	required := len(parts) == 3 && strings.TrimSpace(parts[2]) == "required"
	return SchemaRule{Key: key, Pattern: re, Required: required}, nil
}

// ParseSchemaRules parses multiple rule strings and returns all rules or all errors encountered.
func ParseSchemaRules(raws []string) ([]SchemaRule, error) {
	var rules []SchemaRule
	var errs []string
	for _, raw := range raws {
		rule, err := ParseSchemaRule(raw)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		rules = append(rules, rule)
	}
	if len(errs) > 0 {
		return nil, fmt.Errorf("schema parse errors:\n  %s", strings.Join(errs, "\n  "))
	}
	return rules, nil
}

// ValidateSchema checks secrets against schema rules and returns violations.
func ValidateSchema(secrets map[string]string, rules []SchemaRule) []SchemaViolation {
	var violations []SchemaViolation
	for _, rule := range rules {
		val, exists := secrets[rule.Key]
		if !exists {
			if rule.Required {
				violations = append(violations, SchemaViolation{Key: rule.Key, Message: "required key missing"})
			}
			continue
		}
		if !rule.Pattern.MatchString(val) {
			violations = append(violations, SchemaViolation{
				Key:     rule.Key,
				Message: fmt.Sprintf("value does not match pattern %q", rule.Pattern.String()),
			})
		}
	}
	return violations
}

// FormatSchemaReport returns a human-readable schema validation report.
func FormatSchemaReport(violations []SchemaViolation) string {
	if len(violations) == 0 {
		return "schema: all checks passed"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "schema: %d violation(s)\n", len(violations))
	for _, v := range violations {
		fmt.Fprintf(&sb, "  [%s] %s\n", v.Key, v.Message)
	}
	return strings.TrimRight(sb.String(), "\n")
}
