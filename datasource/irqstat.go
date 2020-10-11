package datasource

import "sort"

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

type sortPair struct {
	key   int
	value uint64
}

func ranking(lst []sortPair) (ret []int) {
	sort.Slice(lst, func(i, j int) bool {
		return lst[i].value > lst[j].value ||
			(lst[i].value == lst[j].value && lst[i].key < lst[j].key)
	})
	ret = make([]int, len(lst))
	for i, item := range lst {
		ret[i] = item.key
	}
	return
}

func (s *IRQStat) CalcCPURanking() []int {
	lst := make([]sortPair, len(s.CPUSum))
	for i, stat := range s.CPUSum {
		lst[i] = sortPair{i, stat}
	}
	return ranking(lst)
}

func (s *IRQStat) CalcIRQSrcRanking() []int {
	lst := make([]sortPair, len(s.IRQSources))
	i := 0
	for num, info := range s.IRQSources {
		lst[i] = sortPair{num, info.Sum}
		i++
	}
	return ranking(lst)
}

func (s *IRQStat) CalcSum() {
	for i, _ := range s.CPUSum {
		s.CPUSum[i] = 0
	}
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
		if !ok { // CPU not found in d
			continue
		}
		for num, _ := range s.IRQSources {
			dstat, ok := d.IRQSources[num]
			if !ok { // IRQ source not found in d
				continue
			}
			if s.IRQSources[num].PerCPU[i] >= dstat.PerCPU[j] {
				s.IRQSources[num].PerCPU[i] -= dstat.PerCPU[j]
			}
		}
	}
	s.CalcSum()
}

func (s *IRQStat) RemoveStaleItems(n *IRQStat) {
	rmList := make([]int, 0)
	for num, _ := range s.IRQSources {
		if _, ok := n.IRQSources[num]; !ok {
			rmList = append(rmList, num)
		}
	}
	for _, num := range rmList {
		delete(s.IRQSources, num)
	}
}
