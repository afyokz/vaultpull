package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter_NoRules(t *testing.T) {
	f := &Filter{}
	input := map[string]string{"DB_HOST": "localhost", "API_KEY": "abc"}
	out := f.Apply(input)
	assert.Equal(t, input, out)
}

func TestFilter_IncludeExact(t *testing.T) {
	f := &Filter{Include: []string{"DB_HOST"}}
	input := map[string]string{"DB_HOST": "localhost", "API_KEY": "abc"}
	out := f.Apply(input)
	assert.Equal(t, map[string]string{"DB_HOST": "localhost"}, out)
}

func TestFilter_IncludeWildcard(t *testing.T) {
	f := &Filter{Include: []string{"DB_*"}}
	input := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432", "API_KEY": "abc"}
	out := f.Apply(input)
	assert.Equal(t, map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}, out)
}

func TestFilter_ExcludeExact(t *testing.T) {
	f := &Filter{Exclude: []string{"API_KEY"}}
	input := map[string]string{"DB_HOST": "localhost", "API_KEY": "abc"}
	out := f.Apply(input)
	assert.Equal(t, map[string]string{"DB_HOST": "localhost"}, out)
}

func TestFilter_ExcludeWildcard(t *testing.T) {
	f := &Filter{Exclude: []string{"SECRET_*"}}
	input := map[string]string{"SECRET_TOKEN": "x", "SECRET_KEY": "y", "DB_HOST": "localhost"}
	out := f.Apply(input)
	assert.Equal(t, map[string]string{"DB_HOST": "localhost"}, out)
}

func TestFilter_IncludeAndExclude(t *testing.T) {
	f := &Filter{Include: []string{"DB_*"}, Exclude: []string{"DB_PASSWORD"}}
	input := map[string]string{"DB_HOST": "localhost", "DB_PASSWORD": "secret", "API_KEY": "abc"}
	out := f.Apply(input)
	assert.Equal(t, map[string]string{"DB_HOST": "localhost"}, out)
}
