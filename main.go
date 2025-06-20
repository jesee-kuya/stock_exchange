package main

import (
	"fmt"
	"log"
	"os"

	e "github.com/jesee-kuya/stock_exchange/engine"
)

// main is the entry point of the stock exchange application. It parses command-line arguments
// to determine the mode of operation: either running the engine or performing a log check.
// In "checker" mode, it loads configuration and log files, then verifies the log using the checker package.
// In the default mode, it runs the engine with the provided configuration file and wait time.
func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  Schedule: go run . <config_file> <log_output_file> <wait_time>")
		fmt.Println("  Check:    go run ./checker <config_file> <log_file>")
		return
	}

	engine()
}

// engine is responsible for running the stock exchange engine in scheduling mode.
// It expects exactly two command-line arguments: the configuration file path and the waiting time.
// The function performs the following steps:
//  1. Validates the number of arguments and prints usage instructions if incorrect.
//  2. Loads the engine configuration from the specified file.
//  3. Runs the engine with the provided waiting time.
//  4. Saves the engine's log to a file with the same name as the configuration file, appended with ".log".
//
// If any step fails, the function logs the error and terminates the program.
func engine() {
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
