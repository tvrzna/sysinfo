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
	Prefailure    bool `json:"prefailure"`
	UpdatedOnline bool `json:"updated_online"`
	Performance   bool `json:"performance"`
	ErrorRate     bool `json:"error_rate"`
	EventCount    bool `json:"event_count"`
	AutoKeep      bool `json:"auto_keep"`
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

type SmartctlOutput struct {
	Smartctl           Smartctl           `json:"smartctl"`
	Devices            []Device           `json:"devices"`
	Device             Device             `json:"device"`
	SmartStatus        SmartStatus        `json:"smart_status"`
	AtaSmartAttributes AtaSmartAttributes `json:"ata_smart_attributes"`
	SmartctlExitStatus int                `json:"smartctl_exit_status"`
}

func (s *SmartctlOutput) parse(data []byte) error {
	return json.Unmarshal(data, s)
}
