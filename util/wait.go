package util

import "time"

func Wait(duration string) {
	d, err := time.ParseDuration(duration)
	if err != nil {
		return // Invalid duration, skip waiting
	}
	time.Sleep(d)
}
