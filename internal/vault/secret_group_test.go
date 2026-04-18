package vault

import (
	"testing"
)

func TestParseGroupRule_Valid(t *testing.T) {
	gr, err := ParseGroupRule("DB_*=database")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gr.Pattern != "DB_*" || gr.Group != "database" {
		t.Errorf("unexpected rule: %+v", gr)
	}
}

func TestParseGroupRule_Invalid(t *testing.T) {
	for _, raw := range []string{"", "noequals", "=group", "pattern="} {
		_, err := ParseGroupRule(raw)
		if err == nil {
			t.Errorf("expected error for %q", raw)
		}
	}
}

func TestGroupSecrets_NoRules(t *testing.T) {
	secrets := map[string]string{"FOO": "1", "BAR": "2"}
	result := GroupSecrets(secrets, nil)
	if len(result["default"]) != 2 {
		t.Errorf("expected 2 in default group, got %d", len(result["default"]))
	}
}

func TestGroupSecrets_MatchesRule(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "APP_KEY": "secret", "DB_PORT": "5432"}
	rules := []GroupRule{{Pattern: "DB_*", Group: "database"}}
	result := GroupSecrets(secrets, rules)
	if len(result["database"]) != 2 {
		t.Errorf("expected 2 in database group, got %d", len(result["database"]))
	}
	if len(result["default"]) != 1 {
		t.Errorf("expected 1 in default group, got %d", len(result["default"]))
	}
}

func TestGroupSecrets_MultipleGroups(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "h", "REDIS_URL": "r", "OTHER": "o"}
	rules := []GroupRule{
		{Pattern: "DB_*", Group: "database"},
		{Pattern: "REDIS_*", Group: "cache"},
	}
	result := GroupSecrets(secrets, rules)
	if result["database"]["DB_HOST"] != "h" {
		t.Error("DB_HOST not in database group")
	}
	if result["cache"]["REDIS_URL"] != "r" {
		t.Error("REDIS_URL not in cache group")
	}
	if result["default"]["OTHER"] != "o" {
		t.Error("OTHER not in default group")
	}
}
