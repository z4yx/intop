package main

import (
	"fmt"
	"time"

	"github.com/z4yx/intop/datasource"
)

func main() {
	ui, err := NewIntopUI()
	defer EndIntopUI()
	if err != nil {
		return
	}

	var baseStat datasource.IRQStat
	var baseTime time.Time
	clearOldStats := func() bool {
		var err error
		currTime := time.Now()
		baseStat, err = datasource.GetCurrentIRQStat()
		if err != nil {
			fmt.Printf("Failed to retrieve interrupts info: %v\n", err)
			return false
		}
		baseTime = currTime
		return true
	}

	if !clearOldStats() {
		return
	}
	for {
		ui.DrawTime(time.Since(baseTime).Seconds())
		curr, err := datasource.GetCurrentIRQStat()
		if err == nil {
			curr.Subtract(&baseStat)
			ui.DrawHeaderLines(curr.CPUName, curr.CPUSum)
			i := 0
			for num, info := range curr.IRQSources {
				ui.DrawIRQSource(i, info.Name, num, info.PerCPU)
				i++
			}
		}
		ui.Refresh()
		key := ui.KeyInput(1000)
		if key == 'q' || key == 'Q' {
			break
		} else if key == 'c' || key == 'C' {
			clearOldStats()
		}
	}

}
