package metric

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	pathNet     = "/sys/class/net/"
	pathRxBytes = "statistics/rx_bytes"
	pathTxBytes = "statistics/tx_bytes"
)

type Netspeed struct {
	Name     string
	Download float64
	Upload   float64
}

func LoadNetspeed(doneCh chan bool, bundle *Bundle) {
	previous := loadNetspeed()
	time.Sleep(200 * time.Millisecond)
	bundle.Netspeed = loadNetspeed()
	for _, n := range bundle.Netspeed {
		var o *Netspeed
		for _, on := range previous {
			if on.Name == n.Name {
				o = on
				break
			}
		}
		if o == nil {
			continue
		}
		n.Download = (n.Download - o.Download) * 5
		n.Upload = (n.Upload - o.Upload) * 5
	}
	doneCh <- true
}

func loadNetspeed() []*Netspeed {
	var result []*Netspeed
	ifaces, _ := os.ReadDir(pathNet)
	for _, iface := range ifaces {
		if iface.Name() == "lo" {
			continue
		}
		strDownload, _ := os.ReadFile(filepath.Join(pathNet, iface.Name(), pathRxBytes))
		strUpload, _ := os.ReadFile(filepath.Join(pathNet, iface.Name(), pathTxBytes))

		download, _ := strconv.ParseFloat(strings.TrimSpace(string(strDownload)), 64)
		upload, _ := strconv.ParseFloat(strings.TrimSpace(string(strUpload)), 64)

		result = append(result, &Netspeed{iface.Name(), download, upload})
	}
	return result
}
