package utils

import (
	"log"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func GetCPUUsage() float64 {
	percent, err := cpu.Percent(0, false)
	if err != nil {
		log.Printf("Error getting CPU usage: %v", err)
		return -1
	}
	return percent[0]
}

func GetMemoryUsage() float64 {
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Error getting memory usage: %v", err)
		return -1
	}
	return v.UsedPercent
}

func GetDiskUsage() float64 {
	d, err := disk.Usage("/")
	if err != nil {
		log.Printf("Error getting disk usage: %v", err)
		return -1
	}
	return d.UsedPercent
}
