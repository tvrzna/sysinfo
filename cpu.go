package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

const (
	pathStat = "/proc/stat"
	cpu      = "cpu"
)

type Cpu struct {
	id         *int
	user       int64
	nice       int64
	system     int64
	idle       int64
	iowait     int64
	irq        int64
	softirq    int64
	steal      int64
	guest      int64
	guest_nice int64
	cores      []*Cpu
}

func (c *Cpu) Usage(previous *Cpu) int {
	if previous != nil {
		total := c.user + c.nice + c.system + c.idle
		totalPrevious := previous.user + previous.nice + previous.system + previous.idle
		return int(100 * (total - totalPrevious + previous.idle - c.idle) / (total - totalPrevious))
	}
	return 0
}

func LoadCpu() *Cpu {
	f, err := os.Open(pathStat)
	if err != nil {
		return nil
	}
	defer f.Close()

	var result *Cpu

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, cpu) {
			cpu := &Cpu{}
			data := strings.Split(strings.ReplaceAll(line, "  ", " "), " ")
			if len(data[0]) == 3 {
				result = cpu
				result.cores = make([]*Cpu, 0)
			} else {
				result.cores = append(result.cores, cpu)
				id, _ := strconv.Atoi(data[0][3:])
				cpu.id = &id
			}
			cpu.user, _ = strconv.ParseInt(data[1], 10, 64)
			cpu.nice, _ = strconv.ParseInt(data[2], 10, 64)
			cpu.system, _ = strconv.ParseInt(data[3], 10, 64)
			cpu.idle, _ = strconv.ParseInt(data[4], 10, 64)
			cpu.iowait, _ = strconv.ParseInt(data[5], 10, 64)
			cpu.irq, _ = strconv.ParseInt(data[6], 10, 64)
			cpu.softirq, _ = strconv.ParseInt(data[7], 10, 64)
			cpu.steal, _ = strconv.ParseInt(data[8], 10, 64)
			cpu.guest, _ = strconv.ParseInt(data[9], 10, 64)
			cpu.guest_nice, _ = strconv.ParseInt(data[10], 10, 64)
		}
	}
	return result
}
