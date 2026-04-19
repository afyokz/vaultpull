package vault

import (
	"testing"
)

func TestParseMergeStrategy_Valid(t *testing.T) {
	cases := []struct {
		input    string
		expected MergeStrategy
	}{
		{"overwrite", MergeStrategyOverwrite},
		{"", MergeStrategyOverwrite},
		{"keep", MergeStrategyKeepExisting},
		{"error", MergeStrategyError},
	}
	for _, tc := range cases {
		got, err := ParseMergeStrategy(tc.input)
		if err != nil {
			t.Errorf("ParseMergeStrategy(%q) unexpected error: %v", tc.input, err)
		}
		if got != tc.expected {
			t.Errorf("ParseMergeStrategy(%q) = %v, want %v", tc.input, got, tc.expected)
		}
	}
}

func TestParseMergeStrategy_Invalid(t *testing.T) {
	_, err := ParseMergeStrategy("unknown")
	if err == nil {
		t.Error("expected error for unknown strategy")
	}
}

func TestMergeWithStrategy_Overwrite(t *testing.T) {
	dst := map[string]string{"A": "old"}
	src := map[string]string{"A": "new", "B": "val"}
	if err := MergeWithStrategy(dst, src, MergeStrategyOverwrite); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["A"] != "new" {
		t.Errorf("expected A=new, got %q", dst["A"])
	}
	if dst["B"] != "val" {
		t.Errorf("expected B=val, got %q", dst["B"])
	}
}

func TestMergeWithStrategy_KeepExisting(t *testing.T) {
	dst := map[string]string{"A": "old"}
	src := map[string]string{"A": "new", "B": "val"}
	if err := MergeWithStrategy(dst, src, MergeStrategyKeepExisting); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["A"] != "old" {
		t.Errorf("expected A=old, got %q", dst["A"])
	}
	if dst["B"] != "val" {
		t.Errorf("expected B=val, got %q", dst["B"])
	}
}

func TestMergeWithStrategy_ErrorOnConflict(t *testing.T) {
	dst := map[string]string{"A": "old"}
	src := map[string]string{"A": "new"}
	err := MergeWithStrategy(dst, src, MergeStrategyError)
	if err == nil {
		t.Error("expected conflict error")
	}
}

func TestMergeWithStrategy_ErrorNoConflict(t *testing.T) {
	dst := map[string]string{"A": "old"}
	src := map[string]string{"B": "new"}
	if err := MergeWithStrategy(dst, src, MergeStrategyError); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if dst["B"] != "new" {
		t.Errorf("expected B=new, got %q", dst["B"])
	}
}
