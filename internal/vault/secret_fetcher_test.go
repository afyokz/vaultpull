package vault

import (
	"testing"
)

func TestMergeSecrets_Empty(t *testing.T) {
	result := MergeSecrets([]FetchedSecret{})
	if len(result) != 0 {
		t.Errorf("expected empty map, got %d entries", len(result))
	}
}

func TestMergeSecrets_Single(t *testing.T) {
	secrets := []FetchedSecret{
		{
			Path: SecretPath{Mount: "secret", Path: "app/db"},
			Data: map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"},
		},
	}

	result := MergeSecrets(secrets)

	if result["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %s", result["DB_HOST"])
	}
	if result["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %s", result["DB_PORT"])
	}
}

func TestMergeSecrets_OverwritesOnConflict(t *testing.T) {
	secrets := []FetchedSecret{
		{
			Path: SecretPath{Mount: "secret", Path: "app/base"},
			Data: map[string]string{"API_KEY": "old-key", "TIMEOUT": "30"},
		},
		{
			Path: SecretPath{Mount: "secret", Path: "app/override"},
			Data: map[string]string{"API_KEY": "new-key"},
		},
	}

	result := MergeSecrets(secrets)

	if result["API_KEY"] != "new-key" {
		t.Errorf("expected API_KEY=new-key after merge, got %s", result["API_KEY"])
	}
	if result["TIMEOUT"] != "30" {
		t.Errorf("expected TIMEOUT=30 to be preserved, got %s", result["TIMEOUT"])
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}
