package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Structure to hold resource usage data
type ResourceUsage struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
}

// Function to fetch resource usage data from the backend API
func fetchResourceUsage() (*ResourceUsage, error) {
	resp, err := http.Get("http://localhost:8080/resource-usage")
	if err != nil {
		return nil, fmt.Errorf("error fetching resource usage: %v", err)
	}
	defer resp.Body.Close()

	var usage ResourceUsage
	if err := json.NewDecoder(resp.Body).Decode(&usage); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &usage, nil
}

// Function to toggle the alert status on the backend
func toggleAlerts(enabled bool) error {
	data := struct {
		EnableAlerts bool `json:"enable_alerts"`
	}{enabled}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %v", err)
	}

	resp, err := http.Post("http://localhost:8080/toggle-alerts", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error sending alert toggle request: %v", err)
	}
	defer resp.Body.Close()

	return nil
}

// Function to change the resource usage limits on the backend
func changeLimits(cpuThreshold, memoryThreshold, diskThreshold float64) error {
	data := struct {
		CPUThreshold    float64 `json:"cpu_threshold"`
		MemoryThreshold float64 `json:"memory_threshold"`
		DiskThreshold   float64 `json:"disk_threshold"`
	}{
		CPUThreshold:    cpuThreshold,
		MemoryThreshold: memoryThreshold,
		DiskThreshold:   diskThreshold,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %v", err)
	}

	resp, err := http.Post("http://localhost:8080/limit-changer", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error sending limit change request: %v", err)
	}
	defer resp.Body.Close()

	return nil
}

func createGraphWindow(a fyne.App) fyne.Window {
	w := a.NewWindow("SysGuard")

	cpuGraph := widget.NewProgressBar()
	cpuGraph.Min = 0
	cpuGraph.Max = 100

	memGraph := widget.NewProgressBar()
	memGraph.Min = 0
	memGraph.Max = 100

	diskGraph := widget.NewProgressBar()
	diskGraph.Min = 0
	diskGraph.Max = 100

	alertCheckbox := widget.NewCheck("Enable Alerts", func(enabled bool) {
		if err := toggleAlerts(enabled); err != nil {
			log.Println("Error toggling alerts:", err)
		}
	})
	alertCheckbox.SetChecked(true) // Set the checkbox to be initially checked

	// Add input fields for changing the limits
	cpuEntry := widget.NewEntry()
	cpuEntry.SetPlaceHolder("Enter CPU Threshold")

	memEntry := widget.NewEntry()
	memEntry.SetPlaceHolder("Enter Memory Threshold")

	diskEntry := widget.NewEntry()
	diskEntry.SetPlaceHolder("Enter Disk Threshold")

	// Button to submit new limits
	changeLimitsButton := widget.NewButton("Change Limits", func() {
		cpuThreshold, err := strconv.ParseFloat(cpuEntry.Text, 64)
		if err != nil {
			log.Println("Invalid CPU threshold")
			return
		}

		memThreshold, err := strconv.ParseFloat(memEntry.Text, 64)
		if err != nil {
			log.Println("Invalid Memory threshold")
			return
		}

		diskThreshold, err := strconv.ParseFloat(diskEntry.Text, 64)
		if err != nil {
			log.Println("Invalid Disk threshold")
			return
		}

		// Send the new limits to the backend
		if err := changeLimits(cpuThreshold, memThreshold, diskThreshold); err != nil {
			log.Println("Error changing limits:", err)
		} else {
			log.Println("Limits updated successfully!")
		}
	})

	go func() {
		for {
			usage, err := fetchResourceUsage()
			if err != nil {
				log.Println("Error fetching resource usage:", err)
				continue
			}

			// Update the graphs with the fetched resource data
			cpuGraph.SetValue(usage.CPUUsage)
			memGraph.SetValue(usage.MemoryUsage)
			diskGraph.SetValue(usage.DiskUsage)

			time.Sleep(time.Second)
		}
	}()

	content := container.NewVBox(
		widget.NewLabel("CPU Usage"),
		cpuGraph,
		widget.NewLabel("Memory Usage"),
		memGraph,
		widget.NewLabel("Disk Usage"),
		diskGraph,
		alertCheckbox, // Add the checkbox for enabling/disabling alerts
		widget.NewLabel("Change Resource Limits"),
		cpuEntry,
		memEntry,
		diskEntry,
		changeLimitsButton,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 300))

	return w
}

func main() {
	a := app.New()
	w := createGraphWindow(a)
	w.ShowAndRun()
}
