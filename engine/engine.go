package engine

// Engine is the main structure for executing and optimizing
// the scheduling process defined in a configuration file
type Engine struct {
	Stock           *Stock
	Processes       []*Process
	Schedule        []string
	Cycle           int
	OptimizeTargets []string
}


type Process struct {
	Name   string
	Input  map[string]int
	Output map[string]int
}