package vault

import (
	"strings"
	"testing"
	"time"
)

func TestCheckRotation_NeedsRotation(t *testing.T) {
	policy := RotationPolicy{MaxAgeDays: 30}
	meta := map[string]SecretMeta{
		"secret/old": {Path: "secret/old", Version: 1, CreatedTime: time.Now().AddDate(0, 0, -60)},
	}
	results := CheckRotation(meta, policy)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].NeedsRotation {
		t.Error("expected NeedsRotation=true")
	}
}

func TestCheckRotation_Fresh(t *testing.T) {
	policy := RotationPolicy{MaxAgeDays: 90}
	meta := map[string]SecretMeta{
		"secret/new": {Path: "secret/new", Version: 2, CreatedTime: time.Now().AddDate(0, 0, -5)},
	}
	results := CheckRotation(meta, policy)
	if results[0].NeedsRotation {
		t.Error("expected NeedsRotation=false")
	}
}

func TestCheckRotation_Empty(t *testing.T) {
	results := CheckRotation(map[string]SecretMeta{}, DefaultRotationPolicy())
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestFormatRotationReport_NoResults(t *testing.T) {
	out := FormatRotationReport([]RotateResult{})
	if out != "No secrets to evaluate." {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatRotationReport_ContainsStatus(t *testing.T) {
	results := []RotateResult{
		{Path: "secret/a", AgeDays: 100, NeedsRotation: true},
		{Path: "secret/b", AgeDays: 10, NeedsRotation: false},
	}
	out := FormatRotationReport(results)
	if !strings.Contains(out, "ROTATION NEEDED") {
		t.Error("expected ROTATION NEEDED in output")
	}
	if !strings.Contains(out, "OK") {
		t.Error("expected OK in output")
	}
}
