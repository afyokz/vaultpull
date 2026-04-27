package vault

import (
	"strings"
	"testing"
	"time"
)

func TestRollbackStore_PushAndList(t *testing.T) {
	s := NewRollbackStore(5)
	s.Push("first", map[string]string{"A": "1"})
	time.Sleep(time.Millisecond)
	s.Push("second", map[string]string{"B": "2"})

	list := s.List()
	if len(list) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(list))
	}
	// newest first
	if list[0].Label != "second" {
		t.Errorf("expected second entry first, got %q", list[0].Label)
	}
}

func TestRollbackStore_Prunes(t *testing.T) {
	s := NewRollbackStore(3)
	for i := 0; i < 5; i++ {
		s.Push("", map[string]string{})
	}
	if len(s.List()) != 3 {
		t.Errorf("expected 3 entries after pruning, got %d", len(s.List()))
	}
}

func TestRollbackStore_GetByID(t *testing.T) {
	s := NewRollbackStore(5)
	entry := s.Push("tagged", map[string]string{"X": "y"})

	got, ok := s.Get(entry.ID)
	if !ok {
		t.Fatal("expected to find entry by ID")
	}
	if got.Label != "tagged" {
		t.Errorf("expected label 'tagged', got %q", got.Label)
	}
}

func TestRollbackStore_GetMissing(t *testing.T) {
	s := NewRollbackStore(5)
	_, ok := s.Get("nonexistent")
	if ok {
		t.Error("expected not found for missing ID")
	}
}

func TestRollbackStore_Latest_Empty(t *testing.T) {
	s := NewRollbackStore(5)
	_, ok := s.Latest()
	if ok {
		t.Error("expected false for empty store")
	}
}

func TestRollbackStore_Latest_ReturnsCopy(t *testing.T) {
	s := NewRollbackStore(5)
	s.Push("v1", map[string]string{"K": "v"})
	entry, ok := s.Latest()
	if !ok {
		t.Fatal("expected entry")
	}
	entry.Secrets["K"] = "mutated"
	latest, _ := s.Latest()
	if latest.Secrets["K"] == "mutated" {
		t.Error("Push should store a copy; mutation affected stored entry")
	}
}

func TestFormatRollbackList_Empty(t *testing.T) {
	out := FormatRollbackList(nil)
	if !strings.Contains(out, "no rollback") {
		t.Errorf("expected empty message, got %q", out)
	}
}

func TestFormatRollbackList_ContainsLabel(t *testing.T) {
	s := NewRollbackStore(5)
	s.Push("my-label", map[string]string{"A": "1", "B": "2"})
	out := FormatRollbackList(s.List())
	if !strings.Contains(out, "my-label") {
		t.Errorf("expected label in output, got %q", out)
	}
	if !strings.Contains(out, "2 keys") {
		t.Errorf("expected key count in output, got %q", out)
	}
}
