package process

// CanRun checks if the process has sufficient resources to run.
// It takes a map of available stocks, where the key is the resource name
// and the value is the quantity available. The function iterates over the
// resources required by the process (p.Needs) and returns false if any
// required resource is missing or insufficient in the provided stocks map.
// Returns true only if all required resources are available in the needed quantities.
func (p *Process) CanRun(stocks map[string]int) bool {
	for resource, required := range p.Needs {
		if available, ok := stocks[resource]; !ok || available < required {
			return false
		}
	}
	return true
}
