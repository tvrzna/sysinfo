package main

import (
	"encoding/json"
	"math"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/tvrzna/sysinfo/metric"
)

type restContext struct {
	m        *sync.Mutex
	lastLoad int64
	sysinfo  *SysinfoDomain
}

func initRestContext() *restContext {
	return &restContext{m: &sync.Mutex{}}
}

func (c *restContext) HandleSysinfoData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	e := json.NewEncoder(w)

	c.m.Lock()
	defer c.m.Unlock()
	if c.lastLoad+450 < time.Now().UnixMilli() {
		c.sysinfo = c.loadSysinfo()
		c.lastLoad = time.Now().UnixMilli()
	}

	e.Encode(c.sysinfo)
}

func (c *restContext) loadSysinfo() *SysinfoDomain {
	result := &SysinfoDomain{}

	bundle := &metric.Bundle{}
	doneCh := make(chan bool, 3)

	go metric.LoadCpu(doneCh, bundle)
	go metric.LoadNetspeed(doneCh, bundle)
	go metric.LoadTop(doneCh, bundle)

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
	result.CPU = &CpuDomain{
		Cores: make([]*CpuCoreDomain, len(bundle.Cpu.Cores)),
	}
	for i := 0; i < len(bundle.Cpu.Cores); i++ {
		result.CPU.Cores[i] = &CpuCoreDomain{
			Id:    bundle.Cpufreq[i].Processor,
			Usage: bundle.Cpu.Cores[i].Usage,
			MHz:   bundle.Cpufreq[i].Freq,
		}
	}

	// Set RAM & Swap
	result.RAM = &MemoryDomain{Used: float64(bundle.Mem.MemTotal - bundle.Mem.MemAvailable), Total: float64(bundle.Mem.MemTotal)}
	result.SWAP = &MemoryDomain{Used: float64(bundle.Mem.SwapTotal - bundle.Mem.SwapFree), Total: float64(bundle.Mem.SwapTotal)}

	result.RAM.tidyValues()
	result.SWAP.tidyValues()

	// Set Loadavgs
	result.Loadavg = &LoadavgDomain{Loadavg1: bundle.Loadavg.Loadavg1, Loadavg5: bundle.Loadavg.Loadavg5, Loadavg15: bundle.Loadavg.Loadavg15}

	// Set temps
	result.Temps = make([]*TempDeviceDomain, 0)
	for _, t := range bundle.Temps {
		device := &TempDeviceDomain{Name: t.Name, Sensors: make([]TempSensorDomain, 0)}
		for _, s := range t.Temps {
			device.Sensors = append(device.Sensors, TempSensorDomain{Name: s.Label, Temp: float32(s.Input) / 1000})
		}
		result.Temps = append(result.Temps, device)
	}

	// Set diskusage
	result.DiskUsage = make([]*DiskUsageDomain, 0)
	for _, d := range bundle.Diskusage {
		diskUsage := &DiskUsageDomain{Path: d.Path, Used: float64(d.UsedSize), Total: float64(d.TotalSize)}
		diskUsage.tidyValues()
		result.DiskUsage = append(result.DiskUsage, diskUsage)
	}

	// Set uptime
	result.Uptime = bundle.Uptime

	// Set netspeed
	result.Netspeed = make([]*NetspeedDomain, 0)
	for _, n := range bundle.Netspeed {
		netspeed := &NetspeedDomain{Name: n.Name, Download: n.Download, Upload: n.Upload}
		netspeed.tidyValues()
		result.Netspeed = append(result.Netspeed, netspeed)
	}

	// Set Proc
	result.Top = make([]*ProcDomain, 0)
	for _, p := range bundle.Top {
		pid := &ProcDomain{PID: p.PID, Comm: p.Comm[1 : len(p.Comm)-1], State: string(p.State), Cpu: p.CpuUsage, RamUsage: float64(p.RamUsage)}
		pid.tidyValues()
		result.Top = append(result.Top, pid)
	}
	sort.Slice(result.Top, func(i, j int) bool {
		if result.Top[i].Cpu == result.Top[j].Cpu {
			return (result.Top[i].RamUsage * (math.Pow(1024, float64(result.Top[i].RamUnit)))) > (result.Top[j].RamUsage * (math.Pow(1024, float64(result.Top[j].RamUnit))))
		}
		return result.Top[i].Cpu > result.Top[j].Cpu
	})
	if len(result.Top) > 20 {
		result.Top = result.Top[:20]
	}

	return result
}
