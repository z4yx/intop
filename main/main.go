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
	resetOldStats := func() bool {
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

	if !resetOldStats() {
		return
	}
	for {
		ui.DrawTime(time.Since(baseTime).Seconds())
		curr, err := datasource.GetCurrentIRQStat()
		if err == nil {
			curr.Subtract(&baseStat)
			ui.SetCPUOrders(curr.CalcCPURanking())
			ui.DrawHeaderLines(curr.CPUName, curr.CPUSum)
			irqOrders := curr.CalcIRQSrcRanking()
			for i, num := range irqOrders {
				info := curr.IRQSources[num]
				ui.DrawIRQSources(i, info.Name, num, info.PerCPU)
			}
			ui.DrawFooterLine(len(irqOrders))
		}
		ui.Refresh()
		key := ui.KeyInput(1000)
		if key == 'q' || key == 'Q' {
			break
		} else if key == 'r' || key == 'R' {
			resetOldStats()
		}
	}

}
