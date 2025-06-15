package process

// Process represents a processing unit in the stock exchange system.
// It defines the name of the process, the required input resources (Needs),
// the output resources produced (Result), and the number of cycles needed to complete the process.
type Process struct {
	Name         string
	Needs        map[string]int
	Result       map[string]int
	Cycle        int
}


