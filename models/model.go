package models

type Engine struct {
	Stock           *Stock
	Processes       []*Process
	Schedule        []string
	Cycle           int
	OptimizeTargets []string
}
