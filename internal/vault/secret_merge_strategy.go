package vault

import "fmt"

// MergeStrategy defines how conflicts are resolved when merging secrets.
type MergeStrategy int

const (
	// MergeStrategyOverwrite replaces existing values with incoming ones.
	MergeStrategyOverwrite MergeStrategy = iota
	// MergeStrategyKeepExisting preserves existing values on conflict.
	MergeStrategyKeepExisting
	// MergeStrategyError returns an error on any conflict.
	MergeStrategyError
)

// ParseMergeStrategy parses a strategy name into a MergeStrategy.
func ParseMergeStrategy(s string) (MergeStrategy, error) {
	switch s {
	case "overwrite", "":
		return MergeStrategyOverwrite, nil
	case "keep":
		return MergeStrategyKeepExisting, nil
	case "error":
		return MergeStrategyError, nil
	default:
		return 0, fmt.Errorf("unknown merge strategy %q: must be overwrite, keep, or error", s)
	}
}

// MergeWithStrategy merges src into dst using the given strategy.
// dst is modified in place.
func MergeWithStrategy(dst, src map[string]string, strategy MergeStrategy) error {
	for k, v := range src {
		existing, exists := dst[k]
		switch strategy {
		case MergeStrategyOverwrite:
			dst[k] = v
		case MergeStrategyKeepExisting:
			if !exists {
				dst[k] = v
			}
			_ = existing
		case MergeStrategyError:
			if exists {
				return fmt.Errorf("conflict: key %q exists in destination (value %q)", k, existing)
			}
			dst[k] = v
		}
	}
	return nil
}
