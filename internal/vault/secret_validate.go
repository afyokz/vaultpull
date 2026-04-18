package vault

import (
	"fmt"
	"strings"
)

// ValidationRule defines a rule for validating secret values.
type ValidationRule struct {
	Key      string
	Required bool
	MinLen   int
	Forbidden []string
}

// ValidationIssue represents a single validation problem.
type ValidationIssue struct {
	Key     string
	Message string
}

func (v ValidationIssue) String() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// ValidateSecrets checks secrets against the provided rules and returns any issues.
func ValidateSecrets(secrets map[string]string, rules []ValidationRule) []ValidationIssue {
	var issues []ValidationIssue

	for _, rule := range rules {
		val, exists := secrets[rule.Key]

		if rule.Required && (!exists || val == "") {
			issues = append(issues, ValidationIssue{Key: rule.Key, Message: "required but missing or empty"})
			continue
		}

		if !exists {
			continue
		}

		if rule.MinLen > 0 && len(val) < rule.MinLen {
			issues = append(issues, ValidationIssue{
				Key:     rule.Key,
				Message: fmt.Sprintf("value too short (min %d chars)", rule.MinLen),
			})
		}

		for _, forbidden := range rule.Forbidden {
			if strings.Contains(val, forbidden) {
				issues = append(issues, ValidationIssue{
					Key:     rule.Key,
					Message: fmt.Sprintf("contains forbidden string %q", forbidden),
				})
			}
		}
	}

	return issues
}

// FormatValidationReport returns a human-readable report of validation issues.
func FormatValidationReport(issues []ValidationIssue) string {
	if len(issues) == 0 {
		return "validation passed: no issues found"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "validation failed: %d issue(s)\n", len(issues))
	for _, issue := range issues {
		fmt.Fprintf(&sb, "  - %s\n", issue)
	}
	return strings.TrimRight(sb.String(), "\n")
}
