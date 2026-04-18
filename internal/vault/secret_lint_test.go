package vault

import (
	"strings"
	"testing"
)

func TestLintSecrets_Clean(t *testing.T) {
	secrets := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"API_KEY":      "abc123",
	}
	issues := LintSecrets(secrets)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestLintSecrets_LowercaseKey(t *testing.T) {
	secrets := map[string]string{"db_host": "localhost"}
	issues := LintSecrets(secrets)
	if len(issues) != 1 || issues[0].Severity != SeverityWarn {
		t.Fatalf("expected 1 warn for lowercase key, got %v", issues)
	}
}

func TestLintSecrets_EmptyValue(t *testing.T) {
	secrets := map[string]string{"TOKEN": ""}
	issues := LintSecrets(secrets)
	if len(issues) != 1 || issues[0].Severity != SeverityWarn {
		t.Fatalf("expected 1 warn for empty value, got %v", issues)
	}
}

func TestLintSecrets_WhitespaceInKey(t *testing.T) {
	secrets := map[string]string{"BAD KEY": "value"}
	issues := LintSecrets(secrets)
	found := false
	for _, i := range issues {
		if i.Severity == SeverityError && strings.Contains(i.Message, "whitespace") {
			found = true
		}
	}
	if !found {
		t.Fatal("expected error for whitespace in key")
	}
}

func TestLintSecrets_LeadingTrailingSpace(t *testing.T) {
	secrets := map[string]string{"SECRET": " value "}
	issues := LintSecrets(secrets)
	if len(issues) != 1 || !strings.Contains(issues[0].Message, "whitespace") {
		t.Fatalf("expected whitespace warn, got %v", issues)
	}
}

func TestFormatLintReport_NoIssues(t *testing.T) {
	out := FormatLintReport(nil)
	if out != "lint: no issues found" {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestFormatLintReport_WithIssues(t *testing.T) {
	issues := []LintIssue{
		{Key: "foo", Message: "key is not uppercase", Severity: SeverityWarn},
	}
	out := FormatLintReport(issues)
	if !strings.Contains(out, "1 issue(s)") || !strings.Contains(out, "foo") {
		t.Fatalf("unexpected report: %s", out)
	}
}
