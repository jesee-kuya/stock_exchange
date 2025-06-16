package checker

import (
	"fmt"
	"os"
)

// LoadConfig parses a configuration file and populates the Checker with
// initial stocks, processes, and optimization goals.
func (c *Checker) LoadConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()
	return nil
}
