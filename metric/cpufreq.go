package metric

import (
	"bufio"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	pathCpuinfo  = "/proc/cpuinfo"
	pathCpuSpeed = "/sys/devices/system/cpu/cpu0/cpufreq"
	processor    = "processor"
	cpuMHz       = "cpu MHz"
	coreId       = "core id"
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

	if len(result) > 0 && result[0].Freq == 0 {
		freq := getFallbackCpuFreq()
		if freq > 0 {
			for _, c := range result {
				c.Freq = freq
			}
		}
	}

	return result
}

func getFallbackCpuFreq() float64 {
	for _, f := range []string{"cpu_cur_freq", "scaling_cur_freq", "bios_limit"} {
		freq, err := loadCpuFreq(f)
		if err == nil && freq > 0 {
			return freq / 1000
		}
	}
	return 0
}

func loadCpuFreq(fileName string) (float64, error) {
	b, err := os.ReadFile(path.Join(pathCpuSpeed, fileName))
	if err != nil {
		return 0, err
	}
	value, err := strconv.ParseFloat(strings.TrimSpace(string(b)), 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}
