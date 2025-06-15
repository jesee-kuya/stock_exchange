package engine

import "github.com/jesee-kuya/stock_exchange/process"

// NewEngine creates and initializes a new Engine instance.
func NewEngine() *Engine {
	return &Engine{
		Stock:           &Stock{Items: make(map[string]int)},
		Processes:       []*process.Process{},
		Schedule:        []string{},
		Cycle:           0,
		OptimizeTargets: []string{},
	}
}
