package config

import (
	"testing"
)

func TestBuildRedactRules_Nil(t *testing.T) {
	rules, err := BuildRedactRules(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected no rules, got %d", len(rules))
	}
}

func TestBuildRedactRules_Disabled(t *testing.T) {
	cfg := &RedactConfig{Enabled: false, UseDefaults: true}
	rules, err := BuildRedactRules(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected no rules when disabled")
	}
}

func TestBuildRedactRules_UseDefaults(t *testing.T) {
	cfg := &RedactConfig{Enabled: true, UseDefaults: true}
	rules, err := BuildRedactRules(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) == 0 {
		t.Error("expected default rules")
	}
}

func TestBuildRedactRules_CustomPattern(t *testing.T) {
	cfg := &RedactConfig{
		Enabled:     true,
		Patterns:    []string{"my_custom_key"},
		Replacement: "[gone]",
	}
	rules, err := BuildRedactRules(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(rules))
	}
}

func TestBuildRedactRules_InvalidPattern(t *testing.T) {
	cfg := &RedactConfig{
		Enabled:  true,
		Patterns: []string{"[bad"},
	}
	_, err := BuildRedactRules(cfg)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}
