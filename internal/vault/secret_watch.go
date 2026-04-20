package vault

import (
	"context"
	"fmt"
	"time"
)

// WatchOption configures the watcher behavior.
type WatchOption struct {
	Interval  time.Duration
	MaxErrors int
}

// DefaultWatchOption returns sensible defaults for watching.
func DefaultWatchOption() WatchOption {
	return WatchOption{
		Interval:  30 * time.Second,
		MaxErrors: 3,
	}
}

// WatchEvent represents a detected change in secrets.
type WatchEvent struct {
	Path    string
	Old     map[string]string
	New     map[string]string
	Err     error
	CheckAt time.Time
}

// HasChanges returns true if the event contains differing secrets.
func (e WatchEvent) HasChanges() bool {
	if len(e.Old) != len(e.New) {
		return true
	}
	for k, v := range e.New {
		if e.Old[k] != v {
			return true
		}
	}
	return false
}

// Watcher polls a secret path and emits events on change.
type Watcher struct {
	fetch  func(ctx context.Context, path string) (map[string]string, error)
	option WatchOption
}

// NewWatcher creates a Watcher using the provided fetch function.
func NewWatcher(fetch func(ctx context.Context, path string) (map[string]string, error), opt WatchOption) *Watcher {
	return &Watcher{fetch: fetch, option: opt}
}

// Watch starts polling the given path and sends events to the returned channel.
// The channel is closed when ctx is cancelled.
func (w *Watcher) Watch(ctx context.Context, path string) <-chan WatchEvent {
	ch := make(chan WatchEvent, 1)
	go func() {
		defer close(ch)
		var prev map[string]string
		errCount := 0
		for {
			current, err := w.fetch(ctx, path)
			event := WatchEvent{Path: path, Old: prev, New: current, Err: err, CheckAt: time.Now()}
			if err != nil {
				errCount++
				if errCount >= w.option.MaxErrors {
					event.Err = fmt.Errorf("max errors reached (%d): %w", w.option.MaxErrors, err)
					ch <- event
					return
				}
				ch <- event
			} else {
				errCount = 0
				if event.HasChanges() {
					ch <- event
				}
				prev = current
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(w.option.Interval):
			}
		}
	}()
	return ch
}
