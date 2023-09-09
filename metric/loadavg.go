package metric

import (
	"os"
	"strconv"
	"strings"
)

const (
	pathLoadavg = "/proc/loadavg"
)

type Loadavg struct {
	Loadavg1  float32
	Loadavg5  float32
	Loadavg15 float32
}

func GetLoadavg() *Loadavg {
	b, _ := os.ReadFile(pathLoadavg)
	data := strings.Split(string(b), " ")
	return &Loadavg{
		Loadavg1:  stringToFloat32(data[0]),
		Loadavg5:  stringToFloat32(data[1]),
		Loadavg15: stringToFloat32(data[2]),
	}
}

func stringToFloat32(val string) float32 {
	value, err := strconv.ParseFloat(val, 32)
	if err == nil {
		return float32(value)
	}
	return 0
}
