package checker

import (
	"github.com/jesee-kuya/stock_exchange/engine"
	"github.com/jesee-kuya/stock_exchange/process"
)

type Checker struct {
	Stocks    map[string]int
	Processes []*process.Process
	Log       []engine.ScheduleEntry
}
