package metric

type Bundle struct {
	Cpu       *Cpu
	Cpufreq   []*Cpufreq
	Loadavg   *Loadavg
	Mem       *MemoryInfo
	Temps     []*TempDevice
	Diskusage []*DiskUsage
	Uptime    uint64
	Netspeed  []*Netspeed
	Top       map[int]*TopProcess
	Diskstats []*Diskstat
}
