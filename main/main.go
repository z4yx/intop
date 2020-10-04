package main

import (
	"time"

	"github.com/z4yx/intop/datasource"
)

func main() {
	ui, err := NewIntopUI()
	defer EndIntopUI()

	if err != nil {
		return
	}

	startTime := time.Now()
	for {
		curr, err := datasource.GetCurrentIRQStat()
		ui.DrawTime(time.Since(startTime).Seconds())
		if err == nil {
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
		}
	}

}
