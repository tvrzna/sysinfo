package smartctl

import (
	"log"
	"time"
)

const (
	scanMinFrequency    int64 = 1800
	getReportsFrequency int64 = 3600
)

type SmartctlContext struct {
	devices     []string
	reports     []*SmartctlOutput
	lastScan    int64
	lastReports int64
}

func CreateSmartctlContext() *SmartctlContext {
	s := &SmartctlContext{}
	go s.runMondayTasker()
	return s
}

func (s *SmartctlContext) isAfter(last, offset int64) bool {
	return (last + offset) <= time.Now().Unix()
}

func (s *SmartctlContext) scan() error {
	var err error
	if s.lastScan == 0 || s.isAfter(s.lastScan, scanMinFrequency) {
		s.devices, err = scanDevices()
		s.lastScan = time.Now().Unix()
	}
	return err
}

func (s *SmartctlContext) startShortTests() error {
	if err := s.scan(); err != nil {
		return err
	}
	for _, d := range s.devices {
		if err := startShortTest(d); err != nil {
			return err
		}
	}
	return nil
}

func (s *SmartctlContext) GetReports() ([]*SmartctlOutput, error) {
	if err := s.scan(); err != nil {
		return nil, err
	}

	if s.lastReports == 0 || s.isAfter(s.lastReports, getReportsFrequency) {
		var result []*SmartctlOutput
		for _, d := range s.devices {
			report, err := getHealthAndAttributes(d)
			if err != nil {
				return nil, err
			}
			result = append(result, report)
		}
		s.reports = result
		s.lastReports = time.Now().Unix()
	}
	return s.reports, nil
}

func (s *SmartctlContext) runMondayTasker() {
	for {
		now := time.Now()
		nextMonday := time.Date(now.Year(), now.Month(), now.Day(), 23, 0, 0, 0, now.Location())
		for nextMonday.Weekday() != time.Monday {
			nextMonday = nextMonday.AddDate(0, 0, 1)
		}
		if now.After(nextMonday) {
			nextMonday = nextMonday.AddDate(0, 0, 7)
		}

		duration := nextMonday.Sub(now)

		time.Sleep(duration)

		err := s.startShortTests()
		if err != nil {
			log.Print(err)
		}

		time.Sleep(1 * time.Minute)
	}
}
