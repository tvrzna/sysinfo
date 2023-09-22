package metric

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
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

func LoadCpu(wg *sync.WaitGroup, bundle *Bundle) {
	previous := loadCpu()
	time.Sleep(500 * time.Millisecond)
	bundle.Cpu = loadCpu()
	for i := 0; i < len(bundle.Cpu.Cores); i++ {
		bundle.Cpu.Cores[i].calcUsage(previous.Cores[i])
	}
	wg.Done()
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
			var name string
			fmt.Sscanf(line, "%s %d %d %d %d %d %d %d %d %d %d", &name, &cpu.user, &cpu.nice, &cpu.system, &cpu.idle, &cpu.iowait, &cpu.irq, &cpu.softirq, &cpu.steal, &cpu.guest, &cpu.guest_nice)
			if len(name) == 3 {
				result = cpu
				result.Cores = make([]*Cpu, 0)
			} else {
				result.Cores = append(result.Cores, cpu)
				id, _ := strconv.Atoi(name[3:])
				cpu.id = &id
			}
		}
	}
	return result
}
