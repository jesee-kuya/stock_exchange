package checker

import (
	"github.com/jesee-kuya/stock_exchange/engine"
	"github.com/jesee-kuya/stock_exchange/process"
)

// Checker represents a structure used to manage and track stock-related data,
// associated processes, and scheduling logs in a stock exchange system.
//
// Fields:
// - Stocks: A map where the keys are stock names (string) and the values are their respective quantities (int).
// - Processes: A slice of pointers to Process objects, representing the processes associated with the stock exchange.
// - Log: A slice of ScheduleEntry objects from the engine package, used to record scheduling events or logs.
type Checker struct {
	Stocks    map[string]int
	Processes []*process.Process
	Log       []engine.ScheduleEntry
}
