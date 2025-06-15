package process

import (
	"github.com/jesee-kuya/stock_exchange/engine"
)

func (p *engine.Process) Run(stocks map[string]int) {
	// Deduct input items from the stock
	for item, requiredQty := range p.Input {
		stocks[item] -= requiredQty
	}

	// Add output items to the stock
	for item, producedQty := range p.Output {
		stocks[item] += producedQty
	}

}
