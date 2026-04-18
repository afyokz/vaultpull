package diff

import (
	"testing"
)

func TestCompute_AllUnchanged(t *testing.T) {
	existing := map[string]string{"FOO": "bar", "BAZ": "qux"}
	incoming := map[string]string{"FOO": "bar", "BAZ": "qux"}
	r := Compute(existing, incoming)
	if r.HasChanges() {
		t.Error("expected no changes")
	}
}

func TestCompute_Added(t *testing.T) {
	r := Compute(map[string]string{}, map[string]string{"NEW_KEY": "val"})
	if len(r.Changes) != 1 || r.Changes[0].Type != Added {
		t.Errorf("expected Added, got %+v", r.Changes)
	}
}

func TestCompute_Modified(t *testing.T) {
	r := Compute(map[string]string{"KEY": "old"}, map[string]string{"KEY": "new"})
	if len(r.Changes) != 1 || r.Changes[0].Type != Modified {
		t.Errorf("expected Modified, got %+v", r.Changes)
	}
	if r.Changes[0].OldVal != "old" || r.Changes[0].NewVal != "new" {
		t.Error("old/new values incorrect")
	}
}

func TestCompute_Removed(t *testing.T) {
	r := Compute(map[string]string{"GONE": "val"}, map[string]string{})
	if len(r.Changes) != 1 || r.Changes[0].Type != Removed {
		t.Errorf("expected Removed, got %+v", r.Changes)
	}
}

func TestResult_Summary(t *testing.T) {
	r := Compute(
		map[string]string{"OLD": "v", "SAME": "v"},
		map[string]string{"NEW": "v", "SAME": "v", "MOD": "new"},
	)
	// OLD removed, NEW added, MOD added (not in existing), SAME unchanged
	summary := r.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestCompute_Mixed(t *testing.T) {
	existing := map[string]string{"A": "1", "B": "2", "C": "3"}
	incoming := map[string]string{"A": "1", "B": "changed", "D": "4"}
	r := Compute(existing, incoming)
	if !r.HasChanges() {
		t.Error("expected changes")
	}
	types := map[ChangeType]int{}
	for _, c := range r.Changes {
		types[c.Type]++
	}
	if types[Unchanged] != 1 || types[Modified] != 1 || types[Removed] != 1 || types[Added] != 1 {
		t.Errorf("unexpected change breakdown: %+v", types)
	}
}
