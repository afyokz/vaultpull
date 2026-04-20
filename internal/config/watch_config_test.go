package config

import (
	"testing"
	"time"
)

func TestBuildWatchOption_Nil(t *testing.T) {
	opt, err := BuildWatchOption(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opt.Interval <= 0 {
		t.Error("expected positive default interval")
	}
	if opt.MaxErrors <= 0 {
		t.Error("expected positive default max errors")
	}
}

func TestBuildWatchOption_ValidInterval(t *testing.T) {
	cfg := &WatchConfig{Interval: "1m", MaxErrors: 5}
	opt, err := BuildWatchOption(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opt.Interval != time.Minute {
		t.Errorf("expected 1m, got %s", opt.Interval)
	}
	if opt.MaxErrors != 5 {
		t.Errorf("expected 5, got %d", opt.MaxErrors)
	}
}

func TestBuildWatchOption_InvalidInterval(t *testing.T) {
	cfg := &WatchConfig{Interval: "notaduration"}
	_, err := BuildWatchOption(cfg)
	if err == nil {
		t.Fatal("expected error for invalid interval")
	}
}

func TestBuildWatchOption_ZeroInterval(t *testing.T) {
	cfg := &WatchConfig{Interval: "0s"}
	_, err := BuildWatchOption(cfg)
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestBuildWatchOption_DefaultMaxErrors(t *testing.T) {
	cfg := &WatchConfig{Interval: "10s"}
	opt, err := BuildWatchOption(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defaults, _ := BuildWatchOption(nil)
	if opt.MaxErrors != defaults.MaxErrors {
		t.Errorf("expected default max errors %d, got %d", defaults.MaxErrors, opt.MaxErrors)
	}
}
