package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jesee-kuya/stock_exchange/process"
)

// ConfigData holds the parsed configuration data for the stock exchange system.
// It includes the initial stock quantities, the list of processes, and the optimization targets.
//
// Fields:
//   - Stocks: a map where the key is the stock name and the value is the initial quantity.
//   - Processes: a slice of pointers to Process structs, each representing a process definition.
//   - OptimizeTargets: a slice of strings specifying the optimization goals extracted from the config file.
type ConfigData struct {
	Stocks          map[string]int
	Processes       []*process.Process
	OptimizeTargets []string
}

// ParseConfig reads a configuration file from the specified path and parses its contents
// into a ConfigData struct. The configuration file is expected to define initial stock
// quantities, process definitions, and optimization targets. Each line in the file is
// interpreted based on its format:
//   - Stock definitions: "name:quantity"
//   - Process definitions: "name:(needs):(results):cycles"
//   - Optimization targets: "optimize:(target1;target2;...)"
// Lines that are empty or start with '#' are ignored as comments.
// Returns a pointer to the populated ConfigData struct or an error if parsing fails.
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

// parseLine analyzes a single line from the configuration file and updates the provided
// ConfigData struct accordingly. It determines the type of configuration entry based on
// the line's format and delegates parsing to the appropriate helper function:
//   - Stock definitions (e.g., "name:quantity") are handled by parseStock.
//   - Process definitions (e.g., "name:(needs):(results):cycles") are handled by parseProcess.
//   - Optimization targets (e.g., "optimize:(target1;target2;...)") are handled by parseOptimize.
// Returns an error if the line format is unrecognized or if parsing fails.
func parseLine(config *ConfigData, line string) error {
	// Check if it's a stock definition (name:quantity)
	if !strings.Contains(line, "(") && strings.Contains(line, ":") && !strings.HasPrefix(line, "optimize:") {
		return parseStock(config, line)
	}

	// Check if it's an optimize line
	if strings.HasPrefix(line, "optimize:") {
		return parseOptimize(config, line)
	}

	// Check if it's a process definition (contains parentheses)
	if strings.Contains(line, "(") && strings.Contains(line, ")") {
		return parseProcess(config, line)
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

// parseProcess parses a process line in format "name:(needs):(results):cycles"
func parseProcess(config *ConfigData, line string) error {
	// Find the first colon to separate name from the rest
	colonIndex := strings.Index(line, ":")
	if colonIndex == -1 {
		return fmt.Errorf("invalid process format: %s", line)
	}

	name := strings.TrimSpace(line[:colonIndex])
	rest := line[colonIndex+1:]

	// Parse the remaining parts: (needs):(results):cycles
	parts := strings.Split(rest, "):")
	if len(parts) != 3 {
		return fmt.Errorf("invalid process format: %s", line)
	}

	// Parse needs
	needs, err := parseResourceMap(parts[0])
	if err != nil {
		return fmt.Errorf("error parsing needs: %w", err)
	}

	// Parse results
	results, err := parseResourceMap(parts[1])
	if err != nil {
		return fmt.Errorf("error parsing results: %w", err)
	}

	// Parse Cycles
	cycles, err := strconv.Atoi(strings.TrimSpace(parts[2]))
	if err != nil {
		return fmt.Errorf("invalid cycle count '%s': %w", parts[2], err)
	}

	proc := &process.Process{
		Name:   name,
		Needs:  needs,
		Result: results,
		Cycle:  cycles,
	}

	config.Processes = append(config.Processes, proc)
	return nil
}

func parseResourceMap(blockStr string) (map[string]int, error) {
	resources := make(map[string]int)
	blockStr = strings.Trim(blockStr, "()")
	if blockStr == "" {
		return resources, nil
	}

	items := strings.Split(blockStr, ";")
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}

		parts := strings.Split(item, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid resource format: %s", item)
		}
		name := strings.TrimSpace(parts[0])
		quantity, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid resource quantity '%v': %w", quantity, err)
		}
		resources[name] = quantity
	}
	return resources, nil
}
