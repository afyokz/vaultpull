package prompt

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/vaultpull/internal/diff"
)

// Confirmer handles interactive user confirmation before applying changes.
type Confirmer struct {
	in  io.Reader
	out io.Writer
}

// NewConfirmer creates a Confirmer reading from stdin and writing to stdout.
func NewConfirmer() *Confirmer {
	return &Confirmer{in: os.Stdin, out: os.Stdout}
}

// NewConfirmerWithIO creates a Confirmer with custom IO (useful for tests).
func NewConfirmerWithIO(in io.Reader, out io.Writer) *Confirmer {
	return &Confirmer{in: in, out: out}
}

// ConfirmDiff prints a summary of the diff and asks the user to confirm.
// Returns true if the user confirms, false otherwise.
func (c *Confirmer) ConfirmDiff(result diff.Result) (bool, error) {
	if len(result.Changes) == 0 {
		fmt.Fprintln(c.out, "No changes detected.")
		return false, nil
	}

	fmt.Fprintln(c.out, "\nPending changes:")
	for _, ch := range result.Changes {
		switch ch.Type {
		case diff.Added:
			fmt.Fprintf(c.out, "  + %s\n", ch.Key)
		case diff.Modified:
			fmt.Fprintf(c.out, "  ~ %s\n", ch.Key)
		case diff.Removed:
			fmt.Fprintf(c.out, "  - %s\n", ch.Key)
		}
	}
	fmt.Fprintf(c.out, "\n%s\n", result.Summary())
	fmt.Fprint(c.out, "Apply changes? [y/N]: ")

	var response string
	_, err := fmt.Fscanln(c.in, &response)
	if err != nil {
		// treat EOF / empty input as no
		return false, nil
	}

	return strings.ToLower(strings.TrimSpace(response)) == "y", nil
}
