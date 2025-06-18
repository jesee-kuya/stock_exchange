# Stock Exchange Simulator

A process chain optimization program that maximizes performance while minimizing delays. This project implements a scheduling system that analyzes task dependencies and resource constraints to find optimal execution sequences.

## Project Overview

The Stock Exchange Simulator is designed to solve resource allocation and scheduling problems by optimizing process execution order. The system reads configuration files that define available resources (stocks) and processes that consume and produce these resources, then generates an optimized execution schedule.

### Key Capabilities
- **Process Chain Optimization**: Analyzes task dependencies and resource constraints
- **Multiple Scheduling Schemes**: Supports both serial and parallel execution strategies
- **Resource Management**: Tracks stock levels and prevents resource conflicts
- **Flexible Optimization**: Can optimize for time, specific products, or combinations
- **Schedule Validation**: Includes verification system to ensure generated schedules are valid

### Use Cases
- Manufacturing process optimization
- Project task scheduling with resource constraints
- Supply chain management
- Any scenario with interdependent tasks and limited resources

## Building and Running

### Prerequisites
- Go 1.19 or higher
- Git (for cloning the repository)

### Building the Project

Clone the repository and build both programs:

```bash
git clone <repository_url>
cd stock-scheduler

# Build the main scheduler
go build -o stock_exchange .

# Build the checker
go build -o checker ./checker
```

### Running the Scheduler

Execute the scheduler with a configuration file and timeout:

```bash
./stock_exchange <config_file> <timeout_seconds>
```

**Parameters:**
- `config_file`: Path to configuration file defining processes and stocks
- `timeout_seconds`: Maximum execution time to prevent infinite loops

**Example:**
```bash
./stock_exchange examples/cabinet_build.txt 30
```

### Running the Checker

Validate a generated schedule against the original configuration:

```bash
./checker <config_file> <log_file>
```

**Parameters:**
- `config_file`: Original configuration file
- `log_file`: Generated schedule log to validate

**Example:**
```bash
./checker examples/cabinet_build.txt examples/cabinet_build.log
```

### Running Both Programs Together

For a complete workflow, run the scheduler followed by the checker:

```bash
# Run scheduler and generate log
./stock_exchange examples/cabinet_build.txt 30

# Validate the generated schedule
./checker examples/cabinet_build.txt examples/cabinet_build.log
```

You can also chain them in a single command:

```bash
./stock_exchange examples/cabinet_build.txt 30 && ./checker examples/cabinet_build.txt examples/cabinet_build.log
```

## File Formats

### Configuration File Format

Configuration files define the initial state and processes using this syntax:

```
# Comments start with hash symbol

# Initial stock definitions
<stock_name>:<quantity>

# Process definitions  
<process_name>:(<input_stock>:<quantity>;...):(output_stock>:<quantity>;...):<cycle_duration>

# Optimization targets
optimize:(<target_stock>|time)
```

#### Configuration Example

```
# Cabinet manufacturing example
board:7

# Process definitions
do_doorknobs:(board:1):(doorknobs:1):15
do_background:(board:2):(background:1):20
do_shelf:(board:1):(shelf:1):10
do_cabinet:(doorknobs:2;background:1;shelf:3):(cabinet:1):30

# Optimize for time and cabinet production
optimize:(time;cabinet)
```

#### Configuration Rules
- Comments begin with `#` and are ignored
- Stock definitions must come before process definitions
- Process inputs and outputs are separated by colons and semicolons
- The optimize line specifies what to maximize (use `time` for time optimization)

### Log File Format

The scheduler generates log files with execution traces:

```
<cycle>:<process_name>
<cycle>:<process_name>
...
No more process doable at cycle <final_cycle>
```

#### Log Example

```
0:do_shelf
0:do_shelf  
0:do_shelf
0:do_doorknobs
0:do_doorknobs
0:do_background
20:do_cabinet
No more process doable at cycle 51
```

## Example Usage and Output

### Complete Workflow Example

**Configuration file** (`examples/simple.txt`):
```
# Simple production line
euro:10

buy_materiel:(euro:8):(material:1):10
build_product:(material:1):(product:1):30
delivery:(product:1):(client_content:1):20

optimize:(time;client_content)
```

