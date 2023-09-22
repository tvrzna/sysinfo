package metric

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	pathDiskstats = "/proc/diskstats"
	pathBlocks    = "/sys/class/block/"
	pathBlockStat = "stat"
	pathBlockSize = "queue/physical_block_size"
)

type Diskstat struct {
	Name       string
	Major      int
	Minor      int
	SectorSize int
	Riops      uint64
	Read       float64
	Wiops      uint64
	Write      float64
}

func LoadDiskstats(wg *sync.WaitGroup, bundle *Bundle) {
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
	wg.Done()
}

func loadDiskstats() []*Diskstat {
	result := make([]*Diskstat, 0)

	f, err := os.Open(pathDiskstats)
	if err != nil {
		return nil
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		data := strings.Fields(scanner.Text())

		block := &Diskstat{Name: data[2]}
		block.Major, _ = strconv.Atoi(data[0])
		block.Minor, _ = strconv.Atoi(data[1])
		if block.Minor > 0 {
			continue
		}

		d, _ := os.ReadFile(filepath.Join(pathBlocks, block.Name, pathBlockSize))
		block.SectorSize, _ = strconv.Atoi(strings.TrimSpace(string(d)))
		if block.SectorSize == 0 {
			continue
		}

		block.Riops, _ = strconv.ParseUint(data[3], 10, 64)
		block.Read, _ = strconv.ParseFloat(data[5], 64)
		block.Wiops, _ = strconv.ParseUint(data[7], 10, 64)
		block.Write, _ = strconv.ParseFloat(data[9], 64)
		result = append(result, block)
	}

	return result
}
