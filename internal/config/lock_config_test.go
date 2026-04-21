package config

import (
	"testing"
	"time"
)

func TestBuildLockOption_Nil(t *testing.T) {
	_, enabled, err := BuildLockOption(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if enabled {
		t.Fatal("expected locking disabled for nil config")
	}
}

func TestBuildLockOption_Disabled(t *testing.T) {
	cfg := &LockConfig{Enabled: false}
	_, enabled, err := BuildLockOption(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if enabled {
		t.Fatal("expected disabled")
	}
}

func TestBuildLockOption_Defaults(t *testing.T) {
	cfg := &LockConfig{Enabled: true}
	opt, enabled, err := BuildLockOption(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if !enabled {
		t.Fatal("expected enabled")
	}
	if opt.StaleTTL != 5*time.Minute {
		t.Errorf("expected default StaleTTL 5m, got %s", opt.StaleTTL)
	}
}

func TestBuildLockOption_CustomStaleTTL(t *testing.T) {
	cfg := &LockConfig{Enabled: true, StaleTTL: "2m"}
	opt, _, err := BuildLockOption(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if opt.StaleTTL != 2*time.Minute {
		t.Errorf("expected 2m, got %s", opt.StaleTTL)
	}
}

func TestBuildLockOption_InvalidStaleTTL(t *testing.T) {
	cfg := &LockConfig{Enabled: true, StaleTTL: "not-a-duration"}
	_, _, err := BuildLockOption(cfg)
	if err == nil {
		t.Fatal("expected error for invalid stale_ttl")
	}
}

func TestBuildLockOption_ZeroStaleTTL(t *testing.T) {
	cfg := &LockConfig{Enabled: true, StaleTTL: "0s"}
	_, _, err := BuildLockOption(cfg)
	if err == nil {
		t.Fatal("expected error for zero stale_ttl")
	}
}
