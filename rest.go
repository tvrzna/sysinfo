package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type SysinfoDomain struct {
	CPU     CpuDomain          `json:"cpu"`
	RAM     MemoryDomain       `json:"ram"`
	SWAP    MemoryDomain       `json:"swap"`
	Loadavg LoadavgDomain      `json:"loadavg"`
	Temps   []TempDeviceDomain `json:"temps"`
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

func HandleSysinfoData(w http.ResponseWriter, r *http.Request) {
	result := SysinfoDomain{}

	pCpu := LoadCpu()
	time.Sleep(200 * time.Millisecond)
	cCpu := LoadCpu()
	cpufreq := LoadCpufreq()
	loadavg := GetLoadavg()
	mem := LoadMemInfo()
	temps := LoadTemps()

	// Set CPU
	result.CPU = CpuDomain{
		Cores: make([]CpuCoreDomain, len(pCpu.cores)),
	}
	for i := 0; i < len(result.CPU.Cores); i++ {
		result.CPU.Cores[i] = CpuCoreDomain{
			Id:    cpufreq[i].Processor,
			Usage: cCpu.cores[i].Usage(pCpu.cores[i]),
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
		device := TempDeviceDomain{Name: t.name, Sensors: make([]TempSensorDomain, 0)}
		for _, s := range t.temps {
			device.Sensors = append(device.Sensors, TempSensorDomain{Name: s.label, Temp: float32(s.input) / 1000})
		}
		result.Temps = append(result.Temps, device)
	}

	w.Header().Set("content-type", "application/json")
	e := json.NewEncoder(w)
	e.Encode(result)
}
