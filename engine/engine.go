package engine

import (
	"github.com/jesee-kuya/stock_exchange/process"
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

// Stock represents the available items in the system.
// The Items map stores item names as keys and their corresponding quantities as values.
type Stock struct {
	Items map[string]int
}

// ScheduleEntry represents a scheduled execution of a process.
// It records the cycle at which the process is run and the name of the process.
type ScheduleEntry struct {
	Cycle       int    // The cycle number when the process is executed
	ProcessName string // The name of the process being executed
}
