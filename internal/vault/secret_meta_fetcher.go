package vault

import (
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// FetchSecretMeta retrieves metadata for a KV v2 secret path.
func FetchSecretMeta(client *vaultapi.Client, mount, subPath string) (SecretMeta, error) {
	metaPath := fmt.Sprintf("%s/metadata/%s", mount, subPath)
	secret, err := client.Logical().Read(metaPath)
	if err != nil {
		return SecretMeta{}, fmt.Errorf("reading metadata at %s: %w", metaPath, err)
	}
	if secret == nil || secret.Data == nil {
		return SecretMeta{}, fmt.Errorf("no metadata found at %s", metaPath)
	}

	versions, ok := secret.Data["versions"].(map[string]interface{})
	if !ok || len(versions) == 0 {
		return SecretMeta{}, fmt.Errorf("no versions in metadata at %s", metaPath)
	}

	// Find the latest version number.
	latestRaw, _ := secret.Data["current_version"].(float64)
	latestKey := fmt.Sprintf("%.0f", latestRaw)

	var createdTime time.Time
	if vData, ok := versions[latestKey].(map[string]interface{}); ok {
		if ts, ok := vData["created_time"].(string); ok {
			createdTime, _ = time.Parse(time.RFC3339Nano, ts)
		}
	}

	return SecretMeta{
		Path:        fmt.Sprintf("%s/%s", mount, subPath),
		Version:     int(latestRaw),
		CreatedTime: createdTime,
	}, nil
}
