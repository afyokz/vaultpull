package vault

import (
	"strings"
	"testing"
	"time"
)

func TestCheckExpiry_Empty(t *testing.T) {
	results := CheckExpiry(DefaultExpiryPolicy(), map[string]SecretMeta{})
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestCheckExpiry_Expired(t *testing.T) {
	policy := DefaultExpiryPolicy()
	meta := map[string]SecretMeta{
		"OLD_KEY": {CreatedAt: time.Now().Add(-100 * 24 * time.Hour)},
	}
	results := CheckExpiry(policy, meta)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Expired {
		t.Error("expected secret to be expired")
	}
}

func TestCheckExpiry_Fresh(t *testing.T) {
	policy := DefaultExpiryPolicy()
	meta := map[string]SecretMeta{
		"NEW_KEY": {CreatedAt: time.Now().Add(-10 * 24 * time.Hour)},
	}
	results := CheckExpiry(policy, meta)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Expired {
		t.Error("expected secret to be fresh")
	}
	if results[0].DaysLeft <= 0 {
		t.Errorf("expected positive days left, got %d", results[0].DaysLeft)
	}
}

func TestCheckExpiry_SkipsZeroTime(t *testing.T) {
	policy := DefaultExpiryPolicy()
	meta := map[string]SecretMeta{
		"NO_DATE": {CreatedAt: time.Time{}},
	}
	results := CheckExpiry(policy, meta)
	if len(results) != 0 {
		t.Error("expected zero-time entries to be skipped")
	}
}

func TestFormatExpiryReport_NoResults(t *testing.T) {
	out := FormatExpiryReport(nil)
	if !strings.Contains(out, "No expiry") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatExpiryReport_ContainsStatus(t *testing.T) {
	results := []ExpiryResult{
		{Key: "MY_SECRET", ExpiresAt: time.Now().Add(5 * 24 * time.Hour), Expired: false, DaysLeft: 5},
	}
	out := FormatExpiryReport(results)
	if !strings.Contains(out, "MY_SECRET") {
		t.Error("expected key in report")
	}
	if !strings.Contains(out, "5 days left") {
		t.Error("expected days left in report")
	}
}
