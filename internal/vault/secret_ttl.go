package vault

import (
	"fmt"
	"strings"
	"time"
)

// TTLPolicy defines how secret TTLs are evaluated.
type TTLPolicy struct {
	WarnThreshold time.Duration
	ErrorThreshold time.Duration
}

// TTLResult holds the TTL evaluation result for a single secret.
type TTLResult struct {
	Key       string
	ExpiresAt time.Time
	Remaining time.Duration
	Status    string // "ok", "warn", "expired"
}

// DefaultTTLPolicy returns a sensible default TTL policy.
func DefaultTTLPolicy() TTLPolicy {
	return TTLPolicy{
		WarnThreshold:  24 * time.Hour,
		ErrorThreshold: 0,
	}
}

// CheckTTL evaluates TTLs for a map of secret keys to expiry times.
func CheckTTL(secrets map[string]time.Time, policy TTLPolicy) []TTLResult {
	now := time.Now()
	results := make([]TTLResult, 0, len(secrets))

	for key, expiresAt := range secrets {
		if expiresAt.IsZero() {
			continue
		}
		remaining := expiresAt.Sub(now)
		status := "ok"
		switch {
		case remaining <= policy.ErrorThreshold:
			status = "expired"
		case remaining <= policy.WarnThreshold:
			status = "warn"
		}
		results = append(results, TTLResult{
			Key:       key,
			ExpiresAt: expiresAt,
			Remaining: remaining,
			Status:    status,
		})
	}
	return results
}

// FormatTTLReport formats TTL check results into a human-readable string.
func FormatTTLReport(results []TTLResult) string {
	if len(results) == 0 {
		return "TTL check: no expiry data found."
	}
	var sb strings.Builder
	sb.WriteString("TTL Report:\n")
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("  [%s] %s — expires %s (in %s)\n",
			strings.ToUpper(r.Status),
			r.Key,
			r.ExpiresAt.Format(time.RFC3339),
			r.Remaining.Round(time.Second),
		))
	}
	return sb.String()
}
