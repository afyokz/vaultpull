package vault

import (
	"fmt"
	"strings"
)

// TransformRule defines a key transformation to apply to fetched secrets.
type TransformRule struct {
	Prefix string
	Suffix string
	Uppercase bool
	Lowercase bool
}

// TransformOption configures a TransformRule.
type TransformOption func(*TransformRule)

func WithPrefix(p string) TransformOption { return func(r *TransformRule) { r.Prefix = p } }
func WithSuffix(s string) TransformOption  { return func(r *TransformRule) { r.Suffix = s } }
func WithUppercase() TransformOption       { return func(r *TransformRule) { r.Uppercase = true } }
func WithLowercase() TransformOption       { return func(r *TransformRule) { r.Lowercase = true } }

// NewTransformRule builds a TransformRule from options.
func NewTransformRule(opts ...TransformOption) (*TransformRule, error) {
	r := &TransformRule{}
	for _, o := range opts {
		o(r)
	}
	if r.Uppercase && r.Lowercase {
		return nil, fmt.Errorf("transform: uppercase and lowercase are mutually exclusive")
	}
	return r, nil
}

// Apply transforms a map of secrets according to the rule.
func (r *TransformRule) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		nk := k
		if r.Uppercase {
			nk = strings.ToUpper(nk)
		} else if r.Lowercase {
			nk = strings.ToLower(nk)
		}
		if r.Prefix != "" {
			nk = r.Prefix + nk
		}
		if r.Suffix != "" {
			nk = nk + r.Suffix
		}
		out[nk] = v
	}
	return out
}
