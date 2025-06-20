package checker

import (
	"github.com/jesee-kuya/stock_exchange/engine"
	"github.com/jesee-kuya/stock_exchange/process"
)

// NewChecker creates and returns a new instance of Checker with initialized fields.
// It sets up an empty map for Stocks, an empty slice for Processes, and an empty log.
//
// This function is typically used to initialize the Checker before use.
func NewChecker() *Checker {
	return &Checker{
		Stocks:    make(map[string]int),
		Processes: []*process.Process{},
		Log:       []engine.ScheduleEntry{},
	}
}
