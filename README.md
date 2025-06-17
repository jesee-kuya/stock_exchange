# Stock Exchange Simulator

A process chain optimization program that maximizes performance while minimizing delays. This project implements a scheduling system that analyzes task dependencies and resource constraints to find optimal execution sequences.

## Overview

The Stock Exchange Simulator reads process definitions from configuration files, analyzes dependencies and resource requirements, then generates optimized execution schedules. The system includes both a main simulator and a checker program to validate generated schedules.

## Features

- **Process Chain Optimization**: Analyzes task dependencies and resource constraints
- **Flexible Scheduling**: Supports both serial and parallel schedule generation schemes
- **Resource Management**: Tracks stock levels and prevents resource conflicts
- **Multiple Optimization Targets**: Can optimize for time, specific products, or combinations
- **Schedule Validation**: Includes checker program to verify generated schedules
- **Configurable Timeout**: Prevents infinite loops with user-defined time limits

## Architecture

### Main Components

1. **Stock Exchange Program**: Generates optimized schedules from configuration files
2. **Checker Program**: Validates generated schedules for correctness
3. **Configuration Parser**: Reads and validates process definitions
4. **Scheduler Engine**: Implements optimization algorithms
5. **Resource Tracker**: Manages stock levels throughout execution

### Scheduling Schemes

- **Serial Schedule Generation**: Selects activities sequentially and schedules them as-soon-as-possible
- **Parallel Schedule Generation**: Schedules multiple activities simultaneously when resources allow

## File Format

### Configuration File Structure

```
# Comments start with #

# Initial stock definitions
<stock_name>:<quantity>

# Process definitions
<process_name>:(<input_stock>:<quantity>;...):(output_stock>:<quantity>;...):<cycle_duration>

# Optimization targets
optimize:(<target_stock>|time)
```

### Example Configuration

```
# Cabinet building example
board:7

do_doorknobs:(board:1):(doorknobs:1):15
do_background:(board:2):(background:1):20
do_shelf:(board:1):(shelf:1):10
do_cabinet:(doorknobs:2;background:1;shelf:3):(cabinet:1):30

optimize:(time;cabinet)
```

## Usage

### Running the Stock Exchange Simulator

```bash
./stock_exchange <config_file> <timeout_seconds>
```

**Parameters:**
- `config_file`: Path to the configuration file defining processes and stocks
- `timeout_seconds`: Maximum execution time to prevent infinite loops

**Example:**
```bash
go run . examples/simple 10
```

**Output:**
```
Main Processes:
 0:buy_materiel
 10:build_product
 40:delivery
No more process doable at cycle 61

Stock:
 euro => 2
 material => 0
 product => 0
 client_content => 1
```

### Running the Checker

```bash
./checker <config_file> <log_file>
```

**Parameters:**
- `config_file`: Original configuration file
- `log_file`: Generated schedule log to validate

**Example:**
```bash
go run ./checker examples/simple examples/simple.log
```

**Output (Success):**
```
Evaluating: 0:buy_materiel
Evaluating: 10:build_product
Evaluating: 40:delivery
Trace completed, no error detected.
```

**Output (Error):**
```
Evaluating: 0:buy_materiel
Evaluating: 10:build_product
Evaluating: 10:build_product
Error detected
at 10:build_product stock insufficient
```

## Log File Format

The simulator generates log files with the following format:
```
<cycle>:<process_name>
<cycle>:<process_name>
...
No more process doable at cycle <final_cycle>
```

Example:
```
0:buy_materiel
10:build_product
40:delivery
No more process doable at cycle 61
```

## Implementation Requirements
-  Golang
- Cofigaration file

### Optimization Strategies
- **Time Optimization**: Minimize total execution time
- **Product Optimization**: Maximize specific product output
- **Combined Optimization**: Balance multiple objectives

## Example Scenarios

### Finite Process Chain
Resources are consumed until no more processes can execute:
```
board:7
do_shelf:(board:1):(shelf:1):10
optimize:(time;shelf)
```

### Infinite Process Chain
Self-sustaining processes that can run indefinitely:
```
money:100
buy_materials:(money:50):(materials:10):5
sell_products:(materials:5):(money:60):10
optimize:(time;money)
```

## Error Handling

The system handles various error conditions:
- **Parse Errors**: Invalid configuration file syntax
- **Resource Errors**: Insufficient stocks for process execution
- **Dependency Errors**: Missing prerequisites for processes
- **Timeout Errors**: Execution exceeds specified time limit

## Building and Installation

```bash
# Clone the repository
git clone <repository_url>
cd stock-exchange-sim

# Build the main program
go build -o stock_exchange .

# Build the checker
go build -o checker ./checker

# Run tests
go test ./...
```

## Project Structure

```
stock-scheduler/
├── go.mod
├── main.go                     # Entrypoint: decides whether to run scheduler or checker

├── engine/                     # Scheduling logic
│   ├── engine.go               # Engine struct definition
│   ├── load_config.go          # func (e *Engine) LoadConfig(path string) error
│   ├── run.go                  # func (e *Engine) Run(waitTime string)
│   ├── save_log.go             # func (e *Engine) SaveLog(path string) error

├── process/                    # Process logic
│   ├── process.go              # Process struct definition
│   ├── can_run.go              # func (p *Process) CanRun(stocks map[string]int) bool
│   ├── run.go                  # func (p *Process) Run(stocks map[string]int)

├── checker/                    # Checker program logic
│   ├── checker.go              # Checker struct definition
│   ├── new_checker.go          # func NewChecker() *Checker
│   ├── load_config.go          # func (c *Checker) LoadConfig(path string) error
│   ├── load_log.go             # func (c *Checker) LoadLog(path string) error
│   ├── verify.go               # func (c *Checker) Verify() error

├── examples/                   # Example config & log files
│   ├── cabinet_build.txt       # Sample config file
│   ├── cabinet_build.log       # Log file generated by scheduler

├── util/                       # Utility helpers (optional)
│   ├── parse.go                # Reusable parsing helpers
│   ├── file.go                 # File I/O utilities
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Implement your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributers
jesee-kuya/(https://github.com/jesee-kuya/)
joseowino(https://github.com/joseowino)
stkisengese(https://github.com/stkisengese)
Baraq23(https://github.com/Baraq23)
