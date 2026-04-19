package vault

import "strings"

// DedupeStrategy controls how duplicate keys are resolved.
type DedupeStrategy string

const (
	DedupeKeepFirst DedupeStrategy = "keep-first"
	DedupeKeepLast  DedupeStrategy = "keep-last"
	DedupeError     DedupeStrategy = "error"
)

// DedupeResult holds the outcome of a deduplication pass.
type DedupeResult struct {
	Secrets    map[string]string
	Duplicates []string
}

// ParseDedupeStrategy parses a strategy string, returning an error if unknown.
func ParseDedupeStrategy(s string) (DedupeStrategy, error) {
	switch DedupeStrategy(strings.ToLower(s)) {
	case DedupeKeepFirst, DedupeKeepLast, DedupeError:
		return DedupeStrategy(strings.ToLower(s)), nil
	}
	return "", fmt.Errorf("unknown dedupe strategy %q: must be keep-first, keep-last, or error", s)
}

// DedupeSecrets merges a slice of secret maps according to the given strategy.
// Later maps in the slice are considered duplicates of earlier ones for the
// same key.
func DedupeSecrets(sources []map[string]string, strategy DedupeStrategy) (DedupeResult, error) {
	seen := make(map[string]bool)
	out := make(map[string]string)
	var dupes []string

	for i, src := range sources {
		for k, v := range src {
			if seen[k] {
				dupes = append(dupes, k)
				switch strategy {
				case DedupeError:
					return DedupeResult{}, fmt.Errorf("duplicate key %q found in source %d", k, i)
				case DedupeKeepLast:
					out[k] = v
				// DedupeKeepFirst: do nothing
				}
			} else {
				out[k] = v
				seen[k] = true
			}
		}
	}

	return DedupeResult{Secrets: out, Duplicates: dupes}, nil
}
