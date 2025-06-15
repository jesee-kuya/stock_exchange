package engine

import (
	"os"
	"strings"
)

// SaveLog persists the simulation log to a specified file.
// Each line of the log follows the format: <cycle>:<process_name>
// It tracks the exact order and timing of process executions for future analysis.
func (e *Engine) SaveLog(path string) error {
	content := strings.Join(e.Schedule, "\n")
	return os.WriteFile(path, []byte(content), 0o644)
}
