package vault

import (
	"fmt"
	"time"
)

// RotateResult holds the outcome of a rotation check for a single secret.
type RotateResult struct {
	Path      string
	LastRotated time.Time
	AgeDays   int
	NeedsRotation bool
}

// RotationPolicy defines when secrets should be rotated.
type RotationPolicy struct {
	MaxAgeDays int
}

// DefaultRotationPolicy returns a policy with a 90-day max age.
func DefaultRotationPolicy() RotationPolicy {
	return RotationPolicy{MaxAgeDays: 90}
}

// CheckRotation evaluates secrets metadata against the policy.
func CheckRotation(meta map[string]SecretMeta, policy RotationPolicy) []RotateResult {
	now := time.Now()
	results := make([]RotateResult, 0, len(meta))
	for path, m := range meta {
		age := int(now.Sub(m.CreatedTime).Hours() / 24)
		results = append(results, RotateResult{
			Path:          path,
			LastRotated:   m.CreatedTime,
			AgeDays:       age,
			NeedsRotation: age >= policy.MaxAgeDays,
		})
	}
	return results
}

// SecretMeta holds metadata about a secret version.
type SecretMeta struct {
	Path        string
	Version     int
	CreatedTime time.Time
}

// FormatRotationReport returns a human-readable rotation report.
func FormatRotationReport(results []RotateResult) string {
	if len(results) == 0 {
		return "No secrets to evaluate."
	}
	out := ""
	for _, r := range results {
		status := "OK"
		if r.NeedsRotation {
			status = "ROTATION NEEDED"
		}
		out += fmt.Sprintf("  [%s] %s (age: %d days)\n", status, r.Path, r.AgeDays)
	}
	return out
}
