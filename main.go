package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jesee-kuya/stock_exchange/checker"
	e "github.com/jesee-kuya/stock_exchange/engine"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  Schedule: go run . <config_file> <log_output_file> <wait_time>")
		fmt.Println("  Check:    go run . ./checker <config_file> <log_file>")
		os.Exit(1)
	}

	mode := args[1]

	switch mode {
	case "./checker":
		if len(args) != 4 {
			fmt.Println("Usage: go run . check <config_file> <log_file>")
			os.Exit(1)
		}

		configPath := args[2]
		logPath := args[3]

		chk := checker.NewChecker()

		if err := chk.LoadConfig(configPath); err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		if err := chk.LoadLog(logPath); err != nil {
			fmt.Printf("Error loading log: %v\n", err)
			os.Exit(1)
		}

		if err := chk.Verify(); err != nil {
			fmt.Printf("Verification failed: %v\n", err)
			os.Exit(1)
		}

	default:
		engine()
	}
}

func engine() {
	fmt.Println(len(os.Args))
	if len(os.Args) != 3 {
		log.Fatal("Usage: run <config_file> <waiting_time>")
		return
	}
	configFile := os.Args[1]
	waitTime := os.Args[2]

	engine := e.NewEngine()
	if err := engine.LoadConfig(configFile); err != nil {
		log.Fatal(err)
	}
	engine.Run(waitTime)
	if err := engine.SaveLog(configFile + ".log"); err != nil {
		log.Fatal(err)
	}
}
