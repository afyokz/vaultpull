package diff

import "fmt"

// ChangeType represents the kind of change for a secret key.
type ChangeType string

const (
	Added    ChangeType = "added"
	Modified ChangeType = "modified"
	Removed  ChangeType = "removed"
	Unchanged ChangeType = "unchanged"
)

// Change describes a single key-level difference.
type Change struct {
	Key    string
	Type   ChangeType
	OldVal string
	NewVal string
}

// Result holds the full diff between existing and incoming secrets.
type Result struct {
	Changes []Change
}

// HasChanges returns true if any non-unchanged entries exist.
func (r *Result) HasChanges() bool {
	for _, c := range r.Changes {
		if c.Type != Unchanged {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary string.
func (r *Result) Summary() string {
	var added, modified, removed int
	for _, c := range r.Changes {
		switch c.Type {
		case Added:
			added++
		case Modified:
			modified++
		case Removed:
			removed++
		}
	}
	return fmt.Sprintf("+%d added, ~%d modified, -%d removed", added, modified, removed)
}

// Compute returns a Result comparing existing env map to incoming secrets.
func Compute(existing, incoming map[string]string) Result {
	var changes []Change

	for k, newVal := range incoming {
		if oldVal, ok := existing[k]; ok {
			if oldVal != newVal {
				changes = append(changes, Change{Key: k, Type: Modified, OldVal: oldVal, NewVal: newVal})
			} else {
				changes = append(changes, Change{Key: k, Type: Unchanged})
			}
		} else {
			changes = append(changes, Change{Key: k, Type: Added, NewVal: newVal})
		}
	}

	for k, oldVal := range existing {
		if _, ok := incoming[k]; !ok {
			changes = append(changes, Change{Key: k, Type: Removed, OldVal: oldVal})
		}
	}

	return Result{Changes: changes}
}
