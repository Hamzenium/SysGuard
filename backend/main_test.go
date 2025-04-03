package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"resource-monitor/backend/handlers"
	"resource-monitor/backend/monitor"
	"resource-monitor/backend/types"
)

func TestToggleAlertHandler(t *testing.T) {
	payload := []byte(`{"enable_alerts": true}`)
	req := httptest.NewRequest(http.MethodPost, "/toggle-alerts", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handlers.ToggleAlertHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %v", w.Code)
	}

	monitor.Mu.Lock()
	defer monitor.Mu.Unlock()
	if !monitor.AlertEnabled {
		t.Errorf("Expected AlertEnabled to be true, got %v", monitor.AlertEnabled)
	}
}

func TestToggleLimitHandler(t *testing.T) {
	payload := []byte(`{"cpu_threshold": 80, "memory_threshold": 70, "disk_threshold": 75}`)
	req := httptest.NewRequest(http.MethodPost, "/limit-changer", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handlers.ToggleLimitHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %v", w.Code)
	}

	monitor.Mu.Lock()
	defer monitor.Mu.Unlock()
	expected := types.SetLimit{CPUThreshold: 80, MemoryThreshold: 70, DiskThreshold: 75}
	if monitor.DefaultLimit != expected {
		t.Errorf("Expected %+v, got %+v", expected, monitor.DefaultLimit)
	}
}

func TestResourceUsageHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/resource-usage", nil)
	w := httptest.NewRecorder()

	handlers.ResourceUsageHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %v", w.Code)
	}

	var usage types.ResourceUsage
	if err := json.NewDecoder(w.Body).Decode(&usage); err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	if usage.CPUUsage < 0 || usage.MemoryUsage < 0 || usage.DiskUsage < 0 {
		t.Errorf("Expected positive resource values, got %+v", usage)
	}
}

// func TestSendAlertLogic(t *testing.T) {
// 	monitor.Mu.Lock()
// 	monitor.AlertEnabled = true
// 	monitor.DefaultLimit = types.SetLimit{
// 		CPUThreshold:    50,
// 		MemoryThreshold: 50,
// 		DiskThreshold:   50,
// 	}
// 	monitor.Mu.Unlock()

// 	// These will trigger alerts (but we don't test the actual notification pop-up)
// 	monitor.SendAlert("CPU", 55)
// 	monitor.SendAlert("Memory", 30) // Should not trigger alert
// 	monitor.SendAlert("Disk", 70)
// }
