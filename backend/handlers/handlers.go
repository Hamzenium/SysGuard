package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"resource-monitor/backend/types"
	"resource-monitor/backend/utils"

	"resource-monitor/backend/monitor"
)

var ShutdownChan = make(chan struct{})

func ResourceUsageHandler(w http.ResponseWriter, r *http.Request) {
	usage := types.ResourceUsage{
		CPUUsage:    utils.GetCPUUsage(),
		MemoryUsage: utils.GetMemoryUsage(),
		DiskUsage:   utils.GetDiskUsage(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usage)
}

func ToggleAlertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		EnableAlerts bool `json:"enable_alerts"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	monitor.Mu.Lock()
	monitor.AlertEnabled = data.EnableAlerts
	monitor.Mu.Unlock()

	log.Printf("Alerts enabled: %v", data.EnableAlerts)
	w.WriteHeader(http.StatusOK)
}

func ToggleLimitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var data types.SetLimit
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	monitor.Mu.Lock()
	monitor.DefaultLimit = data
	monitor.Mu.Unlock()

	log.Printf("Updated limits - CPU: %.2f, Memory: %.2f, Disk: %.2f", data.CPUThreshold, data.MemoryThreshold, data.DiskThreshold)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func ShutdownHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		log.Println("Received shutdown request")
		close(ShutdownChan)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
