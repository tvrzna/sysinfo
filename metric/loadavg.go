package metric

import (
	"fmt"
	"os"
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
	result := &Loadavg{}
	b, _ := os.ReadFile(pathLoadavg)
	fmt.Sscanf(string(b), "%f %f %f", &result.Loadavg1, &result.Loadavg5, &result.Loadavg15)
	return result
}
