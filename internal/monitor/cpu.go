package monitor

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

type CPUStats struct {
	UsagePercent float64
	Cores        int
	Model        string
}

func GetCPU() (*CPUStats, error) {
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, fmt.Errorf("get CPU usage: %w", err)
	}

	if len(percentages) == 0 {
		return nil, fmt.Errorf("CPU usage data unavailable")
	}

	info, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("get CPU information: %w", err)
	}

	stats := &CPUStats{
		UsagePercent: percentages[0],
		Cores:        len(info),
	}

	if len(info) > 0 {
		stats.Model = info[0].ModelName
	}

	return stats, nil
}
