package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

const (
	CPUThreshold    = 80.0
	MemoryThreshold = 90.0
	DiskThreshold   = 80.0
)

type ResourceUsage struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
}

var alertEnabled bool
var mu sync.Mutex // To protect the alertEnabled variable

// Get CPU usage
func getCPUUsage() float64 {
	percent, err := cpu.Percent(0, false)
	if err != nil {
		log.Printf("Error getting CPU usage: %v", err)
		return -1
	}
	return percent[0]
}

// Get Memory usage
func getMemoryUsage() float64 {
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Error getting memory usage: %v", err)
		return -1
	}
	return v.UsedPercent
}

// Get Disk usage
func getDiskUsage() float64 {
	d, err := disk.Usage("/")
	if err != nil {
		log.Printf("Error getting disk usage: %v", err)
		return -1
	}
	return d.UsedPercent
}

// Send macOS notification if usage exceeds threshold
func sendAlert(resource string, usage, threshold float64) {
	mu.Lock()
	defer mu.Unlock()

	if alertEnabled && usage > threshold {
		err := beeep.Alert(fmt.Sprintf("%s Alert", resource), fmt.Sprintf("%s usage is at %.2f%%!", resource, usage), "")
		if err != nil {
			log.Println("Error sending macOS notification:", err)
		}
	}
}

// Toggle the alert status
func toggleAlertHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if r.Method == http.MethodPost {
		var data struct {
			EnableAlerts bool `json:"enable_alerts"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		alertEnabled = data.EnableAlerts
		log.Printf("Alerts enabled: %v", alertEnabled)
	}
	w.WriteHeader(http.StatusOK)
}

// Resource usage handler
func resourceUsageHandler(w http.ResponseWriter, r *http.Request) {
	usage := ResourceUsage{
		CPUUsage:    getCPUUsage(),
		MemoryUsage: getMemoryUsage(),
		DiskUsage:   getDiskUsage(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(usage); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func monitorResources() {
	for {
		cpuUsage := getCPUUsage()
		memUsage := getMemoryUsage()
		diskUsage := getDiskUsage()

		// Check if the resource usage exceeds thresholds and send notifications
		sendAlert("CPU", cpuUsage, CPUThreshold)
		sendAlert("Memory", memUsage, MemoryThreshold)
		sendAlert("Disk", diskUsage, DiskThreshold)

		time.Sleep(5 * time.Second)
	}
}

func main() {
	// Initial alert state
	alertEnabled = true

	// Set up routes
	http.HandleFunc("/toggle-alerts", toggleAlertHandler)
	http.HandleFunc("/resource-usage", resourceUsageHandler)

	// Start monitoring resources in the background
	go monitorResources()

	// Start the backend server
	log.Println("Starting backend server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
