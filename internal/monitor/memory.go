package monitor

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/mem"
)

type MemoryStats struct {
	Total       uint64
	Used        uint64
	Available   uint64
	UsedPercent float64
}

func GetMemory() (*MemoryStats, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("get memory usage: %w", err)
	}

	return &MemoryStats{
		Total:       vm.Total,
		Used:        vm.Used,
		Available:   vm.Available,
		UsedPercent: vm.UsedPercent,
	}, nil
}
