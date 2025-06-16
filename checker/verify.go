package checker

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
	}

	return nil
}
