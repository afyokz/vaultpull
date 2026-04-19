package vault

import (
	"testing"
)

func TestParseRedactRule_Valid(t *testing.T) {
	r, err := ParseRedactRule("password", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Replacement != DefaultRedactReplacement {
		t.Errorf("expected default replacement, got %q", r.Replacement)
	}
}

func TestParseRedactRule_Invalid(t *testing.T) {
	_, err := ParseRedactRule("[invalid", "")
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestRedactSecrets_NoRules(t *testing.T) {
	secrets := map[string]string{"DB_PASSWORD": "s3cr3t"}
	out := RedactSecrets(secrets, nil)
	if out["DB_PASSWORD"] != "s3cr3t" {
		t.Errorf("expected value unchanged, got %q", out["DB_PASSWORD"])
	}
}

func TestRedactSecrets_MatchesKey(t *testing.T) {
	r, _ := ParseRedactRule("password", "")
	secrets := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"DB_HOST":     "localhost",
	}
	out := RedactSecrets(secrets, []*RedactRule{r})
	if out["DB_PASSWORD"] != DefaultRedactReplacement {
		t.Errorf("expected redacted, got %q", out["DB_PASSWORD"])
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected unchanged, got %q", out["DB_HOST"])
	}
}

func TestRedactSecrets_CustomReplacement(t *testing.T) {
	r, _ := ParseRedactRule("token", "[hidden]")
	secrets := map[string]string{"API_TOKEN": "abc123"}
	out := RedactSecrets(secrets, []*RedactRule{r})
	if out["API_TOKEN"] != "[hidden]" {
		t.Errorf("expected [hidden], got %q", out["API_TOKEN"])
	}
}

func TestDefaultSensitivePatterns(t *testing.T) {
	rules := DefaultSensitivePatterns()
	if len(rules) == 0 {
		t.Fatal("expected non-empty default patterns")
	}
	secrets := map[string]string{
		"APP_SECRET": "topsecret",
		"DB_HOST":    "localhost",
	}
	out := RedactSecrets(secrets, rules)
	if out["APP_SECRET"] != DefaultRedactReplacement {
		t.Errorf("expected redacted APP_SECRET")
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST unchanged")
	}
}

func TestContainsSensitive(t *testing.T) {
	if !ContainsSensitive("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD to be sensitive")
	}
	if ContainsSensitive("DB_HOST") {
		t.Error("expected DB_HOST to not be sensitive")
	}
}
