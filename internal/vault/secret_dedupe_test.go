package vault

import (
	"testing"
)

func TestParseDedupeStrategy_Valid(t *testing.T) {
	for _, s := range []string{"keep-first", "keep-last", "error"} {
		_, err := ParseDedupeStrategy(s)
		if err != nil {
			t.Errorf("expected %q to be valid, got %v", s, err)
		}
	}
}

func TestParseDedupeStrategy_Invalid(t *testing.T) {
	_, err := ParseDedupeStrategy("unknown")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestDedupeSecrets_NoDuplicates(t *testing.T) {
	sources := []map[string]string{
		{"A": "1"},
		{"B": "2"},
	}
	res, err := DedupeSecrets(sources, DedupeKeepFirst)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Duplicates) != 0 {
		t.Errorf("expected no duplicates, got %v", res.Duplicates)
	}
	if res.Secrets["A"] != "1" || res.Secrets["B"] != "2" {
		t.Errorf("unexpected secrets: %v", res.Secrets)
	}
}

func TestDedupeSecrets_KeepFirst(t *testing.T) {
	sources := []map[string]string{
		{"KEY": "first"},
		{"KEY": "second"},
	}
	res, err := DedupeSecrets(sources, DedupeKeepFirst)
	if err != nil {
		t.Fatal(err)
	}
	if res.Secrets["KEY"] != "first" {
		t.Errorf("expected 'first', got %q", res.Secrets["KEY"])
	}
	if len(res.Duplicates) != 1 {
		t.Errorf("expected 1 duplicate, got %d", len(res.Duplicates))
	}
}

func TestDedupeSecrets_KeepLast(t *testing.T) {
	sources := []map[string]string{
		{"KEY": "first"},
		{"KEY": "second"},
	}
	res, err := DedupeSecrets(sources, DedupeKeepLast)
	if err != nil {
		t.Fatal(err)
	}
	if res.Secrets["KEY"] != "second" {
		t.Errorf("expected 'second', got %q", res.Secrets["KEY"])
	}
}

func TestDedupeSecrets_ErrorOnConflict(t *testing.T) {
	sources := []map[string]string{
		{"KEY": "first"},
		{"KEY": "second"},
	}
	_, err := DedupeSecrets(sources, DedupeError)
	if err == nil {
		t.Fatal("expected error on duplicate key")
	}
}
