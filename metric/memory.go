package metric

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

const (
	pathMeminfo  = "/proc/meminfo"
	memTotal     = "MemTotal"
	memFree      = "MemFree"
	memAvailable = "MemAvailable"
	swapTotal    = "SwapTotal"
	swapFree     = "SwapFree"
)

type MemoryInfo struct {
	MemTotal     uint64
	MemFree      uint64
	MemAvailable uint64
	SwapTotal    uint64
	SwapFree     uint64
}

func LoadMemInfo() *MemoryInfo {
	f, err := os.Open(pathMeminfo)
	if err != nil {
		return nil
	}
	defer f.Close()

	result := &MemoryInfo{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		delimiter := strings.Index(line, ":")
		if delimiter < 0 {
			continue
		}

		item := line[:delimiter]
		data := strings.TrimSpace(strings.TrimRight(line[delimiter+1:], "kB"))
		value, err := strconv.ParseUint(data, 10, 64)
		if err != nil {
			continue
		}

		switch item {
		case memTotal:
			result.MemTotal = value
		case memFree:
			result.MemFree = value
		case memAvailable:
			result.MemAvailable = value
		case swapTotal:
			result.SwapTotal = value
		case swapFree:
			result.SwapFree = value
		}
	}
	return result
}
