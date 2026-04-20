package vault

import (
	"context"
	"errors"
	"testing"
	"time"
)

func fastOpt() WatchOption {
	return WatchOption{Interval: 10 * time.Millisecond, MaxErrors: 3}
}

func TestWatchEvent_HasChanges_NoChange(t *testing.T) {
	e := WatchEvent{
		Old: map[string]string{"K": "v"},
		New: map[string]string{"K": "v"},
	}
	if e.HasChanges() {
		t.Fatal("expected no changes")
	}
}

func TestWatchEvent_HasChanges_Modified(t *testing.T) {
	e := WatchEvent{
		Old: map[string]string{"K": "old"},
		New: map[string]string{"K": "new"},
	}
	if !e.HasChanges() {
		t.Fatal("expected changes")
	}
}

func TestWatchEvent_HasChanges_Added(t *testing.T) {
	e := WatchEvent{
		Old: map[string]string{},
		New: map[string]string{"K": "v"},
	}
	if !e.HasChanges() {
		t.Fatal("expected changes due to added key")
	}
}

func TestWatch_EmitsOnChange(t *testing.T) {
	calls := 0
	data := []map[string]string{
		{"A": "1"},
		{"A": "2"},
	}
	fetch := func(ctx context.Context, path string) (map[string]string, error) {
		idx := calls
		if idx >= len(data) {
			idx = len(data) - 1
		}
		calls++
		return data[idx], nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	w := NewWatcher(fetch, fastOpt())
	ch := w.Watch(ctx, "secret/app")
	var events []WatchEvent
	for e := range ch {
		events = append(events, e)
	}
	if len(events) == 0 {
		t.Fatal("expected at least one change event")
	}
	if events[0].Path != "secret/app" {
		t.Errorf("unexpected path: %s", events[0].Path)
	}
}

func TestWatch_StopsOnMaxErrors(t *testing.T) {
	fetchErr := errors.New("vault unavailable")
	fetch := func(ctx context.Context, path string) (map[string]string, error) {
		return nil, fetchErr
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	opt := WatchOption{Interval: 5 * time.Millisecond, MaxErrors: 2}
	w := NewWatcher(fetch, opt)
	ch := w.Watch(ctx, "secret/app")
	var last WatchEvent
	for e := range ch {
		last = e
	}
	if last.Err == nil {
		t.Fatal("expected terminal error event")
	}
}
