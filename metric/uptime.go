package metric

import (
	"os"
	"strconv"
	"strings"
)

const (
	pathUptime = "/proc/uptime"
)

func LoadUptime() uint64 {
	b, _ := os.ReadFile(pathUptime)
	data := strings.Split(string(b), " ")
	value, _ := strconv.ParseFloat(data[0], 64)
	return uint64(value)
}
