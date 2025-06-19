package engine

import (
	"fmt"
	"sort"
	"time"

	"github.com/jesee-kuya/stock_exchange/process"
	"github.com/jesee-kuya/stock_exchange/util"
)

// runningProcess represents a process that is currently executing.
// It tracks the process instance and the remaining delay until completion.
type runningProcess struct {
	Process *process.Process // The process being executed
	Delay   int              // Remaining cycles until the process completes
}

// Run executes the stock exchange optimization algorithm for the specified duration.
// It implements a priority-based parallel schedule generation scheme that:
// 1. Schedules processes based on their priority (derived from optimization targets)
// 2. Runs multiple processes concurrently when resources allow
// 3. Respects process dependencies and resource constraints
// 4. Stops when the time limit is reached or no more processes can be scheduled
//
// Parameters:
//   - waitingTime: Maximum execution time in string format (e.g., "10s", "5m")
//
// The function prints the execution schedule and final stock state.
// It generates a log file with the schedule in the format "cycle:process_name".
func (e *Engine) Run(waitingTime string) {
	maxSeconds, err := util.ParseDuration(waitingTime)
	if err != nil {
		fmt.Println("Invalid waiting time format:", err)
		return
	}

	if !e.canRunAny() {
		fmt.Println(" Missing processes\n Exiting... ")
		return
	}

	start := time.Now()
	e.Schedule = []string{}
	e.Cycle = 0
	running := []runningProcess{}

	fmt.Println("Main Processes :")

	targets := map[string]bool{}
	for _, t := range e.OptimizeTargets {
		if _, ok := e.Stock.Items[t]; ok {
			targets[t] = true
		}
	}

	priorities := computePriorities(e.Processes, targets)
	timeExceeded := false

	for {
		// Check time limit but don't break immediately if processes are still running
		if time.Since(start).Seconds() >= float64(maxSeconds) {
			timeExceeded = true
		}

		running = updateRunningProcesses(running, e)

		// Only schedule new processes if time hasn't exceeded
		if !timeExceeded {
			// Get all runnable processes for this cycle
			runnable := []*process.Process{}
			for _, p := range e.Processes {
				if p.CanRun(e.Stock.Items) {
					runnable = append(runnable, p)
				}
			}

			// Sort by priority (descending) and then by name (descending)
			sort.Slice(runnable, func(i, j int) bool {
				pi, pj := priorities[runnable[i].Name], priorities[runnable[j].Name]
				if pi == pj {
					return runnable[i].Name > runnable[j].Name
				}
				return pi > pj
			})

			// Use a copy of stock for simulation
			stockCopy := make(map[string]int)
			for k, v := range e.Stock.Items {
				stockCopy[k] = v
			}

			scheduledCount := make(map[*process.Process]int)
			changed := true
			for changed {
				changed = false
				for _, p := range runnable {
					if p.CanRun(stockCopy) {
						// Consume resources in the simulated stock
						for item, qty := range p.Needs {
							stockCopy[item] -= qty
						}
						scheduledCount[p]++
						changed = true
					}
				}
			}

			// Schedule the processes and update real stock
			scheduledEntries := []string{}
			for _, p := range runnable {
				count := scheduledCount[p]
				for i := 0; i < count; i++ {
					// Update real stock
					for item, qty := range p.Needs {
						e.Stock.Items[item] -= qty
					}
					// Add to running processes
					running = append(running, runningProcess{
						Process: p,
						Delay:   p.Cycle,
					})
					// Create schedule entry
					entry := fmt.Sprintf(" %d:%s", e.Cycle, p.Name)
					scheduledEntries = append(scheduledEntries, entry)
					e.Schedule = append(e.Schedule, entry)
				}
			}

			// Print all entries for this cycle
			for _, entry := range scheduledEntries {
				fmt.Println(entry)
			}

			// Check if we can continue (only if time hasn't exceeded)
			if len(running) == 0 && !e.canRunAny() {
				fmt.Printf("No more process doable at cycle %d\n", e.Cycle+1)
				break
			}
		} else {
			// Time exceeded, just let running processes complete
			if len(running) == 0 {
				break
			}
		}

		e.Cycle++

		// If time exceeded and no processes are running, we can safely exit
		if timeExceeded && len(running) == 0 {
			fmt.Printf("Time limit exceeded after %d cycles\n", e.Cycle)
			break
		}
	}

	printStock(e.Stock)
}

