package datasource

type IRQSource struct {
	Name   string
	Sum    uint64
	PerCPU []uint64
}

type IRQStat struct {
	CPUName    []string
	CPUSum     []uint64
	IRQSources map[int]IRQSource
}
