package config

import (
	"testing"
)

func TestBuildMaskOption_Nil(t *testing.T) {
	if BuildMaskOption(nil) != nil {
		t.Error("expected nil for nil config")
	}
}

func TestBuildMaskOption_Disabled(t *testing.T) {
	cfg := &MaskConfig{Enabled: false, Keys: []string{"TOKEN"}}
	if BuildMaskOption(cfg) != nil {
		t.Error("expected nil when disabled")
	}
}

func TestBuildMaskOption_Defaults(t *testing.T) {
	cfg := &MaskConfig{Enabled: true}
	opt := BuildMaskOption(cfg)
	if opt == nil {
		t.Fatal("expected non-nil option")
	}
	if opt.ShowPrefix != 2 {
		t.Errorf("expected default ShowPrefix=2, got %d", opt.ShowPrefix)
	}
	if opt.Replacement != "****" {
		t.Errorf("expected default replacement, got %q", opt.Replacement)
	}
}

func TestBuildMaskOption_CustomReplacement(t *testing.T) {
	cfg := &MaskConfig{Enabled: true, Replacement: "[hidden]", ShowPrefix: 3, ShowSuffix: 1}
	opt := BuildMaskOption(cfg)
	if opt.Replacement != "[hidden]" {
		t.Errorf("expected [hidden], got %q", opt.Replacement)
	}
	if opt.ShowPrefix != 3 {
		t.Errorf("expected ShowPrefix=3, got %d", opt.ShowPrefix)
	}
}
