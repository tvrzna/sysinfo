package smartctl

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
)

var pathSmartctl string
var pathSudo string

func getExec(args ...string) (*exec.Cmd, error) {
	var err error

	if pathSmartctl == "" {
		pathSmartctl, err = exec.LookPath("smartctl")
		if err != nil {
			return nil, errors.New("smartctl is not installed")
		}
	}

	if os.Getegid() != 0 {
		if pathSudo == "" {
			pathSudo, err = exec.LookPath("sudo")
			if err != nil || pathSudo == "" {
				// Try doas, if sudo is not available
				pathSudo, err = exec.LookPath("doas")
				if err != nil {
					return nil, errors.New("sudo nor doas is available")
				}
			}
		}
	}

	command := pathSmartctl
	if pathSudo != "" {
		command = pathSudo
		args = slices.Insert(args, 0, pathSmartctl)
	}

	args = append(args, "-j")

	return exec.Command(command, args...), nil
}

func runSmartctl(args ...string) (*SmartctlOutput, error) {
	cmd, err := getExec(args...)
	if err != nil {
		return nil, err
	}

	json, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	result := &SmartctlOutput{}
	if err := result.parse(json); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return result, nil

}

func scanDevices() ([]string, error) {
	var result []string

	output, err := runSmartctl("--scan")
	if err != nil {
		return nil, err
	}

	for _, d := range output.Devices {
		result = append(result, d.Name)
	}

	return result, nil
}

func startShortTest(device string) error {
	output, err := runSmartctl("-t", "short", device)
	if err != nil {
		return err
	}
	if output.Smartctl.ExitStatus != 0 {
		return fmt.Errorf("could not start short test on device '%s'", device)
	}
	return nil
}

func getHealthAndAttributes(device string) (*SmartctlOutput, error) {
	output, err := runSmartctl("-AHi", device)
	if err != nil {
		return nil, err
	}
	return output, nil
}
