package vault

import (
	"testing"
)

func TestMaskValue_Empty(t *testing.T) {
	result := MaskValue("", DefaultMaskOption)
	if result != "" {
		t.Errorf("expected empty, got %q", result)
	}
}

func TestMaskValue_ShortValue(t *testing.T) {
	// value shorter than prefix+suffix => full replacement
	result := MaskValue("ab", DefaultMaskOption)
	if result != "****" {
		t.Errorf("expected ****, got %q", result)
	}
}

func TestMaskValue_NormalValue(t *testing.T) {
	result := MaskValue("supersecret", DefaultMaskOption)
	expected := "su****et"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestMaskValue_CustomOption(t *testing.T) {
	opt := MaskOption{ShowPrefix: 1, ShowSuffix: 0, Replacement: "--"}
	result := MaskValue("password", opt)
	expected := "p--"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestMaskSecrets_MasksOnlySpecified(t *testing.T) {
	secrets := map[string]string{
		"API_KEY":  "abcdefgh",
		"APP_NAME": "myapp",
	}
	result := MaskSecrets(secrets, []string{"API_KEY"}, DefaultMaskOption)
	if result["API_KEY"] == "abcdefgh" {
		t.Error("expected API_KEY to be masked")
	}
	if result["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME unchanged, got %q", result["APP_NAME"])
	}
}

func TestMaskSecrets_NoKeys(t *testing.T) {
	secrets := map[string]string{"TOKEN": "secret123"}
	result := MaskSecrets(secrets, nil, DefaultMaskOption)
	if result["TOKEN"] != "secret123" {
		t.Error("expected TOKEN unchanged when no mask keys provided")
	}
}
