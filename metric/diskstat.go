package metric

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	pathBlocks    = "/sys/class/block/"
	pathBlockStat = "stat"
	pathBlockSize = "queue/physical_block_size"
)

type Diskstat struct {
	Name       string
	SectorSize int
	Riops      uint64
	Read       float64
	Wiops      uint64
	Write      float64
}

func LoadDiskstats(doneCh chan bool, bundle *Bundle) {
	previous := loadDiskstats()
	time.Sleep(500 * time.Millisecond)
	bundle.Diskstats = loadDiskstats()

	for _, b := range bundle.Diskstats {
		var p *Diskstat
		for _, i := range previous {
			if b.Name == i.Name {
				p = i
				break
			}
		}
		if p == nil {
			continue
		}
		b.Riops = b.Riops - p.Riops
		b.Read = (b.Read - p.Read) * 2 * float64(512)
		b.Wiops = b.Wiops - p.Wiops
		b.Write = (b.Write - p.Write) * 2 * float64(512)
	}
	doneCh <- true
}

func loadDiskstats() []*Diskstat {
	result := make([]*Diskstat, 0)

	blocks, _ := readDir(pathBlocks)
	for _, b := range blocks {
		block := &Diskstat{Name: b.Name()}
		d, _ := os.ReadFile(filepath.Join(pathBlocks, b.Name(), pathBlockSize))
		block.SectorSize, _ = strconv.Atoi(strings.TrimSpace(string(d)))
		if block.SectorSize == 0 {
			continue
		}

		s, _ := os.ReadFile(filepath.Join(pathBlocks, b.Name(), pathBlockStat))
		data := strings.Fields(string(s))
		block.Riops, _ = strconv.ParseUint(data[0], 10, 64)
		block.Read, _ = strconv.ParseFloat(data[2], 64)
		block.Wiops, _ = strconv.ParseUint(data[4], 10, 64)
		block.Write, _ = strconv.ParseFloat(data[6], 64)
		result = append(result, block)
	}

	return result
}
