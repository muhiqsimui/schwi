package monitor

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/disk"
)

type DiskStats struct {
	Path        string
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
}

func GetDisk(path string) (*DiskStats, error) {
	usage, err := disk.Usage(path)
	if err != nil {
		return nil, fmt.Errorf("get disk usage: %w", err)
	}

	return &DiskStats{
		Path:        path,
		Total:       usage.Total,
		Used:        usage.Used,
		Free:        usage.Free,
		UsedPercent: usage.UsedPercent,
	}, nil
}
