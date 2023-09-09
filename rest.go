package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tvrzna/sysinfo/metric"
)

type MemoryUnit byte

const (
	UnitK MemoryUnit = iota
	UnitM
	UnitG
	UnitT
)

func (b MemoryUnit) String() string {
	return []string{"K", "M", "G", "T"}[int(b)]
}

func (b MemoryUnit) MarshalJSON() ([]byte, error) {
	return []byte("\"" + b.String() + "\""), nil
}

func (b *MemoryUnit) UnmarshalJSON(data []byte) error {
	var v string
	json.Unmarshal(data, &v)
	*b = map[string]MemoryUnit{"K": UnitK, "M": UnitM, "G": UnitG, "T": UnitT}[v]
	return nil
}

type SysinfoDomain struct {
	CPU       CpuDomain          `json:"cpu"`
	RAM       MemoryDomain       `json:"ram"`
	SWAP      MemoryDomain       `json:"swap"`
	Loadavg   LoadavgDomain      `json:"loadavg"`
	Temps     []TempDeviceDomain `json:"temps"`
	DiskUsage []DiskUsageDomain  `json:"diskusage"`
	Uptime    uint64             `json:"uptime"`
}

type CpuDomain struct {
	Cores []CpuCoreDomain `json:"cores"`
}

type CpuCoreDomain struct {
	Id    int     `json:"id"`
	Usage float32 `json:"usage"`
	MHz   float64 `json:"mhz"`
}

type MemoryDomain struct {
	Used  float64    `json:"used"`
	Total float64    `json:"total"`
	Unit  MemoryUnit `json:"unit"`
}

func (m *MemoryDomain) tidyValues() {
	for i := 1; i < 4; i++ {
		val := m.Total / 1024
		if val > 1 {
			m.Total /= 1024
			m.Used /= 1024
			m.Unit = MemoryUnit(byte(i))
		}
	}
}

type LoadavgDomain struct {
	Loadavg1  float32 `json:"loadavg1"`
	Loadavg5  float32 `json:"loadavg5"`
	Loadavg15 float32 `json:"loadavg15"`
}

type TempDeviceDomain struct {
	Name    string             `json:"name"`
	Sensors []TempSensorDomain `json:"sensors"`
}

type TempSensorDomain struct {
	Name string  `json:"name"`
	Temp float32 `json:"temp"`
}

type DiskUsageDomain struct {
	Path    string  `json:"path"`
	UsedGB  float64 `json:"usedgb"`
	TotalGB float64 `json:"totalgb"`
}

func HandleSysinfoData(w http.ResponseWriter, r *http.Request) {
	result := SysinfoDomain{}

	pCpu := metric.LoadCpu()
	time.Sleep(200 * time.Millisecond)
	cCpu := metric.LoadCpu()
	cpufreq := metric.LoadCpufreq()
	loadavg := metric.GetLoadavg()
	mem := metric.LoadMemInfo()
	temps := metric.LoadTemps()
	diskusage := metric.LoadDiskUsage()
	uptime := metric.LoadUptime()

	// Set CPU
	result.CPU = CpuDomain{
		Cores: make([]CpuCoreDomain, len(pCpu.Cores)),
	}
	for i := 0; i < len(cCpu.Cores); i++ {
		result.CPU.Cores[i] = CpuCoreDomain{
			Id:    cpufreq[i].Processor,
			Usage: cCpu.Cores[i].Usage(pCpu.Cores[i]),
			MHz:   cpufreq[i].Freq,
		}
	}

	// Set RAM & Swap
	result.RAM = MemoryDomain{Used: float64(mem.MemTotal - mem.MemAvailable), Total: float64(mem.MemTotal)}
	result.SWAP = MemoryDomain{Used: float64(mem.SwapTotal - mem.SwapFree), Total: float64(mem.SwapTotal)}

	result.RAM.tidyValues()
	result.SWAP.tidyValues()

	// Set Loadavgs
	result.Loadavg = LoadavgDomain{Loadavg1: loadavg.Loadavg1, Loadavg5: loadavg.Loadavg5, Loadavg15: loadavg.Loadavg15}

	// Set temps
	result.Temps = make([]TempDeviceDomain, 0)
	for _, t := range temps {
		device := TempDeviceDomain{Name: t.Name, Sensors: make([]TempSensorDomain, 0)}
		for _, s := range t.Temps {
			device.Sensors = append(device.Sensors, TempSensorDomain{Name: s.Label, Temp: float32(s.Input) / 1000})
		}
		result.Temps = append(result.Temps, device)
	}

	// Set diskusage
	result.DiskUsage = make([]DiskUsageDomain, 0)
	for _, d := range diskusage {
		result.DiskUsage = append(result.DiskUsage, DiskUsageDomain{Path: d.Path, UsedGB: float64(d.UsedSize) / 1024 / 1024 / 1024, TotalGB: float64(d.TotalSize) / 1024 / 1024 / 1024})
	}

	// Set uptime
	result.Uptime = uptime

	w.Header().Set("content-type", "application/json")
	e := json.NewEncoder(w)
	e.Encode(result)
}
