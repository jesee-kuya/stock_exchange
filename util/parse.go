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
		if err := parseLine(config, line); err != nil {
			return nil, fmt.Errorf("error parsing line %d: %w", lineNumber, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	return config, nil
}

// parseLine parses a single line from the configuration file
func parseLine(config *ConfigData, line string) error {
	// Check if it's a stock definition (name:quantity)
	if !strings.Contains(line, "(") && strings.Contains(line, ":") && !strings.HasPrefix(line, "optimize:") {
		return parseStock(config, line)
	}

	// Check if it's an optimize line
	if strings.HasPrefix(line, "optimize:") {
		return parseOptimize(config, line)
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

// parseOptimize parses the optimize line and extracts target list
func parseOptimize(config *ConfigData, line string) error {
	// Remove "optimize:" prefix and parse the targets
	targetsPart := strings.TrimPrefix(line, "optimize:")
	targetsPart = strings.Trim(targetsPart, "()")

	if targetsPart == "" {
		return nil
	}

	// Split by semicolon to get individual targets
	targets := strings.Split(targetsPart, ";")
	for _, target := range targets {
		target = strings.TrimSpace(target)
		if target != "" {
			config.OptimizeTargets = append(config.OptimizeTargets, target)
		}
	}

	return nil
}
