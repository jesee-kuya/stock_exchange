package util

import (
	"fmt"
	"strconv"
)

// ParseDuration parses a string representing a number of cycles (e.g., "10")
// and returns it as an integer. Returns an error if the input is invalid.
func ParseDuration(input string) (int, error) {
	cycles, err := strconv.Atoi(input)
	if err != nil || cycles < 0 {
		return 0, fmt.Errorf("invalid duration: %s", input)
	}
	return cycles, nil
}

