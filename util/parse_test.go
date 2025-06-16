package util

import (
	"testing"

	"github.com/jesee-kuya/stock_exchange/process"
)

func TestParseLine(t *testing.T) {
	// configData := &ConfigData{}
	testCases := []struct {
		name    string
		rawData string
	}{
		{"parse stock", "cabinet:1"},
		{"parse optimize", "optimize:(time;cabinet)"},
		{"parse process", "do_shelf:(board:1):(shelf:1):10"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			configData := &ConfigData{
				Stocks:          map[string]int{},
				Processes:       []*process.Process{},
				OptimizeTargets: []string{},
			}
			err := parseLine(configData, tc.rawData)
			if err != nil {
				t.Errorf("%v", err)
			}
		})
	}
}
