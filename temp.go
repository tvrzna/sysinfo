package main

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
	name  string
	temps []*TempInput
}

func (h *TempDevice) getTempById(id int) *TempInput {
	for _, t := range h.temps {
		if t.id == id {
			return t
		}
	}
	t := &TempInput{id: id}
	h.temps = append(h.temps, t)
	return t
}

type TempInput struct {
	id    int
	label string
	input int
}

func LoadTemps() []*TempDevice {
	var result []*TempDevice

	hws, err := os.ReadDir(pathHwmon)
	if err != nil {
		return nil
	}

	for _, hw := range hws {
		record := &TempDevice{temps: make([]*TempInput, 0)}
		data, err := os.ReadDir(filepath.Join(pathHwmon, hw.Name()))
		if err != nil {
			continue
		}

		for _, f := range data {
			if f.Name() == hwmonName {
				bName, _ := os.ReadFile(filepath.Join(pathHwmon, hw.Name(), f.Name()))
				record.name = strings.TrimSpace(string(bName))
			} else if strings.HasPrefix(f.Name(), hwmonTempPrefix) {
				idEnd := strings.Index(f.Name(), "_")
				id, _ := strconv.Atoi(f.Name()[4:idEnd])
				temp := record.getTempById(id)
				if strings.HasSuffix(f.Name(), hwmonTempLabel) {
					bLabel, _ := os.ReadFile(filepath.Join(pathHwmon, hw.Name(), f.Name()))
					temp.label = strings.TrimSpace(string(bLabel))
				} else if strings.HasSuffix(f.Name(), hwmonTempInput) {
					bInput, _ := os.ReadFile(filepath.Join(pathHwmon, hw.Name(), f.Name()))
					input, _ := strconv.Atoi(strings.TrimSpace(string(bInput)))
					temp.input = input
				}
			}
		}
		if len(record.temps) > 0 {
			result = append(result, record)
		}
	}

	return result
}
