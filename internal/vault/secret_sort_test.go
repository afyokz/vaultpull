package vault

import (
	"testing"
)

func TestSortSecrets_AscByKey(t *testing.T) {
	secrets := map[string]string{
		"ZEBRA": "1",
		"ALPHA": "2",
		"MANGO": "3",
	}
	result := SortSecrets(secrets, DefaultSortOption())
	if result[0].Key != "ALPHA" || result[1].Key != "MANGO" || result[2].Key != "ZEBRA" {
		t.Errorf("unexpected order: %v", result)
	}
}

func TestSortSecrets_DescByKey(t *testing.T) {
	secrets := map[string]string{
		"ZEBRA": "1",
		"ALPHA": "2",
		"MANGO": "3",
	}
	opt := SortOption{Order: SortDesc, ByKey: true}
	result := SortSecrets(secrets, opt)
	if result[0].Key != "ZEBRA" || result[2].Key != "ALPHA" {
		t.Errorf("unexpected order: %v", result)
	}
}

func TestSortSecrets_ByValue(t *testing.T) {
	secrets := map[string]string{
		"A": "cherry",
		"B": "apple",
		"C": "banana",
	}
	opt := SortOption{Order: SortAsc, ByValue: true}
	result := SortSecrets(secrets, opt)
	if result[0].Value != "apple" || result[1].Value != "banana" || result[2].Value != "cherry" {
		t.Errorf("unexpected order by value: %v", result)
	}
}

func TestSortSecrets_Empty(t *testing.T) {
	result := SortSecrets(map[string]string{}, DefaultSortOption())
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestSortSecrets_CaseInsensitive(t *testing.T) {
	secrets := map[string]string{
		"beta":  "1",
		"ALPHA": "2",
	}
	result := SortSecrets(secrets, DefaultSortOption())
	if result[0].Key != "ALPHA" {
		t.Errorf("expected ALPHA first, got %v", result[0].Key)
	}
}

func TestSortSecrets_SingleEntry(t *testing.T) {
	secrets := map[string]string{
		"ONLY": "value",
	}
	result := SortSecrets(secrets, DefaultSortOption())
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if result[0].Key != "ONLY" || result[0].Value != "value" {
		t.Errorf("unexpected result: %v", result[0])
	}
}
