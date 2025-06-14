package process


func (p *models.Process) CanRun(stocks map[string]int) bool {
	for item, requiredQty := range p.Input {
		if stocks[item] < requiredQty {
			return false
		}
	}
	return true
}