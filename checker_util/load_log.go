package checker

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/jesee-kuya/stock_exchange/engine"
)

// LoadLog reads a log file from the specified path and populates the Checker's Log field
// with schedule entries. Each line in the log file is expected to be in the format:
// "<cycle>:<process_name>", where <cycle> is an integer representing the cycle number
// and <process_name> is a string representing the name of the process.
//
// Lines that do not match the expected format or contain invalid cycle numbers are skipped.
//
// Parameters:
//   - path: The file path to the log file.
//
// Returns:
//   - error: An error if the file cannot be opened or if there is an issue during scanning.
//     Returns nil if the log is successfully loaded.
//
// Example log file format:
//
//	1:ProcessA
//	2:ProcessB
//	3:ProcessC
func (c *Checker) LoadLog(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	c.Log = []engine.ScheduleEntry{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip lines that don't match the expected format
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		cycleStr := strings.TrimSpace(parts[0])
		name := strings.TrimSpace(parts[1])

		cycle, err := strconv.Atoi(cycleStr)
		if err != nil {
			continue // skip lines with invalid cycle numbers
		}

		entry := engine.ScheduleEntry{
			Cycle:       cycle,
			ProcessName: name,
		}
		c.Log = append(c.Log, entry)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
