package config

import (
	"testing"
)

func TestBuildRule_Nil(t *testing.T) {
	var tc *TransformConfig
	r, err := tc.BuildRule()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil rule")
	}
}

func TestBuildRule_Prefix(t *testing.T) {
	tc := &TransformConfig{Prefix: "APP_"}
	_, err := tc.BuildRule()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuildRule_MutuallyExclusive(t *testing.T) {
	tc := &TransformConfig{Uppercase: true, Lowercase: true}
	_, err := tc.BuildRule()
	if err == nil {
		t.Fatal("expected error for uppercase+lowercase")
	}
}

func TestBuildRule_Uppercase(t *testing.T) {
	tc := &TransformConfig{Uppercase: true, Prefix: "SVC_"}
	r, err := tc.BuildRule()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := r.Apply(map[string]string{"key": "v"})
	if _, ok := out["SVC_KEY"]; !ok {
		t.Errorf("expected SVC_KEY, got %v", out)
	}
}
