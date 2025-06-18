package engine

import (
	"fmt"
	"sort"
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
		fmt.Println("Invalid waiting time format:", err)
		return
	}

	start := time.Now()
	e.Schedule = []string{}
	e.Cycle = 0
	running := []runningProcess{}

	fmt.Println("Main Processes :")

	targets := map[string]bool{}
	for _, t := range e.OptimizeTargets {
		targets[t] = true
	}
	priorities := computePriorities(e.Processes, targets)

	for {
		if time.Since(start).Seconds() >= float64(maxSeconds) {
			fmt.Printf("Time limit exceeded after %d cycles\n", e.Cycle)
			break
		}

		running = updateRunningProcesses(running, e)

		if entry := scheduleOneProcess(&running, e, priorities); entry != "" {
			e.Schedule = append(e.Schedule, entry)
			fmt.Println(entry)
		}

		if len(running) == 0 && !e.canRunAny() {
			fmt.Printf("No more process doable at cycle %d\n", e.Cycle)
			break
		}

		e.Cycle++
	}

	printStock(e.Stock)
}

func updateRunningProcesses(running []runningProcess, e *Engine) []runningProcess {
	next := []runningProcess{}
	for _, rp := range running {
		rp.Delay--
		if rp.Delay <= 0 {
			for item, qty := range rp.Process.Result {
				e.Stock.Items[item] += qty
			}
		} else {
			next = append(next, rp)
		}
	}
	return next
}


