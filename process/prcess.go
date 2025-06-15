package process

type Process struct {
	Name         string
	Needs        map[string]int
	Result       map[string]int
	Cycle        int
}


