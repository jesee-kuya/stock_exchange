package process

import (
	"github.com/jesee-kuya/stock_exchange/engine"
)

func (p *engine.Process) Run(stocks map[string]int, pending map[int]map[string]int, currentCycle int) {
    // Deduct input items from the stock
    for item, requiredQty := range p.Input {
        stocks[item] -= requiredQty
    }

    // Schedule output items to be added after process.Cycle duration
    dueCycle := currentCycle + p.Cycle
    if pending[dueCycle] == nil {
        pending[dueCycle] = make(map[string]int)
    }
    for item, producedQty := range p.Output {
        pending[dueCycle][item] += producedQty
    }
}
