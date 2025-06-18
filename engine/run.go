package engine

import (
	"fmt"
	"time"

	"github.com/jesee-kuya/stock_exchange/process"
	"github.com/jesee-kuya/stock_exchange/util"
)

func (e *Engine) Run(waitingTime string) {
	maxSeconds, err := util.ParseDuration(waitingTime)
	if err != nil {
		fmt.Printf("Invalid waiting time: %v\n", err)
		return
	}

	startTime := time.Now()

	fmt.Printf("Main Processes:\n")

	type runningProcess struct {
		Process *process.Process
		Delay   int
	}

	var running []runningProcess
	e.Schedule = []string{}
	e.Cycle = 0

	for {
		if time.Since(startTime).Seconds() >= float64(maxSeconds) {
			fmt.Printf("Time limit exceeded after %d cycles\n", e.Cycle)
			break
		}

		var updatedRunning []runningProcess
		for _, rp := range running {
			rp.Delay--
			if rp.Delay == 0 {
				for item, qty := range rp.Process.Result {
					e.Stock.Items[item] += qty
				}
			} else {
				updatedRunning = append(updatedRunning, rp)
			}
		}
		running = updatedRunning

		newProcessesScheduled := true
		for newProcessesScheduled {
			newProcessesScheduled = false
			for _, p := range e.Processes {
				if p.CanRun(e.Stock.Items) {

					for item, qty := range p.Needs {
						e.Stock.Items[item] -= qty
					}

					running = append(running, runningProcess{
						Process: p,
						Delay:   p.Cycle,
					})
					e.Schedule = append(e.Schedule, p.Name)
					fmt.Printf(" %d:%s\n", e.Cycle, p.Name)
					newProcessesScheduled = true
				}
			}
		}

		if len(running) == 0 {

			canStartAny := false
			for _, p := range e.Processes {
				if p.CanRun(e.Stock.Items) {
					canStartAny = true
					break
				}
			}
			if !canStartAny {
				fmt.Printf("No more process doable at cycle %d\n", e.Cycle)
				break
			}
		}

		e.Cycle++
	}

	printStock(e.Stock)
}

func printStock(stock *Stock) {
	fmt.Println("Stock:")
	for item, qty := range stock.Items {
		fmt.Printf(" %s => %d\n", item, qty)
	}
}
