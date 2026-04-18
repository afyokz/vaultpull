package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRender_NoPlaceholders(t *testing.T) {
	r := NewRenderer()
	out, missing, err := r.Render("hello world", map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "hello world" {
		t.Errorf("expected 'hello world', got %q", out)
	}
	if len(missing) != 0 {
		t.Errorf("expected no missing keys, got %v", missing)
	}
}

func TestRender_Substitutes(t *testing.T) {
	r := NewRenderer()
	src := "DB_HOST={{DB_HOST}} DB_PORT={{DB_PORT}}"
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	out, missing, err := r.Render(src, secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "DB_HOST=localhost DB_PORT=5432" {
		t.Errorf("unexpected output: %q", out)
	}
	if len(missing) != 0 {
		t.Errorf("expected no missing, got %v", missing)
	}
}

func TestRender_MissingKeys(t *testing.T) {
	r := NewRenderer()
	src := "{{FOO}} and {{BAR}} and {{FOO}}"
	out, missing, err := r.Render(src, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != src {
		t.Errorf("expected unchanged src, got %q", out)
	}
	if len(missing) != 2 {
		t.Errorf("expected 2 unique missing keys, got %v", missing)
	}
}

func TestRenderFile_ReadsAndRenders(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "tmpl.env")
	_ = os.WriteFile(p, []byte("SECRET={{MY_SECRET}}"), 0644)

	r := NewRenderer()
	out, missing, err := r.RenderFile(p, map[string]string{"MY_SECRET": "abc123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "SECRET=abc123" {
		t.Errorf("unexpected output: %q", out)
	}
	if len(missing) != 0 {
		t.Errorf("expected no missing, got %v", missing)
	}
}

func TestRenderFile_NotFound(t *testing.T) {
	r := NewRenderer()
	_, _, err := r.RenderFile("/nonexistent/path.env", nil)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
