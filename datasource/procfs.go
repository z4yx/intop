package datasource

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

var reIRQName *regexp.Regexp

func parseCPUNames(header string) (CPUName []string) {
	rdr := strings.NewReader(header)
	CPUName = make([]string, 0)
	for {
		var name string
		n, _ := fmt.Fscanf(rdr, "%s", &name)
		if n == 0 {
			break
		}
		CPUName = append(CPUName, name)
	}
	return
}

func parseIRQSource(line string, nCPU int) (ok bool, number int, info IRQSource) {
	ok = false
	rdr := strings.NewReader(line)
	n, err := fmt.Fscanf(rdr, "%d:", &number)
	if n == 0 || err != nil {
		return
	}
	info.PerCPU = make([]uint64, nCPU)
	for i := 0; i < nCPU; i++ {
		n, err = fmt.Fscanf(rdr, "%d", &info.PerCPU[i])
		if n == 0 || err != nil {
			return
		}
	}
	remaining, err := ioutil.ReadAll(rdr)
	if err != nil {
		return
	}
	matches := reIRQName.FindAllSubmatch(remaining, -1)
	if len(matches) != 1 || len(matches[0]) != 2 {
		return
	}
	info.Name = string(matches[0][1])
	ok = true
	return
}

func GetCurrentIRQStat() (stat IRQStat, err error) {
	interrupts, err := ioutil.ReadFile("/proc/interrupts")
	if err != nil {
		return
	}
	lines := strings.Split(string(interrupts), "\n")
	if len(lines) < 1 {
		err = errors.New("Empty file")
		return
	}
	stat.CPUName = parseCPUNames(lines[0])
	stat.CPUSum = make([]uint64, len(stat.CPUName))
	stat.IRQSources = map[int]IRQSource{}
	if reIRQName == nil {
		reIRQName = regexp.MustCompilePOSIX(`.+  (.+)$`) // two spaces as the delimiter
	}
	for i := 1; i < len(lines); i++ {
		ok, num, info := parseIRQSource(lines[i], len(stat.CPUName))
		if !ok {
			break
		}
		stat.IRQSources[num] = info
	}
	stat.CalcSum()
	return
}
