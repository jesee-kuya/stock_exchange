package checker

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/jesee-kuya/stock_exchange/engine"
)

// LoadLog reads the log file and stores the sequence of executed processes as ScheduleEntry.
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
			Cycle: cycle,
			Name:  name,
		}
		c.Log = append(c.Log, entry)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
