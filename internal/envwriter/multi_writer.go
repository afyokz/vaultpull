package envwriter

import (
	"fmt"
	"path/filepath"
)

// GroupedWrite writes each group of secrets to a separate .env file
// under the given directory. The filename is derived from the group name.
func GroupedWrite(dir string, groups map[string]map[string]string, perm uint32) error {
	for group, secrets := range groups {
		filename := fmt.Sprintf(".env.%s", group)
		path := filepath.Join(dir, filename)
		if err := Write(path, secrets, perm); err != nil {
			return fmt.Errorf("writing group %q to %s: %w", group, path, err)
		}
	}
	return nil
}
