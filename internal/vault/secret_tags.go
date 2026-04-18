package vault

import "strings"

// TagFilter holds key=value tag pairs used to filter secrets.
type TagFilter struct {
	tags map[string]string
}

// NewTagFilter parses a slice of "key=value" strings into a TagFilter.
func NewTagFilter(pairs []string) (*TagFilter, error) {
	tags := make(map[string]string, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			return nil, &TagParseError{Raw: p}
		}
		tags[parts[0]] = parts[1]
	}
	return &TagFilter{tags: tags}, nil
}

// TagParseError is returned when a tag pair cannot be parsed.
type TagParseError struct {
	Raw string
}

func (e *TagParseError) Error() string {
	return "invalid tag format (expected key=value): " + e.Raw
}

// Match reports whether the provided metadata tags satisfy all filter tags.
func (f *TagFilter) Match(meta map[string]string) bool {
	if len(f.tags) == 0 {
		return true
	}
	for k, v := range f.tags {
		if meta[k] != v {
			return false
		}
	}
	return true
}

// FilterByTags removes secrets whose metadata does not match the filter.
// secrets is a map of secret key -> flat kv pairs; metaMap maps secret key -> tags.
func FilterByTags(secrets map[string]string, metaMap map[string]map[string]string, f *TagFilter) map[string]string {
	if f == nil || len(f.tags) == 0 {
		return secrets
	}
	result := make(map[string]string)
	for k, v := range secrets {
		if f.Match(metaMap[k]) {
			result[k] = v
		}
	}
	return result
}
