package checker

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// LoadConfig parses a configuration file and populates the Checker with
// initial stocks, processes, and optimization goals.
func (c *Checker) LoadConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse the line based on its format
		if err := c.parseLine(line, lineNumber); err != nil {
			return fmt.Errorf("error parsing line %d: %w", lineNumber, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}
	return nil
}

// parseLine parses a single line from the configuration file
func (c *Checker) parseLine(line string, lineNumber int) error {
	// Check if it's a stock definition (name:quantity)
	if !strings.Contains(line, "(") && strings.Contains(line, ":") && !strings.HasPrefix(line, "optimize:") {
		return c.parseStock(line)
	}
	return fmt.Errorf("unrecognized line format: %s", line)
}

// parseStock parses a stock line in format "name:quantity"
func (c *Checker) parseStock(line string) error {
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid stock format: %s", line)
	}

	name := strings.TrimSpace(parts[0])
	quantityStr := strings.TrimSpace(parts[1])

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		return fmt.Errorf("invalid stock quantity '%s': %w", quantityStr, err)
	}

	c.Stocks[name] = quantity
	return nil
}
