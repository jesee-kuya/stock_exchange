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
	maxCycles, err := util.ParseDuration(waitingTime)
	if err != nil {
		fmt.Printf("Invalid waiting time: %v\n", err)
		return
	}

	fmt.Printf("Starting simulation (max cycles: %d)\n", maxCycles)

	type runningProcess struct {
		Process *process.Process
		Delay   int
	}

	var running []runningProcess
	e.Schedule = []string{}
	e.Cycle = 0

	for e.Cycle < maxCycles {
		fmt.Printf("Cycle %d\n", e.Cycle)

		// Process completion step
		var updatedRunning []runningProcess
		for _, rp := range running {
			rp.Delay--
			if rp.Delay == 0 {
				// Add result items to stock
				for item, qty := range rp.Process.Result {
					e.Stock.Items[item] += qty
					fmt.Printf("  [+] %d %s (from %s)\n", qty, item, rp.Process.Name)
				}
			} else {
				updatedRunning = append(updatedRunning, rp)
			}
		}
		running = updatedRunning

		// Try to schedule new processes
		executed := false
		for _, p := range e.Processes {
			if p.CanRun(e.Stock.Items) {
				// Deduct required resources
				for item, qty := range p.Needs {
					e.Stock.Items[item] -= qty
					fmt.Printf("  [-] %d %s (used by %s)\n", qty, item, p.Name)
				}

				// Start the process
				running = append(running, runningProcess{
					Process: p,
					Delay:   p.Cycle,
				})
				e.Schedule = append(e.Schedule, p.Name)
				fmt.Printf("  [*] Scheduled process: %s\n", p.Name)
				executed = true
			}
		}

		if !executed && len(running) == 0 {
			fmt.Println("No more executable processes. Ending simulation.")
			break
		}

		e.Cycle++
	}

	fmt.Printf("Simulation ended after %d cycles\n", e.Cycle)
}

