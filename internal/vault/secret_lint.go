package vault

import (
	"fmt"
	"strings"
)

// LintSeverity indicates how serious a lint issue is.
type LintSeverity string

const (
	SeverityWarn  LintSeverity = "WARN"
	SeverityError LintSeverity = "ERROR"
)

// LintIssue describes a single problem found in a secret map.
type LintIssue struct {
	Key      string
	Message  string
	Severity LintSeverity
}

func (i LintIssue) String() string {
	return fmt.Sprintf("[%s] %s: %s", i.Severity, i.Key, i.Message)
}

// LintSecrets checks secrets for common problems and returns any issues found.
func LintSecrets(secrets map[string]string) []LintIssue {
	var issues []LintIssue

	for k, v := range secrets {
		if k == "" {
			issues = append(issues, LintIssue{Key: "(empty)", Message: "key is empty", Severity: SeverityError})
			continue
		}
		if strings.ToUpper(k) != k {
			issues = append(issues, LintIssue{Key: k, Message: "key is not uppercase", Severity: SeverityWarn})
		}
		if strings.ContainsAny(k, " \t") {
			issues = append(issues, LintIssue{Key: k, Message: "key contains whitespace", Severity: SeverityError})
		}
		if v == "" {
			issues = append(issues, LintIssue{Key: k, Message: "value is empty", Severity: SeverityWarn})
		}
		if strings.HasPrefix(v, " ") || strings.HasSuffix(v, " ") {
			issues = append(issues, LintIssue{Key: k, Message: "value has leading or trailing whitespace", Severity: SeverityWarn})
		}
	}

	return issues
}

// FormatLintReport returns a human-readable lint report string.
func FormatLintReport(issues []LintIssue) string {
	if len(issues) == 0 {
		return "lint: no issues found"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("lint: %d issue(s) found\n", len(issues)))
	for _, issue := range issues {
		sb.WriteString("  " + issue.String() + "\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}
