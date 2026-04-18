package prompt

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpull/internal/diff"
)

func makeResult(changes []diff.Change) diff.Result {
	return diff.Result{Changes: changes}
}

func TestConfirmDiff_NoChanges(t *testing.T) {
	c := NewConfirmerWithIO(strings.NewReader(""), &strings.Builder{})
	ok, err := c.ConfirmDiff(makeResult(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected false when no changes")
	}
}

func TestConfirmDiff_UserConfirms(t *testing.T) {
	changes := []diff.Change{
		{Key: "DB_HOST", Type: diff.Added},
	}
	c := NewConfirmerWithIO(strings.NewReader("y\n"), &strings.Builder{})
	ok, err := c.ConfirmDiff(makeResult(changes))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected true when user enters 'y'")
	}
}

func TestConfirmDiff_UserDeclines(t *testing.T) {
	changes := []diff.Change{
		{Key: "DB_HOST", Type: diff.Added},
	}
	c := NewConfirmerWithIO(strings.NewReader("n\n"), &strings.Builder{})
	ok, err := c.ConfirmDiff(makeResult(changes))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected false when user enters 'n'")
	}
}

func TestConfirmDiff_EmptyInput(t *testing.T) {
	changes := []diff.Change{
		{Key: "API_KEY", Type: diff.Modified},
	}
	c := NewConfirmerWithIO(strings.NewReader(""), &strings.Builder{})
	ok, err := c.ConfirmDiff(makeResult(changes))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected false on empty input")
	}
}

func TestConfirmDiff_OutputContainsKeys(t *testing.T) {
	changes := []diff.Change{
		{Key: "SECRET_TOKEN", Type: diff.Added},
		{Key: "OLD_KEY", Type: diff.Removed},
	}
	out := &strings.Builder{}
	c := NewConfirmerWithIO(strings.NewReader("n\n"), out)
	_, _ = c.ConfirmDiff(makeResult(changes))
	output := out.String()
	if !strings.Contains(output, "SECRET_TOKEN") {
		t.Error("expected output to contain SECRET_TOKEN")
	}
	if !strings.Contains(output, "OLD_KEY") {
		t.Error("expected output to contain OLD_KEY")
	}
}
