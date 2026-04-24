package vault

import (
	"strings"
	"testing"
	"time"
)

func TestCheckTTL_Empty(t *testing.T) {
	results := CheckTTL(map[string]time.Time{}, DefaultTTLPolicy())
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestCheckTTL_SkipsZeroTime(t *testing.T) {
	secrets := map[string]time.Time{
		"EMPTY_KEY": {},
	}
	results := CheckTTL(secrets, DefaultTTLPolicy())
	if len(results) != 0 {
		t.Fatalf("expected zero-time entries to be skipped, got %d results", len(results))
	}
}

func TestCheckTTL_Expired(t *testing.T) {
	secrets := map[string]time.Time{
		"OLD_TOKEN": time.Now().Add(-1 * time.Hour),
	}
	policy := DefaultTTLPolicy()
	results := CheckTTL(secrets, policy)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != "expired" {
		t.Errorf("expected status 'expired', got %q", results[0].Status)
	}
}

func TestCheckTTL_Warn(t *testing.T) {
	secrets := map[string]time.Time{
		"NEAR_TOKEN": time.Now().Add(12 * time.Hour),
	}
	policy := DefaultTTLPolicy() // warn threshold = 24h
	results := CheckTTL(secrets, policy)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != "warn" {
		t.Errorf("expected status 'warn', got %q", results[0].Status)
	}
}

func TestCheckTTL_OK(t *testing.T) {
	secrets := map[string]time.Time{
		"FRESH_TOKEN": time.Now().Add(48 * time.Hour),
	}
	results := CheckTTL(secrets, DefaultTTLPolicy())
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != "ok" {
		t.Errorf("expected status 'ok', got %q", results[0].Status)
	}
}

func TestFormatTTLReport_NoResults(t *testing.T) {
	out := FormatTTLReport(nil)
	if !strings.Contains(out, "no expiry data") {
		t.Errorf("expected 'no expiry data' message, got: %q", out)
	}
}

func TestFormatTTLReport_ContainsStatus(t *testing.T) {
	results := []TTLResult{
		{Key: "API_KEY", ExpiresAt: time.Now().Add(2 * time.Hour), Remaining: 2 * time.Hour, Status: "warn"},
	}
	out := FormatTTLReport(results)
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN in report, got: %q", out)
	}
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in report, got: %q", out)
	}
}