**Running the scheduler:**
```bash
$ ./stock_exchange examples/simple.txt 60
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

**Generated log file** (`examples/simple.log`):
```
0:buy_materiel
10:build_product
40:delivery
No more process doable at cycle 61
```

**Running the checker:**
```bash
$ ./checker examples/simple.txt examples/simple.log
Evaluating: 0:buy_materiel
Evaluating: 10:build_product
Evaluating: 40:delivery
Trace completed, no error detected.
```

### Error Handling Example

**Invalid log file:**
```
0:buy_materiel
10:build_product
10:build_product  # This creates an error - insufficient materials
40:delivery
```

**Checker output:**
```bash
$ ./checker examples/simple.txt examples/simple_error.log
Evaluating: 0:buy_materiel
Evaluating: 10:build_product
Evaluating: 10:build_product
Error detected
at 10:build_product stock insufficient
```

## Project Structure

```
stock-scheduler/
├── go.mod                      # Go module definition
├── main.go                     # Main entry point
├── engine/                     # Scheduling engine
│   ├── engine.go               # Core engine logic
│   ├── load_config.go          # Configuration loading
│   ├── run.go                  # Execution logic
│   └── save_log.go             # Log file generation
├── process/                    # Process management
│   ├── process.go              # Process definitions
│   ├── can_run.go              # Resource validation
│   └── run.go                  # Process execution
├── checker/                    # Validation system
│   ├── checker.go              # Checker implementation
│   ├── load_config.go          # Config loading for validation
│   ├── load_log.go             # Log file parsing
│   └── verify.go               # Schedule verification
├── examples/                   # Sample configurations
│   ├── cabinet_build.txt       # Cabinet manufacturing example
│   ├── simple.txt              # Basic production line
│   └── infinite.txt            # Self-sustaining processes
└── README.md                   # This documentation
```

## Contributing

We welcome contributions to improve the Stock Exchange Simulator! Please follow these guidelines to ensure smooth collaboration.

### Collaboration Rules

1. **Fork and Branch**: Create a fork of the repository and work on feature branches
2. **Code Quality**: Ensure your code follows Go conventions and includes appropriate comments
3. **Testing**: Add tests for new functionality and ensure existing tests pass
4. **Documentation**: Update documentation for any user-facing changes
5. **Small Commits**: Make focused commits that address single concerns
6. **Pull Request Reviews**: All changes must be reviewed before merging

### Commit Message Format

Use the following prefixes for commit messages to maintain consistency:

| Prefix      | Purpose                       | Example |
| ----------- | ----------------------------- | ------- |
| `feat:`     | New features or functionality | `feat: add parallel scheduling algorithm` |
| `fix:`      | Bug fixes                     | `fix: resolve memory leak in process execution` |
| `refactor:` | Code cleanup or restructuring | `refactor: simplify configuration parser` |
| `docs:`     | Documentation updates         | `docs: update README with new examples` |
| `test:`     | Adding or fixing tests        | `test: add unit tests for resource validation` |
| `chore:`    | Miscellaneous changes         | `chore: update Go module dependencies` |

### Development Workflow

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/yourusername/stock-scheduler.git
   cd stock-scheduler
   ```
3. **Create a feature branch**:
   ```bash
   git checkout -b feat/your-feature-name
   ```
4. **Make your changes** and commit with proper format:
   ```bash
   git add .
   git commit -m "feat: add your new feature description"
   ```
5. **Run tests** to ensure nothing is broken:
   ```bash
   go test ./...
   ```
6. **Push to your fork**:
   ```bash
   git push origin feat/your-feature-name
   ```
7. **Create a pull request** with a clear description of your changes

### Code Style Guidelines

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Document public functions with comments starting with the function name
- Keep functions focused and relatively small
- Handle errors appropriately with clear error messages

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributors

- [jesee-kuya](https://github.com/jesee-kuya/)
- [joseowino](https://github.com/joseowino)
- [stkisengese](https://github.com/stkisengese)
- [Baraq23](https://github.com/Baraq23)