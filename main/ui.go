package main

import (
	"fmt"
	"sort"

	"github.com/gbin/goncurses"
	"github.com/z4yx/intop/datasource"
)

type IntopUI struct {
	win          *goncurses.Window
	orderCPU     []int
	irqNumWidth  int
	irqNameWidth int
}

func (ui *IntopUI) CalcAdaptiveLabelWidth(sources *map[int]datasource.IRQSource) (numWidth int, nameWidth int) {
	const MIN_NAME_WIDTH = 10
	const FULL_NAME_RATIO = float64(0.8)
	lenList := make([]int, 0, len(*sources))

	for num, info := range *sources {
		u, a := ui.IRQLabelWidth(info.Name, num)
		lenList = append(lenList, a)
		if u > numWidth {
			numWidth = u
		}
	}
	sort.Ints(lenList)
	idx := FULL_NAME_RATIO * float64(len(lenList))
	nameWidth = lenList[int(idx)]
	if nameWidth < MIN_NAME_WIDTH {
		nameWidth = MIN_NAME_WIDTH
	}
	return
}

func (ui *IntopUI) SetCPUOrders(orderCPU []int) {
	ui.orderCPU = orderCPU
}

func (ui *IntopUI) SetIRQLabelWidth(numWidth int, nameWidth int) {
	ui.irqNumWidth = numWidth
	ui.irqNameWidth = nameWidth
}

func (ui *IntopUI) DrawHeaderLines(names []string, irqSum []uint64) {
	ui.win.Move(0, ui.irqNumWidth+ui.irqNameWidth)
	for _, idx := range ui.orderCPU {
		ui.win.Printf("% 10s", names[idx])
	}
	ui.win.Move(1, 0)
	ui.win.Printf("%*s", ui.irqNumWidth+ui.irqNameWidth, "IRQ Per CPU:")
	for _, idx := range ui.orderCPU {
		ui.win.Printf("% 10d", irqSum[idx])
	}
}

func (ui *IntopUI) DrawFooterLine(nIRQ int) {
	ui.win.Move(nIRQ+2, 0)
	ui.win.ClearToEOL()
	ui.win.Print("  [q] Quit  [r] Reset Count")
}

func (ui *IntopUI) IRQLabelWidth(name string, number int) (int, int) {
	irqNum := fmt.Sprintf("%d ", number)
	return len(irqNum), len(name)
}

func (ui *IntopUI) DrawIRQSources(index int, name string, number int, irqPerCPU []uint64) {
	ui.win.Move(index+2, 0)

	irqNum := fmt.Sprintf("%-*d", ui.irqNumWidth, number)
	ui.win.AttrOn(goncurses.A_BOLD)
	ui.win.Print(irqNum)
	ui.win.AttrOff(goncurses.A_BOLD)

	if len(name) > ui.irqNameWidth {
		name = name[:ui.irqNameWidth]
	}
	ui.win.Printf("%-*s", ui.irqNameWidth, name)

	for _, idx := range ui.orderCPU {
		ui.win.Printf("% 10d", irqPerCPU[idx])
	}
}

func (ui *IntopUI) DrawTime(t float64) {
	text := fmt.Sprintf("T=%.3fs", t)
	ui.win.MovePrintf(0, 0, "%*s", ui.irqNumWidth+ui.irqNameWidth, text)
}

func (ui *IntopUI) Refresh() {
	goncurses.Update()
	// ui.win.Refresh()
	// y, x := ui.win.MaxYX()
	// ui.win.MovePrintf(0, 0, "%d,%d", y, x)
}

func (ui *IntopUI) KeyInput(timeout int) int {
	ui.win.Timeout(timeout)
	return int(ui.win.GetChar())
}

func NewIntopUI() (ui *IntopUI, err error) {
	ui = new(IntopUI)
	ui.win, err = goncurses.Init()
	if err != nil {
		ui = nil
		return
	}
	goncurses.Cursor(0)
	return
}

func EndIntopUI() {
	goncurses.End()
}
