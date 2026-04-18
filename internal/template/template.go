package template

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Renderer renders a template file by substituting {{VAR}} placeholders
// with values from a provided secrets map.
type Renderer struct {
	placeholder *regexp.Regexp
}

// NewRenderer creates a new Renderer.
func NewRenderer() *Renderer {
	return &Renderer{
		placeholder: regexp.MustCompile(`\{\{([A-Z0-9_]+)\}\}`),
	}
}

// RenderFile reads a template file and substitutes placeholders with secrets.
// Returns the rendered string and a list of any keys that were not found.
func (r *Renderer) RenderFile(path string, secrets map[string]string) (string, []string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", nil, fmt.Errorf("read template: %w", err)
	}
	return r.Render(string(data), secrets)
}

// Render substitutes placeholders in src with values from secrets.
// Returns rendered string and missing keys.
func (r *Renderer) Render(src string, secrets map[string]string) (string, []string, error) {
	var missing []string
	seen := map[string]bool{}

	result := r.placeholder.ReplaceAllStringFunc(src, func(match string) string {
		key := strings.TrimSuffix(strings.TrimPrefix(match, "{{"), "}}")
		val, ok := secrets[key]
		if !ok {
			if !seen[key] {
				missing = append(missing, key)
				seen[key] = true
			}
			return match
		}
		return val
	})

	return result, missing, nil
}
