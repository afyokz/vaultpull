package vault

import (
	"strings"
)

// MaskOption controls how a secret value is masked.
type MaskOption struct {
	ShowPrefix int // number of leading chars to reveal
	ShowSuffix int // number of trailing chars to reveal
	Replacement string
}

// DefaultMaskOption reveals 2 leading and 2 trailing characters.
var DefaultMaskOption = MaskOption{
	ShowPrefix:  2,
	ShowSuffix:  2,
	Replacement: "****",
}

// MaskValue masks a single secret value according to the given option.
func MaskValue(value string, opt MaskOption) string {
	if opt.Replacement == "" {
		opt.Replacement = "****"
	}
	runes := []rune(value)
	n := len(runes)
	if n == 0 {
		return ""
	}
	prefix := opt.ShowPrefix
	suffix := opt.ShowSuffix
	if prefix+suffix >= n {
		return opt.Replacement
	}
	var sb strings.Builder
	sb.WriteString(string(runes[:prefix]))
	sb.WriteString(opt.Replacement)
	sb.WriteString(string(runes[n-suffix:]))
	return sb.String()
}

// MaskSecrets returns a copy of secrets with values masked.
func MaskSecrets(secrets map[string]string, keys []string, opt MaskOption) map[string]string {
	result := make(map[string]string, len(secrets))
	maskSet := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		maskSet[k] = struct{}{}
	}
	for k, v := range secrets {
		if _, ok := maskSet[k]; ok {
			result[k] = MaskValue(v, opt)
		} else {
			result[k] = v
		}
	}
	return result
}
