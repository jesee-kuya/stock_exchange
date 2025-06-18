package engine

import (
	"fmt"
	"time"

	"github.com/jesee-kuya/stock_exchange/process"
	"github.com/jesee-kuya/stock_exchange/util"
)

type runningProcess struct {
	Process *process.Process
	Delay   int
}

func (e *Engine) Run(waitingTime string) {
	maxSeconds, err := util.ParseDuration(waitingTime)
	if err != nil {
		strParseError := fmt.Sprintf("Invalid waiting time format: %v", err)
		e.Schedule = append(e.Schedule, strParseError)
		fmt.Println(strParseError)
		return
	}

	startTime := time.Now()

	fmt.Printf("Main Processes:\n")

	running := []runningProcess{}
	e.Schedule = []string{}
	e.Cycle = 0

	for {
		if time.Since(startTime).Seconds() >= float64(maxSeconds) {
			strTimeFailed := fmt.Sprintf("Time limit exceeded after %d cycles", e.Cycle)
			e.Schedule = append(e.Schedule, strTimeFailed)
			fmt.Println(strTimeFailed)
			break
		}

		running = updatedRunningProcesses(running, *e)
		e.Schedule = append(e.Schedule, scheduler(&running, *e)...)

		if len(running) == 0 {

			canStartAny := false
			for _, p := range e.Processes {
				if p.CanRun(e.Stock.Items) {
					canStartAny = true
					break
				}
			}
			if !canStartAny {
				strProcessFound := fmt.Sprintf("No more process doable at cycle %d", e.Cycle)
				e.Schedule = append(e.Schedule, strProcessFound)
				fmt.Println(strProcessFound)
				break
			}
		}

		e.Cycle++
	}

	printStock(e.Stock)
}


func updatedRunningProcesses(running []runningProcess, e Engine) []runningProcess {
	var updated []runningProcess
	for _, rp := range running {
		rp.Delay--
		if rp.Delay == 0 {
			for item, qty := range rp.Process.Result {
				e.Stock.Items[item] += qty
			}
		} else {
			updated = append(updated, rp)
		}
	}
	return updated
}

func scheduler(running *[]runningProcess, e Engine) []string {
	schedule := []string{}
	newProcessesScheduled := true

	for newProcessesScheduled {
		newProcessesScheduled = false
		for _, p := range e.Processes {
			if p.CanRun(e.Stock.Items) {

				for item, qty := range p.Needs {
					e.Stock.Items[item] -= qty
				}

				*running = append(*running, runningProcess{
					Process: p,
					Delay:   p.Cycle,
				})
				procesDescr := fmt.Sprintf(" %d:%s", e.Cycle, p.Name)
				schedule = append(schedule, procesDescr)
				fmt.Println(procesDescr)
				newProcessesScheduled = true
			}
		}
	}

	return schedule
}

func printStock(stock *Stock) {
	fmt.Println("Stock:")
	for item, qty := range stock.Items {
		fmt.Printf(" %s => %d\n", item, qty)
	}
}
