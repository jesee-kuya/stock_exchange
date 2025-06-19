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
//   - HasOptimizer: a boolean flag to track if an optimizer has already been defined.
type ConfigData struct {
	Stocks          map[string]int
	Processes       []*process.Process
	OptimizeTargets []string
	HasOptimizer    bool
}

// ParseConfig reads a configuration file from the specified path and parses its contents
// into a ConfigData struct. The configuration file is expected to define initial stock
// quantities, process definitions, and optimization targets. Each line in the file is
// interpreted based on its format:
//   - Stock definitions: "name:quantity"
//   - Process definitions: "name:(needs):(results):cycles"
//   - Optimization targets: "optimize:(target1;target2;...)"
//
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
		HasOptimizer:    false,
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
			return nil, fmt.Errorf(" Error while parsing `%s`\nExiting... ", line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf(" Error reading config file: %w\nExiting... ", err)
	}
	return config, nil
}

// parseLine analyzes a single line from the configuration file and updates the provided
// ConfigData struct accordingly. It determines the type of configuration entry based on
// the line's format and delegates parsing to the appropriate helper function:
//   - Stock definitions (e.g., "name:quantity") are handled by parseStock.
//   - Process definitions (e.g., "name:(needs):(results):cycles") are handled by parseProcess.
//   - Optimization targets (e.g., "optimize:(target1;target2;...)") are handled by parseOptimize.
//
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

// parseStock parses a single line of stock data and updates the provided ConfigData.
//
// The expected format for the line is: "item_name:quantity" (e.g., "euro:10").
// It splits the line at the colon, trims whitespace, and converts the quantity to an integer.
// If parsing succeeds, the item and its quantity are added to the config's Stocks map.
//
// Parameters:
//   - config: a pointer to the ConfigData struct to be updated.
//   - line: a string representing one line of stock information.
//
// Returns:
//   - An error if the line format is invalid or the quantity is not a valid integer.
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

// parseOptimize parses an optimization target line and updates the provided ConfigData.
//
// The expected format for the line is: "optimize:(target1;target2;...;targetN)".
// For example: "optimize:(euro;material;energy)"
//
// Behavior:
//   - Checks if an optimizer has already been defined (returns error if duplicate found).
//   - Strips the "optimize:" prefix.
//   - Removes surrounding parentheses.
//   - Splits the remaining string by semicolons to extract individual optimization targets.
//   - Appends non-empty trimmed targets to config.OptimizeTargets.
//   - Sets HasOptimizer flag to true.
//
// Parameters:
//   - config: a pointer to the ConfigData struct to be updated.
//   - line: a string representing the line containing optimization targets.
//
// Returns:
//   - An error if a duplicate optimizer is detected, nil otherwise.
func parseOptimize(config *ConfigData, line string) error {
	// Check if optimizer has already been defined
	if config.HasOptimizer {
		return fmt.Errorf("multiple optimize declarations found")
	}

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

	// Mark that we've processed an optimizer
	config.HasOptimizer = true
	return nil
}

// parseProcess parses a process definition line and updates the provided ConfigData.
//
// The expected format for the line is: "process_name:(needs):(results):cycles".
// For example: "build_product:(material:1):(product:1):30"
//
// Components:
//   - process_name: The unique identifier for the process.
//   - (needs): A parentheses-enclosed list of required resources in format "resource:quantity;resource:quantity;...".
//   - (results): A parentheses-enclosed list of produced resources in format "resource:quantity;resource:quantity;...".
//   - cycles: An integer representing the number of cycles required to complete the process.
//
// Behavior:
//   - Extracts the process name from the beginning of the line up to the first colon.
//   - Splits the remainder into three parts: needs, results, and cycles.
//   - Parses needs and results using parseResourceMap to convert them into resource maps.
//   - Converts the cycles string to an integer.
//   - Creates a new Process struct and appends it to config.Processes.
//
// Parameters:
//   - config: a pointer to the ConfigData struct to be updated.
//   - line: a string representing the process definition line.
//
// Returns:
//   - An error if the line format is invalid, resource parsing fails, or cycles is not a valid integer.
func parseProcess(config *ConfigData, line string) error {
	// Find the first colon to separate name from the rest
	colonIndex := strings.Index(line, ":")
	if colonIndex == -1 {
		return fmt.Errorf("%s", line)
	}

	name := strings.TrimSpace(line[:colonIndex])
	if name == "" {
		return fmt.Errorf("%s", line)
	}
	rest := line[colonIndex+1:]

	// Parse the remaining parts: (needs):(results):cycles
	parts := strings.Split(rest, "):")
	if len(parts) != 3 {
		return fmt.Errorf("%s", line)
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

// parseResourceMap parses a string representing a resource map in the format "(name1:qty1;name2:qty2;...)".
// It returns a map where the keys are resource names and the values are their corresponding quantities.
// The input string should be enclosed in parentheses, with each resource separated by a semicolon and
// each resource specified as "name:quantity". If the input string is empty or contains no resources,
// an empty map is returned. Returns an error if any resource entry is malformed or if a quantity cannot
// be parsed as an integer.
//
// Example input: "(iron:2;coal:3)"
// Example output: map[string]int{"iron": 2, "coal": 3}
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
