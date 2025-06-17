package util

import (
	"testing"

	"github.com/jesee-kuya/stock_exchange/process"
)

// TestParseLine tests the parseLine function with various input strings representing different configuration commands.
// It verifies that parseLine correctly parses stock, optimize, and process commands without returning an error.
// Each test case initializes a fresh ConfigData instance and checks that parsing the rawData string does not produce an error.
func TestParseLine(t *testing.T) {
	// configData := &ConfigData{}
	testCases := []struct {
		name    string
		rawData string
	}{
		{"parse stock", "cabinet:1"},
		{"parse optimize", "optimize:(time;cabinet)"},
		{"parse process", "do_shelf:(board:1):(shelf:1):10"},
		{"parse more items", "do_cabinet:(doorknobs:2;background:1;shelf:3):(cabinet:1):30"},
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