// updateRunningProcesses decrements the delay of all running processes and
// completes those whose delay has reached zero.
// When a process completes, its results are added to the engine's stock.
//
// Parameters:
//   - running: Slice of currently running processes
//   - e: The engine instance to update stock when processes complete
//
// Returns:
//   - Updated slice of running processes with completed ones removed
func updateRunningProcesses(running []runningProcess, e *Engine) []runningProcess {
	next := []runningProcess{}
	for _, rp := range running {
		rp.Delay--
		if rp.Delay <= 0 {
			for item, qty := range rp.Process.Result {
				e.Stock.Items[item] += qty
			}
		} else {
			next = append(next, rp)
		}
	}
	return next
}

// canRunAny checks if any process in the engine can be executed
// with the current stock levels.
//
// Returns:
//   - true if at least one process can run, false otherwise
func (e *Engine) canRunAny() bool {
	for _, p := range e.Processes {
		if p.CanRun(e.Stock.Items) {
			return true
		}
	}
	return false
}

// computePriorities calculates priority values for all processes based on
// their relationship to optimization targets using a breadth-first search approach.
// Processes that directly produce optimization targets get priority 0,
// processes that produce inputs for those get priority 1, and so on.
// This implements a backward chaining algorithm to establish process dependencies.
//
// Priority calculation logic:
// - Direct producers of optimization targets: priority 0 (highest)
// - Processes that support direct producers: priority 1
// - Processes that support priority 1 processes: priority 2
// - And so on...
// - Unreachable processes get the lowest priority (max + 1)
//
// Parameters:
//   - processes: All available processes
//   - targets: Map of optimization target items
//
// Returns:
//   - Map of process names to their priority values (lower = higher priority)
func computePriorities(processes []*process.Process, targets map[string]bool) map[string]int {
	prio := map[string]int{}
	visited := map[string]bool{}
	queue := []struct {
		Name  string
		Depth int
	}{}

	// Initialize with processes that directly produce optimization targets
	for _, p := range processes {
		for result := range p.Result {
			if targets[result] {
				prio[p.Name] = 0
				queue = append(queue, struct {
					Name  string
					Depth int
				}{p.Name, 0})
			}
		}
	}

	// Breadth-first search to assign priorities based on dependency depth
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		visited[curr.Name] = true

		var proc *process.Process
		for _, p := range processes {
			if p.Name == curr.Name {
				proc = p
				break
			}
		}

		// Find processes that produce what this process needs
		for need := range proc.Needs {
			for _, p := range processes {
				if _, ok := p.Result[need]; ok {
					if prio[p.Name] < curr.Depth+1 {
						prio[p.Name] = curr.Depth + 1
						queue = append(queue, struct {
							Name  string
							Depth int
						}{p.Name, curr.Depth + 1})
					}
				}
			}
		}
	}

	// Assign fallback priority for processes not reachable from targets
	max := 0
	for _, v := range prio {
		if v > max {
			max = v
		}
	}
	for _, p := range processes {
		if _, ok := prio[p.Name]; !ok {
			prio[p.Name] = max + 1
		}
	}

	return prio
}

// printStock displays the final state of all stock items in alphabetical order.
// This provides a clear summary of remaining resources after process execution.
//
// Parameters:
//   - stock: The stock instance containing all items and their quantities
func printStock(stock *Stock) {
	fmt.Println("Stock:")
	keys := make([]string, 0, len(stock.Items))
	for k := range stock.Items {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf(" %s => %d\n", k, stock.Items[k])
	}
}
