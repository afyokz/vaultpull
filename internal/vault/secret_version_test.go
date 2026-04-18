package vault

import (
	"testing"
)

func TestParseVersionedPath_NoVersion(t *testing.T) {
	vp, err := ParseVersionedPath("secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vp.Version != 0 {
		t.Errorf("expected version 0, got %d", vp.Version)
	}
	if vp.String() != "secret/data/myapp" {
		t.Errorf("unexpected string: %s", vp.String())
	}
}

func TestParseVersionedPath_WithVersion(t *testing.T) {
	vp, err := ParseVersionedPath("secret/data/myapp@3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vp.Version != 3 {
		t.Errorf("expected version 3, got %d", vp.Version)
	}
	if vp.String() != "secret/data/myapp@3" {
		t.Errorf("unexpected string: %s", vp.String())
	}
}

func TestParseVersionedPath_InvalidVersion(t *testing.T) {
	cases := []string{
		"secret/data/myapp@0",
		"secret/data/myapp@abc",
		"secret/data/myapp@-1",
	}
	for _, c := range cases {
		_, err := ParseVersionedPath(c)
		if err == nil {
			t.Errorf("expected error for %q, got nil", c)
		}
	}
}

func TestVersionParams_Latest(t *testing.T) {
	vp, _ := ParseVersionedPath("secret/data/myapp")
	p := vp.VersionParams()
	if len(p) != 0 {
		t.Errorf("expected empty params for latest, got %v", p)
	}
}

func TestVersionParams_Specific(t *testing.T) {
	vp, _ := ParseVersionedPath("secret/data/myapp@2")
	p := vp.VersionParams()
	if p["version"] != "2" {
		t.Errorf("expected version=2, got %v", p)
	}
}
