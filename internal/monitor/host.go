package monitor

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/host"
)

type HostStats struct {
	Hostname        string
	OS              string
	Platform        string
	PlatformVersion string
	KernelVersion   string
	Uptime          uint64
}

func GetHost() (*HostStats, error) {
	info, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("get host information: %w", err)
	}

	uptime, err := host.Uptime()
	if err != nil {
		return nil, fmt.Errorf("get host uptime: %w", err)
	}

	return &HostStats{
		Hostname:        info.Hostname,
		OS:              info.OS,
		Platform:        info.Platform,
		PlatformVersion: info.PlatformVersion,
		KernelVersion:   info.KernelVersion,
		Uptime:          uptime,
	}, nil
}
