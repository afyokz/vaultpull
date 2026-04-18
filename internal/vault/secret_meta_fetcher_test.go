package vault

import (
	"testing"
	"time"
)

func TestSecretMeta_Fields(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	m := SecretMeta{
		Path:        "secret/myapp/db",
		Version:     3,
		CreatedTime: now,
	}
	if m.Path != "secret/myapp/db" {
		t.Errorf("unexpected path: %s", m.Path)
	}
	if m.Version != 3 {
		t.Errorf("unexpected version: %d", m.Version)
	}
	if !m.CreatedTime.Equal(now) {
		t.Errorf("unexpected created time")
	}
}

func TestCheckRotation_BoundaryExact(t *testing.T) {
	policy := RotationPolicy{MaxAgeDays: 30}
	meta := map[string]SecretMeta{
		"secret/boundary": {
			Path:        "secret/boundary",
			Version:     1,
			CreatedTime: time.Now().AddDate(0, 0, -30),
		},
	}
	results := CheckRotation(meta, policy)
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if !results[0].NeedsRotation {
		t.Error("expected NeedsRotation=true at exact boundary")
	}
}

func TestCheckRotation_MultipleSecrets(t *testing.T) {
	policy := RotationPolicy{MaxAgeDays: 60}
	meta := map[string]SecretMeta{
		"secret/a": {CreatedTime: time.Now().AddDate(0, 0, -10)},
		"secret/b": {CreatedTime: time.Now().AddDate(0, 0, -90)},
		"secret/c": {CreatedTime: time.Now().AddDate(0, 0, -61)},
	}
	results := CheckRotation(meta, policy)
	needed := 0
	for _, r := range results {
		if r.NeedsRotation {
			needed++
		}
	}
	if needed != 2 {
		t.Errorf("expected 2 needing rotation, got %d", needed)
	}
}
