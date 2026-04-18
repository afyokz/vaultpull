package vault

import "strings"

// Filter holds include/exclude key patterns for secret filtering.
type Filter struct {
	Include []string
	Exclude []string
}

// Apply filters a map of secrets, returning only keys that pass the filter rules.
// If Include is non-empty, only matching keys are kept.
// Keys matching any Exclude pattern are always removed.
func (f *Filter) Apply(secrets map[string]string) map[string]string {
	result := make(map[string]string)

	for k, v := range secrets {
		if f.excluded(k) {
			continue
		}
		if len(f.Include) > 0 && !f.included(k) {
			continue
		}
		result[k] = v
	}

	return result
}

func (f *Filter) included(key string) bool {
	for _, pattern := range f.Include {
		if matchPattern(pattern, key) {
			return true
		}
	}
	return false
}

func (f *Filter) excluded(key string) bool {
	for _, pattern := range f.Exclude {
		if matchPattern(pattern, key) {
			return true
		}
	}
	return false
}

// matchPattern supports simple prefix wildcard matching (e.g. "DB_*").
func matchPattern(pattern, key string) bool {
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(key, strings.TrimSuffix(pattern, "*"))
	}
	return pattern == key
}
