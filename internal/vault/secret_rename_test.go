package vault

import (
	"testing"
)

func TestParseRenameRule_Valid(t *testing.T) {
	r, err := ParseRenameRule("OLD_KEY=NEW_KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.From != "OLD_KEY" || r.To != "NEW_KEY" || !r.Exact {
		t.Errorf("unexpected rule: %+v", r)
	}
}

func TestParseRenameRule_Wildcard(t *testing.T) {
	r, err := ParseRenameRule("APP_*=SVC_*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Exact {
		t.Error("expected wildcard rule to be non-exact")
	}
}

func TestParseRenameRule_Invalid(t *testing.T) {
	cases := []string{"NOEQUALS", "=VALUE", "KEY=", ""}
	for _, c := range cases {
		_, err := ParseRenameRule(c)
		if err == nil {
			t.Errorf("expected error for input %q", c)
		}
	}
}

func TestApplyRenames_NoRules(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	out := ApplyRenames(secrets, nil)
	if out["FOO"] != "bar" {
		t.Error("expected unchanged secrets")
	}
}

func TestApplyRenames_ExactMatch(t *testing.T) {
	rule, _ := ParseRenameRule("DB_PASS=DATABASE_PASSWORD")
	out := ApplyRenames(map[string]string{"DB_PASS": "secret", "OTHER": "val"}, []RenameRule{rule})
	if _, ok := out["DB_PASS"]; ok {
		t.Error("old key should not exist")
	}
	if out["DATABASE_PASSWORD"] != "secret" {
		t.Error("new key should have value")
	}
	if out["OTHER"] != "val" {
		t.Error("unmatched key should be unchanged")
	}
}

func TestApplyRenames_WildcardMatch(t *testing.T) {
	rule, _ := ParseRenameRule("APP_*=SVC_*")
	out := ApplyRenames(map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "UNRELATED": "x"}, []RenameRule{rule})
	if out["SVC_HOST"] != "localhost" {
		t.Errorf("expected SVC_HOST, got %v", out)
	}
	if out["SVC_PORT"] != "8080" {
		t.Error("expected SVC_PORT")
	}
	if out["UNRELATED"] != "x" {
		t.Error("unrelated key should be unchanged")
	}
}
