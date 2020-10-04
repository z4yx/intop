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

func (s *IRQStat) CalcSum() {
	for num, info := range s.IRQSources {
		info.Sum = 0
		for i, n := range info.PerCPU {
			s.CPUSum[i] += n
			info.Sum += n
		}
		s.IRQSources[num] = info
	}

}

func (s *IRQStat) Subtract(d *IRQStat) {
	name2idx := map[string]int{}
	for i, name := range d.CPUName {
		name2idx[name] = i
	}
	for i, name := range s.CPUName {
		j, ok := name2idx[name]
		if !ok {
			continue
		}
		for num, _ := range s.IRQSources {
			dstat, ok := d.IRQSources[num]
			if !ok {
				continue
			}
			s.IRQSources[num].PerCPU[i] -= dstat.PerCPU[j]
		}
	}
	s.CalcSum()
}
