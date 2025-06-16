package checker

import (
	"fmt"

	"github.com/jesee-kuya/stock_exchange/process"
)

func (c *Checker) Verify() error {
	stocks := make(map[string]int)
	for k, v := range c.Stocks {
		stocks[k] = v
	}

	pending := make(map[int]map[string]int)
	currentCycle := 0

	for _, entry := range c.Log {
		for cycle := currentCycle; cycle <= entry.Cycle; cycle++ {
			if outputs, ok := pending[cycle]; ok {
				for item, qty := range outputs {
					stocks[item] += qty
				}
				delete(pending, cycle)
			}
		}
		currentCycle = entry.Cycle

		var proc *process.Process
		for _, p := range c.Processes {
			if p.Name == entry.Name {
				proc = p
				break
			}
		}
		if proc == nil {
			return fmt.Errorf("unknown process '%s' at cycle %d", entry.Name, entry.Cycle)
		}

		for item, qty := range proc.Needs {
			if stocks[item] < qty {
				return fmt.Errorf("insufficient stock for '%s' at cycle %d: need %d %s, have %d",
					proc.Name, entry.Cycle, qty, item, stocks[item])
			}
		}

		for item, qty := range proc.Needs {
			stocks[item] -= qty
		}

		dueCycle := entry.Cycle + proc.Cycle
		if pending[dueCycle] == nil {
			pending[dueCycle] = make(map[string]int)
		}
		for item, qty := range proc.Result {
			pending[dueCycle][item] += qty
		}
	}

	for _, outputs := range pending {
		for item, qty := range outputs {
			stocks[item] += qty
		}
	}

	fmt.Println("Trace completed. No error detected.")
	fmt.Println("Final stocks:")
	for item, qty := range stocks {
		fmt.Printf("  %s => %d\n", item, qty)
	}

	return nil
}
