package engine

import (
	"fmt"
	"sort"
	"time"

	"github.com/jesee-kuya/stock_exchange/process"
	"github.com/jesee-kuya/stock_exchange/util"
)

type runningProcess struct {
	Process *process.Process
	Delay   int
}

func (e *Engine) Run(waitingTime string) {
	maxSeconds, err := util.ParseDuration(waitingTime)
	if err != nil {
		fmt.Println("Invalid waiting time format:", err)
		return
	}

	start := time.Now()
	e.Schedule = []string{}
	e.Cycle = 0
	running := []runningProcess{}

	fmt.Println("Main Processes :")

	targets := map[string]bool{}
	for _, t := range e.OptimizeTargets {
		targets[t] = true
	}
	priorities := computePriorities(e.Processes, targets)

	for {
		if time.Since(start).Seconds() >= float64(maxSeconds) {
			fmt.Printf("Time limit exceeded after %d cycles\n", e.Cycle)
			break
		}

		running = updateRunningProcesses(running, e)

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

		if len(running) == 0 && !e.canRunAny() {
			fmt.Printf("No more process doable at cycle %d\n", e.Cycle + 1)
			break
		}

		e.Cycle++
	}

	printStock(e.Stock)
}

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

func (e *Engine) canRunAny() bool {
	for _, p := range e.Processes {
		if p.CanRun(e.Stock.Items) {
			return true
		}
	}
	return false
}

func computePriorities(processes []*process.Process, targets map[string]bool) map[string]int {
	prio := map[string]int{}
	visited := map[string]bool{}
	queue := []struct {
		Name  string
		Depth int
	}{}

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

	// Fallback for unreachable
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
