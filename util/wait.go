package util

import (
	"fmt"
	"strconv"
)

// ParseDuration parses the input string as an integer representing a duration in cycles.
// It returns the parsed integer value and an error if the input is not a valid non-negative integer.
// If the input is invalid or negative, an error describing the issue is returned.
func ParseDuration(input string) (int, error) {
	cycles, err := strconv.Atoi(input)
	if err != nil || cycles < 0 {
		return 0, fmt.Errorf("invalid duration: %s", input)
	}
	return cycles, nil
}

