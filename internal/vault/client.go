package vault

import (
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client.
type Client struct {
	api *vaultapi.Client
}

// NewClient creates a new Vault client with the given address and token.
func NewClient(addr, token string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = addr

	api, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	api.SetToken(token)

	return &Client{api: api}, nil
}

// GetSecrets fetches key-value secrets from the given KV v2 path.
func (c *Client) GetSecrets(path string) (map[string]string, error) {
	secret, err := c.api.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("vault read error: %w", err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at path: %s", path)
	}

	// KV v2 nests data under secret.Data["data"]
	raw, ok := secret.Data["data"]
	if !ok {
		// KV v1 fallback
		return flattenData(secret.Data)
	}

	nestedMap, ok := raw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected secret data format at path: %s", path)
	}
	return flattenData(nestedMap)
}

func flattenData(data map[string]interface{}) (map[string]string, error) {
	result := make(map[string]string, len(data))
	for k, v := range data {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result, nil
}
