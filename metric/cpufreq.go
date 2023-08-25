package metric

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

const (
	pathCpuinfo = "/proc/cpuinfo"
	processor   = "processor"
	cpuMHz      = "cpu MHz"
	coreId      = "core id"
)

type Cpufreq struct {
	Processor int
	CoreId    int
	Freq      float64
}

func LoadCpufreq() []*Cpufreq {
	f, err := os.Open(pathCpuinfo)
	if err != nil {
		return nil
	}
	defer f.Close()

	var result []*Cpufreq

	var p *Cpufreq
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		delimiter := strings.Index(line, ":")
		if delimiter < 0 {
			continue
		}

		item := strings.TrimSpace(line[:delimiter])
		data := strings.TrimSpace(line[delimiter+1:])

		switch item {
		case processor:
			if p != nil {
				result = append(result, p)
			}
			p = &Cpufreq{}
			p.Processor, _ = strconv.Atoi(data)
		case cpuMHz:
			p.Freq, _ = strconv.ParseFloat(data, 64)
		case coreId:
			p.CoreId, _ = strconv.Atoi(data)
		}
	}
	if p != nil {
		result = append(result, p)
	}
	return result
}
