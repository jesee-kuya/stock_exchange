package process

func (p *Process) CanRun(stocks map[string]int) bool {
	for resource, required := range p.Needs {
		if available, ok := stocks[resource]; !ok || available < required {
			return false
		}
	}
	return true
}
