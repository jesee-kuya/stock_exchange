package engine

import (
	"github.com/jesee-kuya/stock_exchange/util"
)

// LoadConfig loads and parses a configuration file to initialize the Engine with
// stock items, processes, and optimization targets.
//
// This method reads a configuration file from the specified path and populates
// the Engine's internal state with the parsed data. The configuration file should
// contain stock definitions, process definitions, and optimization targets in the
// expected format.
//
// Configuration file format:
//   - Stock definitions: "item_name:quantity" (e.g., "euro:10")
//   - Process definitions: "name:(needs):(results):cycles" (e.g., "build_product:(material:1):(product:1):30")
//   - Optimization targets: "optimize:(target1;target2;...)" (e.g., "optimize:(time;client_content)")
//   - Comments: Lines starting with '#' are ignored
//   - Empty lines are ignored
//
// Behavior:
//   - Delegates file parsing to util.ParseConfig() which validates the file format.
//   - Initializes the Engine's Stock with the parsed stock items and quantities.
//   - Populates the Engine's Processes slice with the parsed process definitions.
//   - Sets the Engine's OptimizeTargets with the parsed optimization goals.
//   - Maintains the original structure and relationships between parsed data.
//
// Parameters:
//   - path: a string representing the file path to the configuration file.
//
// Returns:
//   - An error if the file cannot be read, parsed, or contains invalid format.
//   - nil if the configuration is successfully loaded and the Engine is initialized.
//
// Example usage:
//
//	engine := &Engine{}
//	err := engine.LoadConfig("config/example.conf")
//	if err != nil {
//	    log.Fatal("Failed to load config:", err)
//	}
func (e *Engine) LoadConfig(path string) error {
	config, err := util.ParseConfig(path)
	if err != nil {
		return err
	}
	e.Stock = &Stock{Items: config.Stocks}
	e.Processes = config.Processes
	e.OptimizeTargets = config.OptimizeTargets
	return nil
}
