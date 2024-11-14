package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/tvrzna/sysinfo/metric"
	"github.com/tvrzna/sysinfo/metric/smartctl"
)

type restContext struct {
	conf     *config
	m        *sync.Mutex
	s        *smartctl.SmartctlContext
	lastLoad int64
	sysinfo  *SysinfoDomain
}

func initRestContext(conf *config) *restContext {
	var s *smartctl.SmartctlContext
	if conf.widgetsIndex["smartctl"] {
		s = smartctl.CreateSmartctlContext()
	}

	return &restContext{m: &sync.Mutex{}, conf: conf, s: s}
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

	// Load metrics
	bundle := c.loadMetrics()

	// Set CPU
	if c.conf.widgetsIndex["cpu"] {
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
	}

	// Set RAM & Swap
	if c.conf.widgetsIndex["memory"] {
		result.RAM = &MemoryDomain{Used: float64(bundle.Mem.MemTotal - bundle.Mem.MemAvailable), Total: float64(bundle.Mem.MemTotal)}
		result.SWAP = &MemoryDomain{Used: float64(bundle.Mem.SwapTotal - bundle.Mem.SwapFree), Total: float64(bundle.Mem.SwapTotal)}

		result.RAM.tidyValues()
		result.SWAP.tidyValues()
	}

	// Set system
	if c.conf.widgetsIndex["system"] {
		result.System = &SystemDomain{
			&LoadavgDomain{Loadavg1: bundle.System.Loadavg.Loadavg1, Loadavg5: bundle.System.Loadavg.Loadavg5, Loadavg15: bundle.System.Loadavg.Loadavg15},
			bundle.System.Uptime,
			bundle.System.Updates,
			bundle.System.Hostname,
			bundle.System.OsType,
			bundle.System.OsRelease,
		}
	}

	// Set temps
	if c.conf.widgetsIndex["temps"] {
		result.Temps = make([]*TempDeviceDomain, 0)
		for _, t := range bundle.Temps {
			device := &TempDeviceDomain{Name: t.Name, Sensors: make([]TempSensorDomain, 0)}
			for _, s := range t.Temps {
				device.Sensors = append(device.Sensors, TempSensorDomain{Name: s.Label, Temp: float32(s.Input) / 1000})
			}
			result.Temps = append(result.Temps, device)
		}
	}

	// Set diskusage
	if c.conf.widgetsIndex["diskusage"] {
		result.DiskUsage = make([]*DiskUsageDomain, 0)
		for _, d := range bundle.Diskusage {
			diskUsage := &DiskUsageDomain{Path: d.Path, Used: float64(d.UsedSize), Total: float64(d.TotalSize)}
			diskUsage.tidyValues()
			result.DiskUsage = append(result.DiskUsage, diskUsage)
		}
	}

	// Set netspeed
	if c.conf.widgetsIndex["netspeed"] {
		result.Netspeed = make([]*NetspeedDomain, 0)
		for _, n := range bundle.Netspeed {
			netspeed := &NetspeedDomain{Name: n.Name, Download: n.Download, Upload: n.Upload}
			netspeed.tidyValues()
			result.Netspeed = append(result.Netspeed, netspeed)
		}
	}

	// Set top
	if c.conf.widgetsIndex["top"] {
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
	}

	// Set diskstats
	if c.conf.widgetsIndex["diskstats"] {
		result.Diskstats = make([]*DiskstatDomain, 0)
		for _, d := range bundle.Diskstats {
			diskstat := &DiskstatDomain{Name: d.Name, Riops: d.Riops, Read: d.Read, Wiops: d.Wiops, Write: d.Write}
			diskstat.tidyValues()
			result.Diskstats = append(result.Diskstats, diskstat)
		}
		sort.Slice(result.Diskstats, func(i, j int) bool {
			return strings.ToLower(result.Diskstats[i].Name) < strings.ToLower(result.Diskstats[j].Name)
		})
	}

	// Set smartctl data
	if c.conf.widgetsIndex["smartctl"] {
		result.Smartctl = make([]*SmartctlDomain, 0)
		for _, s := range bundle.Smartctl {
			device := &SmartctlDomain{Name: s.Device.Name, Model: s.ModelName, SmartStatusPassed: s.SmartStatus.Passed, PowerOnTime: s.PowerOnTime.Hours, PowerCycleCount: s.PowerCycleCount, Temperature: s.Temperature.Current, Attributes: make([]SmartctlAttributeDomain, 0)}
			for _, a := range s.AtaSmartAttributes.Table {
				attr := SmartctlAttributeDomain{Name: a.Name, Value: a.Value, Worst: a.Worst, Raw: a.Raw.Value, Flags: a.Flags.String}
				device.Attributes = append(device.Attributes, attr)
			}
			if s.NVMeSmartHealthInformationLog != nil {
				n := s.NVMeSmartHealthInformationLog
				device.Nvme = &SmartctlNVME{n.CriticalWarning, n.Temperature, n.AvailableSpare, n.AvailableSpareThreshold, n.PercentageUsed, n.DataUnitsRead, n.DataUnitsWritten, n.HostReadCommands, n.HostWriteCommands, n.ControllerBusyTime, n.PowerCycles, n.PowerOnHours, n.UnsafeShutdowns, n.MediaErrors, n.NumErrLogEntries, n.WarningTempTime, n.CriticalCompTime, n.TemperatureSensors}
			}
			result.Smartctl = append(result.Smartctl, device)
		}
	}

	return result
}

func (c *restContext) loadMetrics() *metric.Bundle {
	parallelMetrics := make([]func(*sync.WaitGroup, *metric.Bundle), 0)
	if c.conf.widgetsIndex["cpu"] {
		parallelMetrics = append(parallelMetrics, metric.LoadCpu)
	}
	if c.conf.widgetsIndex["netspeed"] {
		parallelMetrics = append(parallelMetrics, metric.LoadNetspeed)
	}
	if c.conf.widgetsIndex["top"] {
		parallelMetrics = append(parallelMetrics, metric.LoadTop)
	}
	if c.conf.widgetsIndex["diskstats"] {
		parallelMetrics = append(parallelMetrics, metric.LoadDiskstats)
	}

	bundle := &metric.Bundle{}
	wg := &sync.WaitGroup{}
	for _, f := range parallelMetrics {
		wg.Add(1)
		go f(wg, bundle)
	}
	wg.Wait()

	if c.conf.widgetsIndex["cpu"] {
		bundle.Cpufreq = metric.LoadCpufreq()
	}
	if c.conf.widgetsIndex["system"] {
		bundle.System = &metric.System{
			Loadavg:   metric.GetLoadavg(),
			Uptime:    metric.LoadUptime(),
			Updates:   metric.LoadPkgUpdates(),
			Hostname:  metric.LoadHostname(),
			OsType:    metric.LoadOsType(),
			OsRelease: metric.LoadOsRelease(),
		}
	}
	if c.conf.widgetsIndex["memory"] {
		bundle.Mem = metric.LoadMemInfo()
	}
	if c.conf.widgetsIndex["temps"] {
		bundle.Temps = metric.LoadTemps()
	}
	if c.conf.widgetsIndex["diskusage"] {
		bundle.Diskusage = metric.LoadDiskUsage()
	}
	if c.conf.widgetsIndex["smartctl"] {
		var err error
		bundle.Smartctl, err = c.s.GetReports()
		if err != nil {
			log.Print(err)
		}
	}

	return bundle
}
