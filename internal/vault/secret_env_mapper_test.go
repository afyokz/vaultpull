package vault

import (
	"testing"
)

func TestParseMappingRule_Valid(t *testing.T) {
	rule, err := ParseMappingRule("secret/app:APP")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.Path != "secret/app" || rule.Prefix != "APP" {
		t.Errorf("got %+v", rule)
	}
}

func TestParseMappingRule_Invalid(t *testing.T) {
	cases := []string{"", "nocolon", ":PREFIX", "path:"}
	for _, c := range cases {
		_, err := ParseMappingRule(c)
		if err == nil {
			t.Errorf("expected error for %q", c)
		}
	}
}

func TestApplyMappings_NoRules(t *testing.T) {
	byPath := map[string]map[string]string{
		"secret/app": {"KEY": "val"},
	}
	out := ApplyMappings(byPath, nil)
	if out["KEY"] != "val" {
		t.Errorf("expected KEY=val, got %v", out)
	}
}

func TestApplyMappings_AppliesPrefix(t *testing.T) {
	byPath := map[string]map[string]string{
		"secret/app": {"TOKEN": "abc"},
	}
	rules := []MappingRule{{Path: "secret/app", Prefix: "APP"}}
	out := ApplyMappings(byPath, rules)
	if out["APP_TOKEN"] != "abc" {
		t.Errorf("expected APP_TOKEN=abc, got %v", out)
	}
	if _, ok := out["TOKEN"]; ok {
		t.Error("unprefixed key should not exist")
	}
}

func TestApplyMappings_UnmatchedPathPassthrough(t *testing.T) {
	byPath := map[string]map[string]string{
		"secret/app":   {"TOKEN": "abc"},
		"secret/other": {"FOO": "bar"},
	}
	rules := []MappingRule{{Path: "secret/app", Prefix: "APP"}}
	out := ApplyMappings(byPath, rules)
	if out["APP_TOKEN"] != "abc" {
		t.Errorf("expected APP_TOKEN=abc")
	}
	if out["FOO"] != "bar" {
		t.Errorf("expected FOO=bar for unmatched path")
	}
}
