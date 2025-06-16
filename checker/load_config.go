package checker

import (
	"github.com/jesee-kuya/stock_exchange/util"
)

// LoadConfig parses a configuration file and populates the Checker with
// initial stocks, processes, and optimization goals.
func (c *Checker) LoadConfig(path string) error {
	configData, err := util.ParseConfig(path)
	if err != nil {
		return err
	}
	// Populate checker fields
	c.Stocks = configData.Stocks

	return nil
}
