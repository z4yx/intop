package main

import (
	"fmt"

	"github.com/gbin/goncurses"
)

const IRQ_NAME_WIDTH = 16

type IntopUI struct {
	win      *goncurses.Window
	orderCPU []int
	// orderIRQ []int
}

func (ui *IntopUI) SetCPUOrders(orderCPU []int) {
	ui.orderCPU = orderCPU
	// ui.orderIRQ = orderIRQ
}

func (ui *IntopUI) DrawHeaderLines(names []string, irqSum []uint64) {
	ui.win.Move(0, IRQ_NAME_WIDTH)
	for _, idx := range ui.orderCPU {
		ui.win.Printf("% 10s", names[idx])
	}
	ui.win.Move(1, 0)
	ui.win.Printf("%*s", IRQ_NAME_WIDTH, "IRQ Per CPU:")
	for _, idx := range ui.orderCPU {
		ui.win.Printf("% 10d", irqSum[idx])
	}
}

func (ui *IntopUI) DrawFooterLine(nIRQ int) {
	ui.win.Move(nIRQ+2, 0)
	ui.win.ClearToEOL()
	ui.win.Print("  [q] Quit  [r] Reset Statistics")
}

func (ui *IntopUI) DrawIRQSources(index int, name string, number int, irqPerCPU []uint64) {
	ui.win.Move(index+2, 0)
	label := fmt.Sprintf("%d ", number)
	width := IRQ_NAME_WIDTH - len(label)
	ui.win.AttrOn(goncurses.A_BOLD)
	ui.win.Print(label)
	ui.win.AttrOff(goncurses.A_BOLD)
	if len(name) > width {
		name = name[:width]
	}
	ui.win.Printf("%-*s", width, name)
	for _, idx := range ui.orderCPU {
		ui.win.Printf("% 10d", irqPerCPU[idx])
	}
}

func (ui *IntopUI) DrawTime(t float64) {
	text := fmt.Sprintf("T=%.3fs", t)
	ui.win.MovePrintf(0, 0, "%*s", IRQ_NAME_WIDTH, text)
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
	// ui.win.HLine(1, 1, goncurses.ACS_BOARD, 60)
	// ui.win.HLine(30, 1, goncurses.ACS_BOARD, 60)
	return
}

func EndIntopUI() {
	goncurses.End()
}
