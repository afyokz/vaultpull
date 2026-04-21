package vault

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ImportFormat represents the format of an import source file.
type ImportFormat string

const (
	ImportFormatDotenv ImportFormat = "dotenv"
	ImportFormatJSON   ImportFormat = "json"
	ImportFormatShell  ImportFormat = "shell"
)

// ParseImportFormat parses and validates an import format string.
func ParseImportFormat(s string) (ImportFormat, error) {
	switch ImportFormat(strings.ToLower(s)) {
	case ImportFormatDotenv:
		return ImportFormatDotenv, nil
	case ImportFormatJSON:
		return ImportFormatJSON, nil
	case ImportFormatShell:
		return ImportFormatShell, nil
	default:
		return "", fmt.Errorf("unsupported import format %q: must be dotenv, json, or shell", s)
	}
}

// ImportSecrets reads secrets from a file (or stdin if path is "-") in the
// given format and returns them as a flat key/value map.
func ImportSecrets(path string, format ImportFormat) (map[string]string, error) {
	var r *os.File
	if path == "-" {
		r = os.Stdin
	} else {
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("import: open %q: %w", path, err)
		}
		defer f.Close()
		r = f
	}

	switch format {
	case ImportFormatDotenv:
		return importDotenv(r)
	case ImportFormatJSON:
		return importJSON(r)
	case ImportFormatShell:
		return importShell(r)
	default:
		return nil, fmt.Errorf("import: unknown format %q", format)
	}
}

func importDotenv(f *os.File) (map[string]string, error) {
	out := map[string]string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		if key != "" {
			out[key] = val
		}
	}
	return out, scanner.Err()
}

func importJSON(f *os.File) (map[string]string, error) {
	raw := map[string]interface{}{}
	if err := json.NewDecoder(f).Decode(&raw); err != nil {
		return nil, fmt.Errorf("import: decode json: %w", err)
	}
	out := map[string]string{}
	for k, v := range raw {
		out[k] = fmt.Sprintf("%v", v)
	}
	return out, nil
}

func importShell(f *os.File) (map[string]string, error) {
	out := map[string]string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `'"'`)
		if key != "" {
			out[key] = val
		}
	}
	return out, scanner.Err()
}
