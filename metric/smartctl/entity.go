package smartctl

import (
	"encoding/json"
)

type Smartctl struct {
	Version    []int `json:"version"`
	ExitStatus int   `json:"exit_status"`
}

type Device struct {
	Name     string `json:"name"`
	InfoName string `json:"info_name"`
	Type     string `json:"type"`
	Protocol string `json:"protocol"`
}

type SmartStatus struct {
	Passed bool `json:"passed"`
}

type Flags struct {
	String        string `json:"string"`
	Prefailure    bool   `json:"prefailure"`
	UpdatedOnline bool   `json:"updated_online"`
	Performance   bool   `json:"performance"`
	ErrorRate     bool   `json:"error_rate"`
	EventCount    bool   `json:"event_count"`
	AutoKeep      bool   `json:"auto_keep"`
}

type Raw struct {
	Value  int    `json:"value"`
	String string `json:"string"`
}

type Attribute struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Value      int    `json:"value"`
	Worst      int    `json:"worst"`
	Thresh     int    `json:"thresh"`
	WhenFailed string `json:"when_failed"`
	Flags      Flags  `json:"flags"`
	Raw        Raw    `json:"raw"`
}

type AtaSmartAttributes struct {
	Revision int         `json:"revision"`
	Table    []Attribute `json:"table"`
}

type NVMeSmartHealthInformationLog struct {
	CriticalWarning         int   `json:"critical_warning"`
	Temperature             int   `json:"temperature"`
	AvailableSpare          int   `json:"available_spare"`
	AvailableSpareThreshold int   `json:"available_spare_threshold"`
	PercentageUsed          int   `json:"percentage_used"`
	DataUnitsRead           int   `json:"data_units_read"`
	DataUnitsWritten        int   `json:"data_units_written"`
	HostReadCommands        int   `json:"host_read_commands"`
	HostWriteCommands       int   `json:"host_write_commands"`
	ControllerBusyTime      int   `json:"controller_busy_time"`
	PowerCycles             int   `json:"power_cycles"`
	PowerOnHours            int   `json:"power_on_hours"`
	UnsafeShutdowns         int   `json:"unsafe_shutdowns"`
	MediaErrors             int   `json:"media_errors"`
	NumErrLogEntries        int   `json:"num_err_log_entries"`
	WarningTempTime         int   `json:"warning_temp_time"`
	CriticalCompTime        int   `json:"critical_comp_time"`
	TemperatureSensors      []int `json:"temperature_sensors"`
}

type PowerOnTime struct {
	Hours int `json:"hours"`
}

type Temperature struct {
	Current int `json:"current"`
}

type SmartctlOutput struct {
	Smartctl                      Smartctl                       `json:"smartctl"`
	Devices                       []Device                       `json:"devices"`
	Device                        Device                         `json:"device"`
	ModelName                     string                         `json:"model_name"`
	SmartStatus                   SmartStatus                    `json:"smart_status"`
	NVMeSmartHealthInformationLog *NVMeSmartHealthInformationLog `json:"nvme_smart_health_information_log"`
	AtaSmartAttributes            AtaSmartAttributes             `json:"ata_smart_attributes"`
	SmartctlExitStatus            int                            `json:"smartctl_exit_status"`
	PowerOnTime                   PowerOnTime                    `json:"power_on_time"`
	PowerCycleCount               int                            `json:"power_cycle_count"`
	Temperature                   Temperature                    `json:"temperature"`
}

func (s *SmartctlOutput) parse(data []byte) error {
	return json.Unmarshal(data, s)
}
