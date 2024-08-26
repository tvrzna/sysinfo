package main

import "encoding/json"

type MemoryUnit byte

const (
	UnitB MemoryUnit = iota
	UnitK
	UnitM
	UnitG
	UnitT
)

func (b MemoryUnit) String() string {
	return []string{"B", "K", "M", "G", "T", "P"}[int(b)]
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
	CPU       *CpuDomain          `json:"cpu,omitempty"`
	RAM       *MemoryDomain       `json:"ram,omitempty"`
	SWAP      *MemoryDomain       `json:"swap,omitempty"`
	Temps     []*TempDeviceDomain `json:"temps,omitempty"`
	DiskUsage []*DiskUsageDomain  `json:"diskusage,omitempty"`
	Netspeed  []*NetspeedDomain   `json:"netspeed,omitempty"`
	Top       []*ProcDomain       `json:"top,omitempty"`
	Diskstats []*DiskstatDomain   `json:"diskstats,omitempty"`
	Smartctl  []*SmartctlDomain   `json:"smartctl,omitempty"`
	System    *SystemDomain       `json:"system,omitempty"`
}

type CpuDomain struct {
	Cores []*CpuCoreDomain `json:"cores"`
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

type NetspeedDomain struct {
	Name            string     `json:"name"`
	Download        float64    `json:"download"`
	DownloadUnit    MemoryUnit `json:"downloadUnit"`
	DownloadPercent float64    `json:"downloadPercent"`
	Upload          float64    `json:"upload"`
	UploadUnit      MemoryUnit `json:"uploadUnit"`
	UploadPercent   float64    `json:"uploadPercent"`
}

type ProcDomain struct {
	PID      int        `json:"pid"`
	Comm     string     `json:"comm"`
	State    string     `json:"state"`
	Cpu      float32    `json:"cpu"`
	RamUsage float64    `json:"ram"`
	RamUnit  MemoryUnit `json:"ramUnit"`
}

type DiskstatDomain struct {
	Name         string     `json:"name"`
	Riops        uint64     `json:"riops"`
	Read         float64    `json:"read"`
	ReadUnit     MemoryUnit `json:"readUnit"`
	ReadPercent  float64    `json:"readPercent"`
	Wiops        uint64     `json:"wiops"`
	Write        float64    `json:"write"`
	WriteUnit    MemoryUnit `json:"writeUnit"`
	WritePercent float64    `json:"writePercent"`
}

type SmartctlDomain struct {
	Name              string                    `json:"name"`
	Model             string                    `json:"model"`
	SmartStatusPassed bool                      `json:"smartStatusPassed"`
	PowerOnTime       int                       `json:"powerOnTime"`
	PowerCycleCount   int                       `json:"powerCycleCount"`
	Temperature       int                       `json:"temp"`
	Attributes        []SmartctlAttributeDomain `json:"attributes"`
	Nvme              *SmartctlNVME             `json:"nvme,omitempty"`
}

type SmartctlAttributeDomain struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
	Worst int    `json:"worst"`
	Raw   int    `json:"raw"`
	Flags string `json:"flags"`
}

type SmartctlNVME struct {
	CriticalWarning         int   `json:"criticalWarning"`
	Temperature             int   `json:"temperature"`
	AvailableSpare          int   `json:"availableSpare"`
	AvailableSpareThreshold int   `json:"availableSpareThreshold"`
	PercentageUsed          int   `json:"percentageUsed"`
	DataUnitsRead           int   `json:"dataUnitsRead"`
	DataUnitsWritten        int   `json:"dataUnitsWritten"`
	HostReadCommands        int   `json:"hostReadCommands"`
	HostWriteCommands       int   `json:"hostWriteCommands"`
	ControllerBusyTime      int   `json:"controllerBusyTime"`
	PowerCycles             int   `json:"powerCycles"`
	PowerOnHours            int   `json:"powerOnHours"`
	UnsafeShutdowns         int   `json:"unsafeShutdowns"`
	MediaErrors             int   `json:"mediaErrors"`
	NumErrLogEntries        int   `json:"numErrLogEntries"`
	WarningTempTime         int   `json:"warningTempTime"`
	CriticalCompTime        int   `json:"criticalCompTime"`
	TemperatureSensors      []int `json:"temperatureSensors"`
}

type SystemDomain struct {
	Loadavg   *LoadavgDomain `json:"loadavg,omitempty"`
	Uptime    uint64         `json:"uptime,omitempty"`
	Updates   int            `json:"updates"`
	Hostname  string         `json:"hostname"`
	OsType    string         `json:"ostype"`
	OsRelease string         `json:"osrelease"`
}

func (n *NetspeedDomain) tidyValues() {
	n.DownloadPercent = n.Download / 104857600 * 50
	n.UploadPercent = n.Upload / 104857600 * 50

	n.Download, n.DownloadUnit = tidyPrefix(n.Download, 0)
	n.Upload, n.UploadUnit = tidyPrefix(n.Upload, 0)
}

func (d *DiskUsageDomain) tidyValues() {
	if d.Used > 0 && d.Total > 0 {
		d.Percent = d.Used / d.Total * 100
	}
	d.Used, d.UsedUnit = tidyPrefix(d.Used, 0)
	d.Total, d.TotalUnit = tidyPrefix(d.Total, 0)
}

func (p *ProcDomain) tidyValues() {
	p.RamUsage, p.RamUnit = tidyPrefix(p.RamUsage, 0)
}

func (d *DiskstatDomain) tidyValues() {
	if d.Read > 0 {
		d.ReadPercent = d.Read / 104857600 * 100
	}
	if d.Write > 0 {
		d.WritePercent = d.Write / 104857600 * 100
	}
	d.Read, d.ReadUnit = tidyPrefix(d.Read, 0)
	d.Write, d.WriteUnit = tidyPrefix(d.Write, 0)
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
