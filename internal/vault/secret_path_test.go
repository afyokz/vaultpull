package vault

import (
	"testing"
)

func TestParseSecretPath_Valid(t *testing.T) {
	cases := []struct {
		input    string
		Mount    string
		SubPath  string
		fullPath string
	}{
		{"secret/myapp", "secret", "data/myapp", "secret/data/myapp"},
		{"secret/data/myapp", "secret", "data/myapp", "secret/data/myapp"},
		{"kvv2/data/team/service", "kvv2", "data/team/service", "kvv2/data/team/service"},
		{"kvv2/team/service", "kvv2", "data/team/service", "kvv2/data/team/service"},
		{"secret/metadata/myapp", "secret", "metadata/myapp", "secret/metadata/myapp"},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			sp, err := ParseSecretPath(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if sp.Mount != tc.Mount {
				t.Errorf("Mount: got %q, want %q", sp.Mount, tc.Mount)
			}
			if sp.SubPath != tc.SubPath {
				t.Errorf("SubPath: got %q, want %q", sp.SubPath, tc.SubPath)
			}
			if sp.FullPath() != tc.fullPath {
				t.Errorf("FullPath: got %q, want %q", sp.FullPath(), tc.fullPath)
			}
		})
	}
}

func TestParseSecretPath_Invalid(t *testing.T) {
	cases := []string{
		"",
		"   ",
		"noseparator",
		"mount/",
	}

	for _, tc := range cases {
		t.Run(tc, func(t *testing.T) {
			_, err := ParseSecretPath(tc)
			if err == nil {
				t.Fatalf("expected error for input %q, got nil", tc)
			}
		})
	}
}

func TestSecretPath_String(t *testing.T) {
	sp := &SecretPath{Mount: "secret", SubPath: "data/myapp"}
	if sp.String() != "secret/data/myapp" {
		t.Errorf("String: got %q, want %q", sp.String(), "secret/data/myapp")
	}
}
