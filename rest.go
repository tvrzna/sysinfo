package main

import (
	"encoding/json"
	"net/http"

	"github.com/tvrzna/sysinfo/metric"
)

type MemoryUnit byte

const (
	UnitB MemoryUnit = iota
	UnitK
	UnitM
	UnitG
	UnitT
)

func (b MemoryUnit) String() string {
	return []string{"B", "K", "M", "G", "T"}[int(b)]
}

func (b MemoryUnit) MarshalJSON() ([]byte, error) {
	return []byte("\"" + b.String() + "\""), nil
}

func (b *MemoryUnit) UnmarshalJSON(data []byte) error {
	var v string
	json.Unmarshal(data, &v)
	*b = map[string]MemoryUnit{"B": UnitB, "K": UnitK, "M": UnitM, "G": UnitG, "T": UnitT}[v]
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
	Used      float64    `json:"used"`
	UsedUnit  MemoryUnit `json:"usedUnit"`
	Total     float64    `json:"total"`
	TotalUnit MemoryUnit `json:"totalUnit"`
	Percent   float64    `json:"percent"`
}

func (m *MemoryDomain) tidyValues() {
	if m.Used > 0 && m.Total > 0 {
		m.Percent = m.Used / m.Total * 100
	}
	m.Total, m.TotalUnit = tidyPrefix(m.Total, 1)
	m.Used, m.UsedUnit = tidyPrefix(m.Used, 1)
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
	Path      string     `json:"path"`
	Used      float64    `json:"used"`
	UsedUnit  MemoryUnit `json:"usedUnit"`
	Total     float64    `json:"total"`
	TotalUnit MemoryUnit `json:"totalUnit"`
	Percent   float64    `json:"percent"`
}

func (d *DiskUsageDomain) tidyValues() {
	if d.Used > 0 && d.Total > 0 {
		d.Percent = d.Used / d.Total * 100
	}
	d.Used, d.UsedUnit = tidyPrefix(d.Used, 0)
	d.Total, d.TotalUnit = tidyPrefix(d.Total, 0)
}

func tidyPrefix(value float64, start byte) (float64, MemoryUnit) {
	result := value
	resultUnit := MemoryUnit(start)
	for i := start; i < 5; i++ {
		if val := result / 1024; val < 1 {
			break
		} else {
			result = val
			resultUnit = MemoryUnit(i + 1)
		}
	}
	return result, resultUnit
}

func HandleSysinfoData(w http.ResponseWriter, r *http.Request) {
	result := SysinfoDomain{}

	bundle := &metric.Bundle{}
	doneCh := make(chan bool, 1)

	go metric.LoadCpu(doneCh, bundle)

	for i := 0; i < cap(doneCh); i++ {
		<-doneCh
	}

	bundle.Cpufreq = metric.LoadCpufreq()
	bundle.Loadavg = metric.GetLoadavg()
	bundle.Mem = metric.LoadMemInfo()
	bundle.Temps = metric.LoadTemps()
	bundle.Diskusage = metric.LoadDiskUsage()
	bundle.Uptime = metric.LoadUptime()

	// Set CPU
	result.CPU = CpuDomain{
		Cores: make([]CpuCoreDomain, len(bundle.Cpu.Cores)),
	}
	for i := 0; i < len(bundle.Cpu.Cores); i++ {
		result.CPU.Cores[i] = CpuCoreDomain{
			Id:    bundle.Cpufreq[i].Processor,
			Usage: bundle.Cpu.Cores[i].Usage,
			MHz:   bundle.Cpufreq[i].Freq,
		}
	}

	// Set RAM & Swap
	result.RAM = MemoryDomain{Used: float64(bundle.Mem.MemTotal - bundle.Mem.MemAvailable), Total: float64(bundle.Mem.MemTotal)}
	result.SWAP = MemoryDomain{Used: float64(bundle.Mem.SwapTotal - bundle.Mem.SwapFree), Total: float64(bundle.Mem.SwapTotal)}

	result.RAM.tidyValues()
	result.SWAP.tidyValues()

	// Set Loadavgs
	result.Loadavg = LoadavgDomain{Loadavg1: bundle.Loadavg.Loadavg1, Loadavg5: bundle.Loadavg.Loadavg5, Loadavg15: bundle.Loadavg.Loadavg15}

	// Set temps
	result.Temps = make([]TempDeviceDomain, 0)
	for _, t := range bundle.Temps {
		device := TempDeviceDomain{Name: t.Name, Sensors: make([]TempSensorDomain, 0)}
		for _, s := range t.Temps {
			device.Sensors = append(device.Sensors, TempSensorDomain{Name: s.Label, Temp: float32(s.Input) / 1000})
		}
		result.Temps = append(result.Temps, device)
	}

	// Set diskusage
	result.DiskUsage = make([]DiskUsageDomain, 0)
	for _, d := range bundle.Diskusage {
		diskUsage := DiskUsageDomain{Path: d.Path, Used: float64(d.UsedSize), Total: float64(d.TotalSize)}
		diskUsage.tidyValues()
		result.DiskUsage = append(result.DiskUsage, diskUsage)
	}

	// Set uptime
	result.Uptime = bundle.Uptime

	w.Header().Set("content-type", "application/json")
	e := json.NewEncoder(w)
	e.Encode(result)
}
