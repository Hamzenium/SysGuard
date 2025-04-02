package monitor

import (
	"fmt"
	"log"
	"resource-monitor/backend/types"
	"resource-monitor/backend/utils"
	"sync"
	"time"

	"github.com/gen2brain/beeep"
)

var (
	AlertEnabled bool
	Mu           sync.Mutex
	DefaultLimit = types.SetLimit{CPUThreshold: 90, MemoryThreshold: 90, DiskThreshold: 90}
)

func sendAlert(resource string, usage float64) {
	Mu.Lock()
	defer Mu.Unlock()

	if !AlertEnabled {
		return
	}

	var threshold float64
	switch resource {
	case "CPU":
		threshold = DefaultLimit.CPUThreshold
	case "Memory":
		threshold = DefaultLimit.MemoryThreshold
	case "Disk":
		threshold = DefaultLimit.DiskThreshold
	}

	if usage > threshold {
		err := beeep.Alert(fmt.Sprintf("%s Alert", resource), fmt.Sprintf("%s usage is at %.2f%%!", resource, usage), "")
		if err != nil {
			log.Println("Error sending macOS notification:", err)
		}
	}
}

func MonitorResources() {
	for {
		cpuUsage := utils.GetCPUUsage()
		memUsage := utils.GetMemoryUsage()
		diskUsage := utils.GetDiskUsage()

		sendAlert("CPU", cpuUsage)
		sendAlert("Memory", memUsage)
		sendAlert("Disk", diskUsage)

		time.Sleep(5 * time.Second)
	}
}
