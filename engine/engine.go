package engine

import (
	"fmt"

	"github.com/jesee-kuya/stock_exchange/process"
	"github.com/jesee-kuya/stock_exchange/util"
)

// Engine is the main structure for executing and optimizing
// the scheduling process defined in a configuration file
type Engine struct {
	Stock           *Stock
	Processes       []*process.Process
	Schedule        []string
	Cycle           int
	OptimizeTargets []string
}

// Stock represents the available items and their quantities.
type Stock struct {
	Items map[string]int
}


type ScheduleEntry struct {
	Cycle       int
	ProcessName string
}


// LoadConfig loads the engine configuration from the specified file path.
// It parses the configuration file, initializes the Stock, Processes, and OptimizeTargets
// fields of the Engine based on the parsed data, and returns an error if parsing fails.
//
// Parameters:
//   - path: The file path to the configuration file.
//
// Returns:
//   - error: An error if the configuration could not be parsed, otherwise nil.
func (e *Engine) LoadConfig(path string) error {
	config, err := util.ParseConfig(path)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	e.Stock = &Stock{Items: config.Stocks}
	e.Processes = config.Processes
	e.OptimizeTargets = config.OptimizeTargets

	return nil
}
