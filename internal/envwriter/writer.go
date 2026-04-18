package envwriter

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/user/vaultpull/internal/backup"
)

// Write merges secrets into the .env file at path, backing up first.
// Existing keys not present in secrets are preserved.
func Write(path string, secrets map[string]string) error {
	bm := backup.NewManager(5)
	if _, err := bm.Create(path); err != nil {
		return fmt.Errorf("envwriter: backup: %w", err)
	}

	existing, err := readEnvFile(path)
	if err != nil {
		return fmt.Errorf("envwriter: read: %w", err)
	}

	for k, v := range secrets {
		existing[k] = v
	}

	return writeEnvFile(path, existing)
}

// readEnvFile parses a .env file into a map. Returns empty map if file absent.
func readEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	defer f.Close()

	result := map[string]string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return result, scanner.Err()
}

// writeEnvFile serialises a map to a .env file with 0600 permissions.
func writeEnvFile(path string, data map[string]string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for k, v := range data {
		fmt.Fprintf(w, "%s=%s\n", k, v)
	}
	return w.Flush()
}
