package checker

import (
	"github.com/jesee-kuya/stock_exchange/util"
)

// LoadConfig loads the configuration data from the specified file path.
// It uses the util.ParseConfig function to parse the configuration file.
// The parsed configuration data is then used to populate the Stocks and Processes fields of the Checker instance.
//
// Parameters:
//   - path: A string representing the file path to the configuration file.
//
// Returns:
//   - error: An error if the configuration file cannot be parsed or loaded successfully, otherwise nil.
func (c *Checker) LoadConfig(path string) error {
	configData, err := util.ParseConfig(path)
	if err != nil {
		return err
	}
	
	c.Stocks = configData.Stocks
	c.Processes = configData.Processes
	

	return nil
}
