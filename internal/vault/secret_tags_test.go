package vault

import (
	"testing"
)

func TestNewTagFilter_Valid(t *testing.T) {
	f, err := NewTagFilter([]string{"env=prod", "team=platform"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(f.tags))
	}
}

func TestNewTagFilter_Invalid(t *testing.T) {
	_, err := NewTagFilter([]string{"badformat"})
	if err == nil {
		t.Fatal("expected error for bad format")
	}
}

func TestNewTagFilter_EmptyKey(t *testing.T) {
	_, err := NewTagFilter([]string{"=value"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestTagFilter_Match_NoRules(t *testing.T) {
	f := &TagFilter{tags: map[string]string{}}
	if !f.Match(map[string]string{"env": "dev"}) {
		t.Fatal("empty filter should match anything")
	}
}

func TestTagFilter_Match_Hit(t *testing.T) {
	f, _ := NewTagFilter([]string{"env=prod"})
	if !f.Match(map[string]string{"env": "prod", "region": "us"}) {
		t.Fatal("expected match")
	}
}

func TestTagFilter_Match_Miss(t *testing.T) {
	f, _ := NewTagFilter([]string{"env=prod"})
	if f.Match(map[string]string{"env": "dev"}) {
		t.Fatal("expected no match")
	}
}

func TestFilterByTags_RemovesNonMatching(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "x", "API_KEY": "y"}
	meta := map[string]map[string]string{
		"DB_PASS": {"env": "prod"},
		"API_KEY": {"env": "dev"},
	}
	f, _ := NewTagFilter([]string{"env=prod"})
	result := FilterByTags(secrets, meta, f)
	if len(result) != 1 {
		t.Fatalf("expected 1 secret, got %d", len(result))
	}
	if _, ok := result["DB_PASS"]; !ok {
		t.Fatal("expected DB_PASS in result")
	}
}

func TestFilterByTags_NilFilter(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	result := FilterByTags(secrets, nil, nil)
	if len(result) != 2 {
		t.Fatalf("expected 2 secrets, got %d", len(result))
	}
}
