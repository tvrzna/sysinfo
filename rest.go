package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tvrzna/sysinfo/metric"
)

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
	Used  float32 `json:"used"`
	Total float32 `json:"total"`
	Usage float32 `json:"usage"`
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
	for i := 0; i < len(result.CPU.Cores); i++ {
		result.CPU.Cores[i] = CpuCoreDomain{
			Id:    cpufreq[i].Processor,
			Usage: cCpu.Cores[i].Usage(pCpu.Cores[i]),
			MHz:   cpufreq[i].Freq,
		}
	}

	// Set RAM & Swap
	result.RAM = MemoryDomain{Used: mem.MemTotalGB() - mem.MemAvailableGB(), Total: mem.MemTotalGB()}
	result.RAM.Usage = result.RAM.Used / result.RAM.Total * 100
	result.SWAP = MemoryDomain{Used: mem.SwapTotalGB() - mem.SwapFreeGB(), Total: mem.SwapTotalGB()}
	result.SWAP.Usage = result.SWAP.Used / result.SWAP.Total * 100

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
