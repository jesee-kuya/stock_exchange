package engine

import (
	"fmt"

	"github.com/jesee-kuya/stock_exchange/process"
	"github.com/jesee-kuya/stock_exchange/util"
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

// Stock represents the available items and their quantities.
type Stock struct {
	Items map[string]int
}


type ScheduleEntry struct {
	Cycle       int
	ProcessName string
}

func (e *Engine) LoadConfig(path string) error {
	config, err := util.ParseConfig(path)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	e.Stock = &Stock{Items: config.Stocks}
	e.Processes = config.Processes
	e.OptimizeTargets = config.OptimizeTargets

	return nil
}

func (e *Engine) Run(waitingTime string) {
    e.Cycle = 0
    maxIdleCycles := 1000 
    idleCycles := 0

    for {
        executed := false

        for _, proc := range e.Processes {
            // Check if process can run (enough stock for needs)
            canRun := true
            for item, qty := range proc.Needs {
                if e.Stock.Items[item] < qty {
                    canRun = false
                    break
                }
            }
            if !canRun {
                continue
            }

            // Deduct needs from stock
            for item, qty := range proc.Needs {
                e.Stock.Items[item] -= qty
            }

            for item, qty := range proc.Result {
                e.Stock.Items[item] += qty
            }

            // Log execution
            fmt.Printf("Cycle %d: Executed process %s\n", e.Cycle, proc.Name)
            executed = true
        }

        e.Cycle++

        // Simulate waiting time between cycles
        util.Wait(waitingTime)

        if !executed {
            idleCycles++
        } else {
            idleCycles = 0
        }

        // Stop if no process was executed for maxIdleCycles
        if idleCycles >= maxIdleCycles {
            fmt.Println("No more executable processes or max idle cycles reached. Stopping simulation.")
            break
        }
    }
}
