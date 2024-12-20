package metric

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	pathProc    = "/proc/"
	pathTopStat = "stat"
)

type TopProcess struct {
	PID      int
	Comm     string
	State    byte
	PPID     int
	Utime    uint64
	Stime    uint64
	Cutime   int64
	Cstime   int64
	RSS      int64
	Now      uint64
	CpuUsage float32
	RamUsage uint64
}

var pageSize = uint64(os.Getpagesize())

func (t *TopProcess) calc(p *TopProcess) {
	if p != nil {
		t.CpuUsage = float32(1 / (float32(t.Now-p.Now) / 1000) * float32(t.getTotal()-p.getTotal()))
	}
	t.RamUsage = uint64(t.RSS) * pageSize
}

func (t *TopProcess) getTotal() uint64 {
	return t.Utime + t.Stime + uint64(t.Cutime+t.Cstime)
}

func LoadTop(wg *sync.WaitGroup, bundle *Bundle) {
	previous := loadTop()
	time.Sleep(400 * time.Millisecond)
	bundle.Top = loadTop()

	for _, p := range bundle.Top {
		p.calc(previous[p.PID])
	}
	wg.Done()
}

func loadTop() map[int]*TopProcess {
	result := make(map[int]*TopProcess)
	proc, _ := readDir(pathProc)
	for _, pid := range proc {
		if pid.IsDir() {
			if _, err := strconv.Atoi(pid.Name()); err != nil {
				continue
			}
			p := &TopProcess{}

			b, err := os.ReadFile(filepath.Join(pathProc, pid.Name(), pathTopStat))
			if err != nil {
				continue
			}
			data := strings.SplitN(string(b), " ", 26)

			p.PID, _ = strconv.Atoi(data[0])
			p.Comm = data[1]
			p.State = data[2][0]
			p.PPID, _ = strconv.Atoi(data[3])
			p.Utime, _ = strconv.ParseUint(data[13], 10, 64)
			p.Stime, _ = strconv.ParseUint(data[14], 10, 64)
			p.Cutime, _ = strconv.ParseInt(data[15], 10, 64)
			p.Cstime, _ = strconv.ParseInt(data[16], 10, 64)
			p.RSS, _ = strconv.ParseInt(data[23], 10, 64)
			p.Now = uint64(time.Now().UnixMilli())

			result[p.PID] = p
		}
	}
	return result
}

func readDir(path string) ([]os.DirEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.ReadDir(-1)
}
