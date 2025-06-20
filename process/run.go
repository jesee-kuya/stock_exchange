package process

// Run processes the consumption of required items from the current stock and schedules the produced items
// to be added to the stock after a specified number of cycles. It deducts the quantities specified in p.Needs
// from the stocks map, and adds the quantities specified in p.Result to the pending map for the due cycle.
// Parameters:
//   - stocks: a map representing the current available quantities of each item.
//   - pending: a map where produced items are scheduled to be added to stocks in future cycles.
//   - currentCycle: the current cycle number, used to determine when produced items become available.
func (p *Process) Run(stocks map[string]int, pending map[int]map[string]int, currentCycle int) {
	for item, requiredQty := range p.Needs {
		stocks[item] -= requiredQty
	}

	dueCycle := currentCycle + p.Cycle
	if _, exists := pending[dueCycle]; !exists {
		pending[dueCycle] = make(map[string]int)
	}
	for item, producedQty := range p.Result {
		pending[dueCycle][item] += producedQty
	}
}
