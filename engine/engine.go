package engine

import (
	"fmt"

	"github.com/jesee-kuya/stock_exchange/process"
	"github.com/jesee-kuya/stock_exchange/util"
)

// Engine represents the core simulation engine for the stock exchange system.
// It manages the available stock, the set of processes that can be executed,
//
// Fields:
//   - Stock: Pointer to the Stock struct representing available items and their quantities.
//   - Processes: Slice of pointers to Process structs, representing all possible processes.
//   - Schedule: Slice of strings recording the names of scheduled processes in order of execution.
//   - Cycle: Integer tracking the current simulation cycle.
//   - OptimizeTargets: Slice of strings specifying the optimization targets for the simulation.
type Engine struct {
	Stock           *Stock
	Processes       []*process.Process
	Schedule        []string
	Cycle           int
	OptimizeTargets []string
}


// Stock represents the available items in the system.
// The Items map stores item names as keys and their corresponding quantities as values.
type Stock struct {
	Items map[string]int
}

// ScheduleEntry represents a scheduled execution of a process.
// It records the cycle at which the process is run and the name of the process.
type ScheduleEntry struct {
	Cycle       int    // The cycle number when the process is executed
	ProcessName string // The name of the process being executed
}

// LoadConfig loads and parses a configuration file to initialize the Engine's state.
//
// It performs the following actions:
// - Parses the configuration file at the given path.
// - Initializes the Stock with item quantities from the config.
// - Loads the list of processes from the config.
// - Sets the optimization targets.
//
// Returns an error if parsing the configuration fails.
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

// Run executes the simulation of process scheduling and execution over a specified duration.
//
// It simulates each cycle of the engine until either the specified maximum duration is reached
// or no more processes can be executed.
//
// Parameters:
//   - waitingTime: a string representing the maximum simulation time, which will be parsed into cycles.
//
// Behavior:
//   - Parses the waitingTime into an integer number of cycles.
//   - Initializes internal engine state: schedule list and cycle counter.
//   - For each cycle:
//       - Completes any running processes and updates the stock with their output.
//       - Attempts to start new processes if their input needs can be fulfilled from the current stock.
//       - Deducts used items from stock and schedules the process for execution based on its delay.
//   - The simulation ends early if no process can be started and none are still running.
//
// Notes:
//   - Each process has a defined delay (`Cycle`) representing how long it takes to complete.
//   - Results and actions are printed to standard output at every step.
func (e *Engine) Run(waitingTime string) {
	maxCycles, err := util.ParseDuration(waitingTime)
	if err != nil {
		fmt.Printf("Invalid Cycle: %v\n", err)
		return
	}

	fmt.Printf("Main Processes:\n")

	type runningProcess struct {
		Process *process.Process
		Delay   int
	}

	var running []runningProcess
	e.Schedule = []string{}
	e.Cycle = 0

	for e.Cycle < maxCycles {
		fmt.Printf("Cycle %d\n", e.Cycle)

		var updatedRunning []runningProcess
		for _, rp := range running {
			rp.Delay--
			if rp.Delay == 0 {
				for item, qty := range rp.Process.Result {
					e.Stock.Items[item] += qty
					fmt.Printf("  [+] %d %s (from %s)\n", qty, item, rp.Process.Name)
				}
			} else {
				updatedRunning = append(updatedRunning, rp)
			}
		}
		running = updatedRunning
		executed := false
		for _, p := range e.Processes {
			if p.CanRun(e.Stock.Items) {
				for item, qty := range p.Needs {
					e.Stock.Items[item] -= qty
					fmt.Printf("  [-] %d %s (used by %s)\n", qty, item, p.Name)
				}
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


