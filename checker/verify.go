package checker

func (c *Checker) Verify() error {
	stocks := make(map[string]int)
	for k, v := range c.Stocks {
		stocks[k] = v
	}

	return nil
}
