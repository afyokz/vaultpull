package vault

import (
	"strings"
	"testing"
)

func TestValidateSecrets_NoRules(t *testing.T) {
	issues := ValidateSecrets(map[string]string{"KEY": "val"}, nil)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestValidateSecrets_RequiredMissing(t *testing.T) {
	rules := []ValidationRule{{Key: "DB_PASS", Required: true}}
	issues := ValidateSecrets(map[string]string{}, rules)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "DB_PASS" {
		t.Errorf("unexpected key %s", issues[0].Key)
	}
}

func TestValidateSecrets_RequiredPresent(t *testing.T) {
	rules := []ValidationRule{{Key: "DB_PASS", Required: true}}
	issues := ValidateSecrets(map[string]string{"DB_PASS": "secret"}, rules)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestValidateSecrets_MinLen(t *testing.T) {
	rules := []ValidationRule{{Key: "TOKEN", MinLen: 10}}
	issues := ValidateSecrets(map[string]string{"TOKEN": "short"}, rules)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if !strings.Contains(issues[0].Message, "too short") {
		t.Errorf("unexpected message: %s", issues[0].Message)
	}
}

func TestValidateSecrets_Forbidden(t *testing.T) {
	rules := []ValidationRule{{Key: "API_KEY", Forbidden: []string{"changeme", "default"}}}
	issues := ValidateSecrets(map[string]string{"API_KEY": "changeme123"}, rules)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestFormatValidationReport_NoIssues(t *testing.T) {
	out := FormatValidationReport(nil)
	if !strings.Contains(out, "passed") {
		t.Errorf("expected passed message, got: %s", out)
	}
}

func TestFormatValidationReport_WithIssues(t *testing.T) {
	issues := []ValidationIssue{{Key: "FOO", Message: "required but missing or empty"}}
	out := FormatValidationReport(issues)
	if !strings.Contains(out, "failed") {
		t.Errorf("expected failed in output, got: %s", out)
	}
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected key FOO in output")
	}
}
