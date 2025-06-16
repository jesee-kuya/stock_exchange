package checker

func (c *Checker) Verify() error {
	stocks := make(map[string]int)
	for k, v := range c.Stocks {
		stocks[k] = v
	}

	pending := make(map[int]map[string]int)
	currentCycle := 0

	return nil
}
