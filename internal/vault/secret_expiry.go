package vault

import (
	"fmt"
	"time"
)

// ExpiryPolicy defines when a secret is considered expired.
type ExpiryPolicy struct {
	MaxAge time.Duration
}

// ExpiryResult holds the expiry status for a single secret.
type ExpiryResult struct {
	Key       string
	ExpiresAt time.Time
	Expired   bool
	DaysLeft  int
}

// DefaultExpiryPolicy returns a policy with a 90-day max age.
func DefaultExpiryPolicy() ExpiryPolicy {
	return ExpiryPolicy{MaxAge: 90 * 24 * time.Hour}
}

// CheckExpiry evaluates secret metadata against the expiry policy.
func CheckExpiry(policy ExpiryPolicy, meta map[string]SecretMeta) []ExpiryResult {
	now := time.Now()
	results := make([]ExpiryResult, 0, len(meta))
	for key, m := range meta {
		if m.CreatedAt.IsZero() {
			continue
		}
		expiresAt := m.CreatedAt.Add(policy.MaxAge)
		expired := now.After(expiresAt)
		daysLeft := int(time.Until(expiresAt).Hours() / 24)
		results = append(results, ExpiryResult{
			Key:       key,
			ExpiresAt: expiresAt,
			Expired:   expired,
			DaysLeft:  daysLeft,
		})
	}
	return results
}

// FormatExpiryReport returns a human-readable summary of expiry results.
func FormatExpiryReport(results []ExpiryResult) string {
	if len(results) == 0 {
		return "No expiry data available.\n"
	}
	out := ""
	for _, r := range results {
		status := fmt.Sprintf("%d days left", r.DaysLeft)
		if r.Expired {
			status = "EXPIRED"
		}
		out += fmt.Sprintf("  %-40s %s (expires %s)\n", r.Key, status, r.ExpiresAt.Format("2006-01-02"))
	}
	return out
}
