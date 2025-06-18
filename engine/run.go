package engine

import (
	"fmt"

	"github.com/jesee-kuya/stock_exchange/process"
	"github.com/jesee-kuya/stock_exchange/util"
)

// Run executes the main simulation loop of the Engine for a specified duration.
// The waitingTime parameter is a string representing the maximum simulation time,
// which is parsed into a number of cycles. In each cycle, the function processes
// running processes, updates the stock with completed process results, and attempts
// to schedule new processes if their resource requirements are met. The simulation
// ends when the maximum number of cycles is reached or no more processes can be executed.
// The function prints detailed information about each cycle, including resource usage,
// process scheduling, and simulation status.
func (e *Engine) Run(waitingTime string) {
	maxCycles, err := util.ParseDuration(waitingTime)
	if err != nil {
		fmt.Printf("Invalid Cycle: %v\n", err)
		return
	}

	type runningProcess struct {
		Process *process.Process
		Delay   int
	}

	type scheduledProcess struct {
		Cycle int
		Name  string
	}

	var running []runningProcess
	var started []scheduledProcess

	e.Schedule = []string{}
	e.Cycle = 0

	for e.Cycle < maxCycles {
		// Finish any running process
		var updatedRunning []runningProcess
		for _, rp := range running {
			rp.Delay--
			if rp.Delay == 0 {
				for item, qty := range rp.Process.Result {
					e.Stock.Items[item] += qty
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
				for item, qty := range p.Needs {
					e.Stock.Items[item] -= qty
				}

				running = append(running, runningProcess{
					Process: p,
					Delay:   p.Cycle,
				})

				started = append(started, scheduledProcess{
					Cycle: e.Cycle,
					Name:  p.Name,
				})

				e.Schedule = append(e.Schedule, p.Name)
				executed = true
			}
		}

		if !executed && len(running) == 0 {
			break
		}

		e.Cycle++
	}

	// Output in expected format
	fmt.Println("Main Processes:")
	for _, s := range started {
		fmt.Printf(" %d:%s\n", s.Cycle, s.Name)
	}

	fmt.Printf("No more process doable at cycle %d\n", e.Cycle)

	fmt.Println("Stock:")
	for item, qty := range e.Stock.Items {
		fmt.Printf(" %s => %d\n", item, qty)
	}
}
