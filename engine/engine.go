package engine

import "github.com/jesee-kuya/stock_exchange/process"

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
