package metric

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	pathHwmon       = "/sys/class/hwmon"
	hwmonName       = "name"
	hwmonTempPrefix = "temp"
	hwmonTempInput  = "_input"
	hwmonTempLabel  = "_label"
)

type TempDevice struct {
	Name  string
	Temps []*TempInput
}

func (h *TempDevice) getTempById(id int) *TempInput {
	for _, t := range h.Temps {
		if t.Id == id {
			return t
		}
	}
	t := &TempInput{Id: id}
	h.Temps = append(h.Temps, t)
	return t
}

type TempInput struct {
	Id    int
	Label string
	Input int
}

func LoadTemps() []*TempDevice {
	var result []*TempDevice

	hws, err := os.ReadDir(pathHwmon)
	if err != nil {
		return nil
	}

	for _, hw := range hws {
		record := &TempDevice{Temps: make([]*TempInput, 0)}
		data, err := os.ReadDir(filepath.Join(pathHwmon, hw.Name()))
		if err != nil {
			continue
		}

		for _, f := range data {
			if f.Name() == hwmonName {
				bName, _ := os.ReadFile(filepath.Join(pathHwmon, hw.Name(), f.Name()))
				record.Name = strings.TrimSpace(string(bName))
			} else if strings.HasPrefix(f.Name(), hwmonTempPrefix) {
				idEnd := strings.Index(f.Name(), "_")
				id, _ := strconv.Atoi(f.Name()[4:idEnd])
				temp := record.getTempById(id)
				if strings.HasSuffix(f.Name(), hwmonTempLabel) {
					bLabel, _ := os.ReadFile(filepath.Join(pathHwmon, hw.Name(), f.Name()))
					temp.Label = strings.TrimSpace(string(bLabel))
				} else if strings.HasSuffix(f.Name(), hwmonTempInput) {
					bInput, _ := os.ReadFile(filepath.Join(pathHwmon, hw.Name(), f.Name()))
					input, _ := strconv.Atoi(strings.TrimSpace(string(bInput)))
					temp.Input = input
				}
			}
		}
		if len(record.Temps) > 0 {
			result = append(result, record)
		}
	}

	return result
}
