package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry represents a single audit log entry.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	Path      string    `json:"path"`
	Keys      []string  `json:"keys"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
}

// Logger writes audit entries to a file.
type Logger struct {
	filePath string
}

// NewLogger creates a new Logger that appends to filePath.
func NewLogger(filePath string) *Logger {
	return &Logger{filePath: filePath}
}

// Log writes an audit entry to the log file.
func (l *Logger) Log(op, path string, keys []string, err error) error {
	e := Entry{
		Timestamp: time.Now().UTC(),
		Operation: op,
		Path:      path,
		Keys:      keys,
		Success:   err == nil,
	}
	if err != nil {
		e.Error = err.Error()
	}

	f, openErr := os.OpenFile(l.filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if openErr != nil {
		return fmt.Errorf("audit: open log file: %w", openErr)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if encErr := enc.Encode(e); encErr != nil {
		return fmt.Errorf("audit: encode entry: %w", encErr)
	}
	return nil
}
