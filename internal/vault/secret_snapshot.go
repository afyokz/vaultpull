package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of fetched secrets.
type Snapshot struct {
	CapturedAt time.Time         `json:"captured_at"`
	Secrets    map[string]string `json:"secrets"`
}

// SaveSnapshot writes secrets to a JSON snapshot file.
func SaveSnapshot(path string, secrets map[string]string) error {
	snap := Snapshot{
		CapturedAt: time.Now().UTC(),
		Secrets:    secrets,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("snapshot: write failed: %w", err)
	}
	return nil
}

// LoadSnapshot reads a snapshot file and returns the Snapshot.
func LoadSnapshot(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read failed: %w", err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal failed: %w", err)
	}
	return &snap, nil
}

// DiffSnapshot compares a snapshot against current secrets and returns added,
// removed, and changed keys.
func DiffSnapshot(snap *Snapshot, current map[string]string) (added, removed, changed []string) {
	for k, v := range current {
		old, ok := snap.Secrets[k]
		if !ok {
			added = append(added, k)
		} else if old != v {
			changed = append(changed, k)
		}
	}
	for k := range snap.Secrets {
		if _, ok := current[k]; !ok {
			removed = append(removed, k)
		}
	}
	return
}
