package engine

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jesee-kuya/stock_exchange/process"
)

// Engine is the main structure for executing and optimizing
// the scheduling process defined in a configuration file
type Engine struct {
	Stock           *Stock
	Processes       []*process.Process
	Schedule        []string
	Cycle           int
	OptimizeTargets []string
}

// Stock represents the available items and their quantities.
type Stock struct {
	Items map[string]int
}

func (e *Engine) LoadConfig(path string) error {
	file, err := os.Open(path)

	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

	}

	return nil
}
