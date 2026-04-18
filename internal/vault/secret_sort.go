package vault

import (
	"sort"
	"strings"
)

// SortOrder defines the ordering direction.
type SortOrder int

const (
	SortAsc  SortOrder = iota
	SortDesc
)

// SortOption configures how secrets are sorted.
type SortOption struct {
	Order  SortOrder
	ByKey  bool
	ByValue bool
}

// DefaultSortOption returns ascending key-based sort.
func DefaultSortOption() SortOption {
	return SortOption{Order: SortAsc, ByKey: true}
}

// SortSecrets returns a new map's keys sorted into a slice of key-value pairs.
func SortSecrets(secrets map[string]string, opt SortOption) []KeyValue {
	pairs := make([]KeyValue, 0, len(secrets))
	for k, v := range secrets {
		pairs = append(pairs, KeyValue{Key: k, Value: v})
	}

	sort.Slice(pairs, func(i, j int) bool {
		var less bool
		if opt.ByValue {
			less = strings.ToLower(pairs[i].Value) < strings.ToLower(pairs[j].Value)
		} else {
			less = strings.ToLower(pairs[i].Key) < strings.ToLower(pairs[j].Key)
		}
		if opt.Order == SortDesc {
			return !less
		}
		return less
	})

	return pairs
}

// KeyValue holds a secret key-value pair.
type KeyValue struct {
	Key   string
	Value string
}
