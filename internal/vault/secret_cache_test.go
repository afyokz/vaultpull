package vault

import (
	"testing"
	"time"
)

func TestSecretCache_MissOnEmpty(t *testing.T) {
	c := NewSecretCache(5 * time.Minute)
	_, ok := c.Get("secret/app")
	if ok {
		t.Fatal("expected cache miss on empty cache")
	}
}

func TestSecretCache_HitAfterSet(t *testing.T) {
	c := NewSecretCache(5 * time.Minute)
	secrets := map[string]string{"KEY": "value"}
	c.Set("secret/app", secrets)

	got, ok := c.Get("secret/app")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if got["KEY"] != "value" {
		t.Errorf("expected value %q, got %q", "value", got["KEY"])
	}
}

func TestSecretCache_ReturnsCopy(t *testing.T) {
	c := NewSecretCache(5 * time.Minute)
	c.Set("secret/app", map[string]string{"A": "1"})

	got, _ := c.Get("secret/app")
	got["A"] = "mutated"

	got2, _ := c.Get("secret/app")
	if got2["A"] != "1" {
		t.Error("cache returned a reference instead of a copy")
	}
}

func TestSecretCache_ExpiredEntry(t *testing.T) {
	c := NewSecretCache(10 * time.Millisecond)
	c.Set("secret/app", map[string]string{"X": "y"})
	time.Sleep(20 * time.Millisecond)

	_, ok := c.Get("secret/app")
	if ok {
		t.Fatal("expected cache miss after TTL expiry")
	}
}

func TestSecretCache_Invalidate(t *testing.T) {
	c := NewSecretCache(5 * time.Minute)
	c.Set("secret/app", map[string]string{"K": "v"})
	c.Invalidate("secret/app")

	_, ok := c.Get("secret/app")
	if ok {
		t.Fatal("expected cache miss after invalidation")
	}
}

func TestSecretCache_Flush(t *testing.T) {
	c := NewSecretCache(5 * time.Minute)
	c.Set("secret/a", map[string]string{"A": "1"})
	c.Set("secret/b", map[string]string{"B": "2"})
	c.Flush()

	if _, ok := c.Get("secret/a"); ok {
		t.Error("expected miss for secret/a after flush")
	}
	if _, ok := c.Get("secret/b"); ok {
		t.Error("expected miss for secret/b after flush")
	}
}
