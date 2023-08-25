package metric

import (
	"bufio"
	"os"
	"strings"
	"syscall"
)

type EnSizeUnit byte

const (
	pathMounts = "/proc/mounts"
	TB
)

type DiskUsage struct {
	Device    string
	Path      string
	UsedSize  uint64
	TotalSize uint64
}

func LoadDiskUsage() []*DiskUsage {
	supportedFS := []string{"ext2", "ext3", "ext4", "btrfs", "xfs"}
	var result []*DiskUsage

	f, err := os.Open(pathMounts)
	if err != nil {
		return nil
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), " ")
		if len(data) > 3 && arrContains(supportedFS, data[2]) {
			u := &DiskUsage{Device: data[0], Path: data[1]}

			var stat syscall.Statfs_t
			syscall.Statfs(u.Path, &stat)

			u.TotalSize = stat.Blocks * uint64(stat.Bsize)
			u.UsedSize = u.TotalSize - (stat.Bfree * uint64(stat.Bsize))

			result = append(result, u)
		}
	}

	return result
}

func arrContains(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}
