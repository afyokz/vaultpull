package vault

import (
	"sync"
	"time"
)

// CacheEntry holds a cached secret map with an expiry timestamp.
type CacheEntry struct {
	Secrets   map[string]string
	FetchedAt time.Time
	TTL       time.Duration
}

// IsExpired returns true if the cache entry is past its TTL.
func (e *CacheEntry) IsExpired() bool {
	return time.Since(e.FetchedAt) > e.TTL
}

// SecretCache is a simple in-memory TTL cache keyed by secret path.
type SecretCache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
	ttl     time.Duration
}

// NewSecretCache creates a SecretCache with the given TTL.
func NewSecretCache(ttl time.Duration) *SecretCache {
	return &SecretCache{
		entries: make(map[string]*CacheEntry),
		ttl:     ttl,
	}
}

// Get returns cached secrets for a path, or (nil, false) if absent/expired.
func (c *SecretCache) Get(path string) (map[string]string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[path]
	if !ok || e.IsExpired() {
		return nil, false
	}
	copy := make(map[string]string, len(e.Secrets))
	for k, v := range e.Secrets {
		copy[k] = v
	}
	return copy, true
}

// Set stores secrets for a path, resetting the TTL.
func (c *SecretCache) Set(path string, secrets map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}
	c.entries[path] = &CacheEntry{
		Secrets:   copy,
		FetchedAt: time.Now(),
		TTL:       c.ttl,
	}
}

// Invalidate removes a single path from the cache.
func (c *SecretCache) Invalidate(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, path)
}

// Flush clears all entries from the cache.
func (c *SecretCache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*CacheEntry)
}
