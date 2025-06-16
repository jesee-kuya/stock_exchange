package util

import "testing"

func TestParseLine(t *testing.T) {
	// configData := &ConfigData{}
	testCases := []struct {
		name    string
		rawData string
	}{
		{"parse stock", "cabinet:1"},
		{"parse optimize", "optimize:(time;cabinet)"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			configData := &ConfigData{
				Stocks: map[string]int{},
			}
			err := parseLine(configData, tc.rawData)
			if err != nil {
				t.Errorf("%v", err)
			}
		})
	}
	t.Run("parse stock", func(t *testing.T) {
		ConfigData := &ConfigData{
			Stocks: map[string]int{},
		}
		rawData := "cabinet:1"

		err := parseLine(ConfigData, rawData)
		if err != nil {
			t.Errorf("%v", err)
		}
	})
}
