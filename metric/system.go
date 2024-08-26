package metric

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	pathLoadavg   = "/proc/loadavg"
	pathUptime    = "/proc/uptime"
	pathHostname  = "/proc/sys/kernel/hostname"
	pathOsType    = "/proc/sys/kernel/ostype"
	pathOsRelease = "/proc/sys/kernel/osrelease"
)

type System struct {
	Loadavg   *Loadavg
	Uptime    uint64
	Updates   int
	Hostname  string
	OsType    string
	OsRelease string
}

type Loadavg struct {
	Loadavg1  float32
	Loadavg5  float32
	Loadavg15 float32
}

func GetLoadavg() *Loadavg {
	result := &Loadavg{}
	b, _ := os.ReadFile(pathLoadavg)
	fmt.Sscanf(string(b), "%f %f %f", &result.Loadavg1, &result.Loadavg5, &result.Loadavg15)
	return result
}

func LoadUptime() uint64 {
	b, _ := os.ReadFile(pathUptime)
	data := strings.Split(string(b), " ")
	value, _ := strconv.ParseFloat(data[0], 64)
	return uint64(value)
}

func LoadHostname() string {
	b, _ := os.ReadFile(pathHostname)
	return strings.TrimSpace(string(b))
}

func LoadOsType() string {
	b, _ := os.ReadFile(pathOsType)
	return strings.TrimSpace(string(b))
}

func LoadOsRelease() string {
	b, _ := os.ReadFile(pathOsRelease)
	return strings.TrimSpace(string(b))
}
