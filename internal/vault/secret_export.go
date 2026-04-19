package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ExportFormat defines the output format for secret export.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatDotenv ExportFormat = "dotenv"
	FormatShell ExportFormat = "shell"
)

// ExportOptions controls how secrets are exported.
type ExportOptions struct {
	Format  ExportFormat
	OutFile string // empty means stdout
	Redact  bool
}

// ExportSecrets serializes secrets to the given format and writes to file or stdout.
func ExportSecrets(secrets map[string]string, opts ExportOptions) error {
	var out string
	var err error

	switch opts.Format {
	case FormatJSON:
		out, err = exportJSON(secrets, opts.Redact)
	case FormatShell:
		out = exportShell(secrets, opts.Redact)
	case FormatDotenv, "":
		out = exportDotenv(secrets, opts.Redact)
	default:
		return fmt.Errorf("unsupported export format: %s", opts.Format)
	}
	if err != nil {
		return err
	}

	if opts.OutFile == "" {
		fmt.Print(out)
		return nil
	}
	return os.WriteFile(opts.OutFile, []byte(out), 0600)
}

func maybeRedact(v string, redact bool) string {
	if redact {
		return "***"
	}
	return v
}

func exportJSON(secrets map[string]string, redact bool) (string, error) {
	m := make(map[string]string, len(secrets))
	for k, v := range secrets {
		m[k] = maybeRedact(v, redact)
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b) + "\n", nil
}

func exportDotenv(secrets map[string]string, redact bool) string {
	var sb strings.Builder
	for k, v := range secrets {
		fmt.Fprintf(&sb, "%s=%s\n", k, maybeRedact(v, redact))
	}
	return sb.String()
}

func exportShell(secrets map[string]string, redact bool) string {
	var sb strings.Builder
	for k, v := range secrets {
		fmt.Fprintf(&sb, "export %s=%q\n", k, maybeRedact(v, redact))
	}
	return sb.String()
}
