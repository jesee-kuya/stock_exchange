package engine

import "github.com/jesee-kuya/stock_exchange/process"

// NewEngine creates and returns a pointer to a new Engine instance with initialized fields.
// The Stock field is initialized with an empty map of items.
// The Processes field is initialized as an empty slice of *process.Process.
// The Schedule field is initialized as an empty slice of strings.
// The Cycle field is set to 0.
// The OptimizeTargets field is initialized as an empty slice of strings.
func NewEngine() *Engine {
	return &Engine{
		Stock:           &Stock{Items: make(map[string]int)},
		Processes:       []*process.Process{},
		Schedule:        []string{},
		Cycle:           0,
		OptimizeTargets: []string{},
	}
}
