package config

import (
	"testing"
)

func TestBuildRollbackStore_Nil(t *testing.T) {
	store, err := BuildRollbackStore(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store != nil {
		t.Error("expected nil store for nil config")
	}
}

func TestBuildRollbackStore_Disabled(t *testing.T) {
	cfg := &RollbackConfig{Enabled: false, MaxSize: 5}
	store, err := BuildRollbackStore(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store != nil {
		t.Error("expected nil store when disabled")
	}
}

func TestBuildRollbackStore_Defaults(t *testing.T) {
	cfg := &RollbackConfig{Enabled: true}
	store, err := BuildRollbackStore(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store == nil {
		t.Fatal("expected non-nil store")
	}
}

func TestBuildRollbackStore_CustomMaxSize(t *testing.T) {
	cfg := &RollbackConfig{Enabled: true, MaxSize: 3}
	store, err := BuildRollbackStore(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store == nil {
		t.Fatal("expected non-nil store")
	}
	// Push 5 entries, expect only 3 retained
	for i := 0; i < 5; i++ {
		store.Push("", map[string]string{})
	}
	if len(store.List()) != 3 {
		t.Errorf("expected 3 entries, got %d", len(store.List()))
	}
}

func TestBuildRollbackStore_NegativeMaxSize(t *testing.T) {
	cfg := &RollbackConfig{Enabled: true, MaxSize: -1}
	_, err := BuildRollbackStore(cfg)
	if err == nil {
		t.Error("expected error for negative max_size")
	}
}
