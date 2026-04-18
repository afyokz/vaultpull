package vault

import (
	"context"
	"fmt"
)

// FetchedSecret holds the resolved path and its flattened key-value pairs.
type FetchedSecret struct {
	Path SecretPath
	Data map[string]string
}

// FetchSecrets retrieves secrets for all provided paths using the given client.
func FetchSecrets(ctx context.Context, client *Client, paths []SecretPath) ([]FetchedSecret, error) {
	results := make([]FetchedSecret, 0, len(paths))

	for _, p := range paths {
		data, err := client.GetSecret(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("fetching secret %s: %w", p, err)
		}

		results = append(results, FetchedSecret{
			Path: p,
			Data: data,
		})
	}

	return results, nil
}

// MergeSecrets merges all fetched secrets into a single map.
// Later entries overwrite earlier ones on key conflict.
func MergeSecrets(secrets []FetchedSecret) map[string]string {
	merged := make(map[string]string)
	for _, s := range secrets {
		for k, v := range s.Data {
			merged[k] = v
		}
	}
	return merged
}
