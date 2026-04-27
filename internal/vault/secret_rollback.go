package vault

import (
	"fmt"
	"sort"
	"time"
)

// RollbackEntry represents a single rollback point for a set of secrets.
type RollbackEntry struct {
	ID        string
	Timestamp time.Time
	Secrets   map[string]string
	Label     string
}

// RollbackStore holds an ordered list of rollback entries.
type RollbackStore struct {
	entries []RollbackEntry
	maxSize int
}

// DefaultRollbackMaxSize is the default number of rollback entries to retain.
const DefaultRollbackMaxSize = 10

// NewRollbackStore creates a new RollbackStore with the given max size.
func NewRollbackStore(maxSize int) *RollbackStore {
	if maxSize <= 0 {
		maxSize = DefaultRollbackMaxSize
	}
	return &RollbackStore{maxSize: maxSize}
}

// Push adds a new rollback entry, pruning oldest if over capacity.
func (s *RollbackStore) Push(label string, secrets map[string]string) RollbackEntry {
	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}
	entry := RollbackEntry{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Timestamp: time.Now().UTC(),
		Secrets:   copy,
		Label:     label,
	}
	s.entries = append(s.entries, entry)
	if len(s.entries) > s.maxSize {
		s.entries = s.entries[len(s.entries)-s.maxSize:]
	}
	return entry
}

// List returns all entries sorted newest-first.
func (s *RollbackStore) List() []RollbackEntry {
	sorted := make([]RollbackEntry, len(s.entries))
	copy(sorted, s.entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp.After(sorted[j].Timestamp)
	})
	return sorted
}

// Get retrieves an entry by ID, returning false if not found.
func (s *RollbackStore) Get(id string) (RollbackEntry, bool) {
	for _, e := range s.entries {
		if e.ID == id {
			return e, true
		}
	}
	return RollbackEntry{}, false
}

// Latest returns the most recent entry, or false if empty.
func (s *RollbackStore) Latest() (RollbackEntry, bool) {
	if len(s.entries) == 0 {
		return RollbackEntry{}, false
	}
	list := s.List()
	return list[0], true
}

// FormatRollbackList formats entries for display.
func FormatRollbackList(entries []RollbackEntry) string {
	if len(entries) == 0 {
		return "no rollback entries found\n"
	}
	out := ""
	for _, e := range entries {
		label := e.Label
		if label == "" {
			label = "(no label)"
		}
		out += fmt.Sprintf("[%s] %s — %s (%d keys)\n",
			e.ID, e.Timestamp.Format(time.RFC3339), label, len(e.Secrets))
	}
	return out
}
