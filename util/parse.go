package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jesee-kuya/stock_exchange/process"
)

// ConfigData represents the parsed configuration data
type ConfigData struct {
	Stocks          map[string]int
	Processes       []*process.Process
	OptimizeTargets []string
}

// LoadConfig parses a configuration file and populates the Checker with
// initial stocks, processes, and optimization goals.
func ParseConfig(path string) (*ConfigData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	config := &ConfigData{
		Stocks:          make(map[string]int),
		Processes:       make([]*process.Process, 0),
		OptimizeTargets: make([]string, 0),
	}

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
		if err := parseLine(config, line, lineNumber); err != nil {
			return nil, fmt.Errorf("error parsing line %d: %w", lineNumber, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	return config, nil
}

// parseLine parses a single line from the configuration file
func parseLine(config *ConfigData, line string, lineNumber int) error {
	// Check if it's a stock definition (name:quantity)
	if !strings.Contains(line, "(") && strings.Contains(line, ":") && !strings.HasPrefix(line, "optimize:") {
		return parseStock(config, line)
	}
	return fmt.Errorf("unrecognized line format: %s", line)
}

// / parseStock parses a stock line in format "name:quantity"
func parseStock(config *ConfigData, line string) error {
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

	config.Stocks[name] = quantity
	return nil
}
