package metric

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
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
	Cores      []*Cpu
	Usage      float32
}

func LoadCpu(doneCh chan bool, bundle *Bundle) {
	previous := loadCpu()
	time.Sleep(200 * time.Millisecond)
	bundle.Cpu = loadCpu()
	for i := 0; i < len(bundle.Cpu.Cores); i++ {
		bundle.Cpu.Cores[i].calcUsage(previous.Cores[i])
	}
	doneCh <- true
}

func (c *Cpu) calcUsage(previous *Cpu) {
	if previous != nil {
		total := c.getTotal()
		totalPrevious := previous.getTotal()

		val := (total - totalPrevious + previous.idle - c.idle)
		valDiv := (total - totalPrevious)

		if val > 0 && valDiv > 0 {
			c.Usage = float32(100 * (total - totalPrevious + previous.idle - c.idle) / (total - totalPrevious))
			return
		}
	}
}

func (c *Cpu) getTotal() int64 {
	return c.user + c.nice + c.system + c.idle
}

func loadCpu() *Cpu {
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
				result.Cores = make([]*Cpu, 0)
			} else {
				result.Cores = append(result.Cores, cpu)
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
