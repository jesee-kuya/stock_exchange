package checker

import (
	"github.com/jesee-kuya/stock_exchange/engine"
	"github.com/jesee-kuya/stock_exchange/process"
)

// NewChecker creates and initializes a new Checker instance.
// It returns a pointer to a Checker with properly initialized fields.
func NewChecker() *Checker {
	return &Checker{
		Stocks:    make(map[string]int),
		Processes: []*process.Process{},
		Log:       []engine.ScheduleEntry{},
	}
}
