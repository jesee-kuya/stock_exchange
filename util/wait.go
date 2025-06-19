package util

import (
	"fmt"
	"strconv"
)

// ParseDuration parses the input string as an float64 representing a duration in cycles.
// It returns the parsed float value and an error if the input is not a valid non-negative number.
// If the input is invalid or negative, an error describing the issue is returned.
func ParseDuration(input string) (float64, error) {
	cycles, err := strconv.ParseFloat(input, 64)
	if err != nil || cycles < 0 {
		return 0, fmt.Errorf("invalid duration: %s", input)
	}
	return cycles, nil
}
