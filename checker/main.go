package main

import (
	"fmt"
	"os"

	checker "github.com/jesee-kuya/stock_exchange/checker_util"
)

func main() {
	args := os.Args
	if len(args) != 3 {
		fmt.Println("Usage: go run ./checker <config_file> <log_file>")
		return
	}

	configPath := args[1]
	logPath := args[2]

	chk := checker.NewChecker()

	if err := chk.LoadConfig(configPath); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	if err := chk.LoadLog(logPath); err != nil {
		fmt.Printf("Error loading log: %v\n", err)
		return
	}

	if err := chk.Verify(); err != nil {
		fmt.Printf("Verification failed: %v\n", err)
		return
	}
}
