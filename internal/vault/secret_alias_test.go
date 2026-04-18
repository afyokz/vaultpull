package vault

import (
	"testing"
)

func TestParseAliasRule_Valid(t *testing.T) {
	orig, alias, err := ParseAliasRule("DB_PASS=DATABASE_PASSWORD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if orig != "DB_PASS" || alias != "DATABASE_PASSWORD" {
		t.Fatalf("got orig=%q alias=%q", orig, alias)
	}
}

func TestParseAliasRule_Invalid(t *testing.T) {
	cases := []string{"NODASH", "=ALIAS", "ORIG=", ""}
	for _, c := range cases {
		_, _, err := ParseAliasRule(c)
		if err == nil {
			t.Errorf("expected error for rule %q", c)
		}
	}
}

func TestParseAliasRules_Valid(t *testing.T) {
	am, err := ParseAliasRules([]string{"A=B", "C=D"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if am["A"] != "B" || am["C"] != "D" {
		t.Fatalf("unexpected map: %v", am)
	}
}

func TestApplyAliases_NoRules(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	out := ApplyAliases(secrets, AliasMap{})
	if out["KEY"] != "val" {
		t.Fatal("expected key to be preserved")
	}
}

func TestApplyAliases_Renames(t *testing.T) {
	secrets := map[string]string{"OLD_KEY": "secret", "KEEP": "yes"}
	am := AliasMap{"OLD_KEY": "NEW_KEY"}
	out := ApplyAliases(secrets, am)
	if _, ok := out["OLD_KEY"]; ok {
		t.Fatal("original key should be removed")
	}
	if out["NEW_KEY"] != "secret" {
		t.Fatal("alias key should hold original value")
	}
	if out["KEEP"] != "yes" {
		t.Fatal("unaliased key should be preserved")
	}
}

func TestApplyAliases_NoConflict(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	am := AliasMap{"A": "C"}
	out := ApplyAliases(secrets, am)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}
