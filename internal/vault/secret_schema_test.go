package vault

import (
	"strings"
	"testing"
)

func TestParseSchemaRule_Valid(t *testing.T) {
	r, err := ParseSchemaRule("API_KEY:[A-Za-z0-9]{32,}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Key != "API_KEY" {
		t.Errorf("expected key API_KEY, got %s", r.Key)
	}
	if r.Required {
		t.Error("expected not required")
	}
}

func TestParseSchemaRule_Required(t *testing.T) {
	r, err := ParseSchemaRule("TOKEN:.+:required")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Required {
		t.Error("expected required=true")
	}
}

func TestParseSchemaRule_Invalid(t *testing.T) {
	if _, err := ParseSchemaRule("NOPATTERN"); err == nil {
		t.Error("expected error for missing pattern")
	}
	if _, err := ParseSchemaRule(":pattern"); err == nil {
		t.Error("expected error for empty key")
	}
	if _, err := ParseSchemaRule("KEY:[invalid"); err == nil {
		t.Error("expected error for bad regex")
	}
}

func TestValidateSchema_NoViolations(t *testing.T) {
	rule, _ := ParseSchemaRule("DB_URL:.+")
	secrets := map[string]string{"DB_URL": "postgres://localhost/db"}
	v := ValidateSchema(secrets, []SchemaRule{rule})
	if len(v) != 0 {
		t.Errorf("expected no violations, got %v", v)
	}
}

func TestValidateSchema_PatternMismatch(t *testing.T) {
	rule, _ := ParseSchemaRule("PORT:[0-9]+")
	secrets := map[string]string{"PORT": "not-a-number"}
	v := ValidateSchema(secrets, []SchemaRule{rule})
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Key != "PORT" {
		t.Errorf("unexpected key %s", v[0].Key)
	}
}

func TestValidateSchema_RequiredMissing(t *testing.T) {
	rule, _ := ParseSchemaRule("SECRET:.+:required")
	v := ValidateSchema(map[string]string{}, []SchemaRule{rule})
	if len(v) != 1 || v[0].Key != "SECRET" {
		t.Errorf("expected missing required violation, got %v", v)
	}
}

func TestFormatSchemaReport_Clean(t *testing.T) {
	out := FormatSchemaReport(nil)
	if !strings.Contains(out, "passed") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatSchemaReport_WithViolations(t *testing.T) {
	v := []SchemaViolation{{Key: "FOO", Message: "bad value"}}
	out := FormatSchemaReport(v)
	if !strings.Contains(out, "FOO") || !strings.Contains(out, "1 violation") {
		t.Errorf("unexpected output: %s", out)
	}
}
