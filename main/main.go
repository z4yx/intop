package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akamensky/argparse"
	"github.com/z4yx/intop/datasource"
)

type options struct {
	excludedIRQ map[int]bool
	filteredCPU map[int]bool
	labelWidth  int
}

func parseOptions() (options, error) {
	var opt options
	parser := argparse.NewParser("intop", "Real-time visualization of /proc/interrupts")
	excludedIrqList := parser.IntList("e", "exclude", &argparse.Options{
		Required: false,
		Help:     "exclude the given IRQ number",
	})
	filteredCPUList := parser.IntList("c", "cpu", &argparse.Options{
		Required: false,
		Help:     "display specified CPUs only",
	})
	labelWidth := parser.Int("w", "width", &argparse.Options{
		Required: false,
		Default:  -1,
		Help:     "width of IRQ name column",
	})
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		return opt, err
	}

	opt.excludedIRQ = map[int]bool{}
	for _, v := range *excludedIrqList {
		opt.excludedIRQ[v] = true
	}
	if len(*filteredCPUList) > 0 {
		opt.filteredCPU = map[int]bool{}
		for _, v := range *filteredCPUList {
			opt.filteredCPU[v] = true
		}

	}
	opt.labelWidth = *labelWidth
	return opt, nil
}

func filterCPUs(opt *options, cpus []int) []int {
	if opt.filteredCPU == nil {
		return cpus
	}
	out := make([]int, 0, len(cpus))
	for _, c := range cpus {
		if _, exist := opt.filteredCPU[c]; exist {
			out = append(out, c)
		}
	}
	return out
}

func setLabelWidth(ui *IntopUI, opt *options, baseStat *datasource.IRQStat) {
	num, name := ui.CalcAdaptiveLabelWidth(&baseStat.IRQSources)
	if opt.labelWidth <= 0 {
		ui.SetIRQLabelWidth(num, name)
	} else {
		ui.SetIRQLabelWidth(num, opt.labelWidth)
	}
}

func main() {
	opt, err := parseOptions()
	if err != nil {
		return
	}

	ui, err := NewIntopUI()
	defer EndIntopUI()
	if err != nil {
		return
	}

	sigTrap := make(chan os.Signal)
	signal.Notify(sigTrap, os.Interrupt, syscall.SIGTERM)
	go (func() {
		<-sigTrap
		EndIntopUI()
		os.Exit(0)
	})()

	var baseStat, recentStat datasource.IRQStat
	resetOldStats := func() bool {
		var err error
		baseStat, err = datasource.GetCurrentIRQStat(opt.excludedIRQ)
		if err != nil {
			fmt.Printf("Failed to retrieve interrupts info: %v\n", err)
			return false
		}
		recentStat = baseStat.Clone()
		return true
	}

	if !resetOldStats() {
		return
	}
	setLabelWidth(ui, &opt, &baseStat)

	for {
		ui.DrawTime(recentStat.AcqTime.Sub(baseStat.AcqTime).Seconds())

		baseStat.RemoveStaleItems(&recentStat)
		recentStat.Subtract(&baseStat)

		cpusInOrder := recentStat.CalcCPURanking()
		cpusInOrder = filterCPUs(&opt, cpusInOrder)
		ui.SetCPUOrders(cpusInOrder)
		ui.DrawHeaderLines(recentStat.CPUName, recentStat.CPUSum)
		irqOrders := recentStat.CalcIRQSrcRanking()
		for i, num := range irqOrders {
			info := recentStat.IRQSources[num]
			ui.DrawIRQSources(i, info.Name, num, info.PerCPU)
		}
		ui.DrawFooterLine(len(irqOrders))
		ui.Refresh()

		waitMs := 1000 - time.Since(baseStat.AcqTime).Milliseconds()%1000
		key := ui.KeyInput(int(waitMs))
		if key == 'q' || key == 'Q' {
			break
		} else if key == 'r' || key == 'R' {
			resetOldStats()
		} else if key == 'd' { // for debugging
			// delete(recentStat.IRQSources, 131)
			// baseStat.RemoveStaleItems(&recentStat)
		} else {
			// Update statistics
			curr, err := datasource.GetCurrentIRQStat(opt.excludedIRQ)
			if err == nil {
				recentStat = curr
			}
		}
	}

}
