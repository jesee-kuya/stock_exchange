package process

func (p *Process) Run(stocks map[string]int, pending map[int]map[string]int, currentCycle int) {
	// Deduct required input items from the stock
	for item, requiredQty := range p.Needs {
		stocks[item] -= requiredQty
	}

	// Schedule result items to be added after the process cycle delay
	dueCycle := currentCycle + p.Cycle
	if _, exists := pending[dueCycle]; !exists {
		pending[dueCycle] = make(map[string]int)
	}
	for item, producedQty := range p.Result {
		pending[dueCycle][item] += producedQty
	}
}
