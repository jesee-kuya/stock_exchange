package checker

import (
	"fmt"

	"github.com/jesee-kuya/stock_exchange/process"
)

// Verify checks the validity of the process execution log against the initial stock levels and process definitions.
// It simulates the consumption and production of items by processes over cycles, ensuring that:
//   - Each process in the log exists in the list of known processes.
//   - Sufficient stock is available for each process's needs at the time it is executed.
//   - Outputs from processes are applied after the required number of cycles.
// If any inconsistency is found (such as an unknown process or insufficient stock), an error is returned
// describing the issue and the cycle at which it occurred. If the log is valid, it returns nil.
func (c *Checker) Verify() error {
	stocks := make(map[string]int)
	for k, v := range c.Stocks {
		stocks[k] = v
	}

	// Pending outputs map: cycle -> items
	pending := make(map[int]map[string]int)
	currentCycle := 0

	for _, entry := range c.Log {
		fmt.Printf("Evaluating: %d:%s\n", entry.Cycle, entry.ProcessName)

		// Apply any pending outputs from prior cycles
		for cycle := currentCycle; cycle <= entry.Cycle; cycle++ {
			if outputs, ok := pending[cycle]; ok {
				for item, qty := range outputs {
					stocks[item] += qty
				}
				delete(pending, cycle)
			}
		}
		currentCycle = entry.Cycle

		// Find the process
		var proc *process.Process
		for _, p := range c.Processes {
			if p.Name == entry.ProcessName {
				proc = p
				break
			}
		}
		if proc == nil {
			return fmt.Errorf("unknown process '%s' at cycle %d", entry.ProcessName, entry.Cycle)
		}

		// Check if enough stock exists
		for item, qty := range proc.Needs {
			if stocks[item] < qty {
				return fmt.Errorf("insufficient stock for '%s' at cycle %d: need %d %s, have %d",
					proc.Name, entry.Cycle, qty, item, stocks[item])
			}
		}

		// Deduct input from stocks
		for item, qty := range proc.Needs {
			stocks[item] -= qty
		}

		// Schedule outputs
		dueCycle := entry.Cycle + proc.Cycle
		if pending[dueCycle] == nil {
			pending[dueCycle] = make(map[string]int)
		}
		for item, qty := range proc.Result {
			pending[dueCycle][item] += qty
		}
	}

	// Flush remaining pending outputs
	for _, outputs := range pending {
		for item, qty := range outputs {
			stocks[item] += qty
		}
	}

	fmt.Println("Trace completed. No error detected.")
	return nil
}
