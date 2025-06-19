package engine

import (
	"github.com/jesee-kuya/stock_exchange/util"
)

func (e *Engine) LoadConfig(path string) error {
	config, err := util.ParseConfig(path)
	if err != nil {
		return err
	}

	e.Stock = &Stock{Items: config.Stocks}
	e.Processes = config.Processes
	e.OptimizeTargets = config.OptimizeTargets

	return nil
}
