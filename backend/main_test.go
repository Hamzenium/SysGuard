package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test toggleAlertHandler
func TestToggleAlertHandler(t *testing.T) {
	payload := []byte(`{"enable_alerts": true}`)
	req := httptest.NewRequest(http.MethodPost, "/toggle-alerts", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	toggleAlertHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %v", w.Code)
	}

	if !alertEnabled {
		t.Errorf("Expected alertEnabled to be true, got %v", alertEnabled)
	}
}

// Test toggleLimitHandler
func TestToggleLimitHandler(t *testing.T) {
	payload := []byte(`{"cpu_threshold": 80, "memory_threshold": 70, "disk_threshold": 75}`)
	req := httptest.NewRequest(http.MethodPost, "/limit-changer", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	toggleLimitHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %v", w.Code)
	}

	expectedLimits := SetLimit{
		CPUThreshold:    80,
		MemoryThreshold: 70,
		DiskThreshold:   75,
	}

	if defaultLimit != expectedLimits {
		t.Errorf("Expected limits %+v, got %+v", expectedLimits, defaultLimit)
	}
}

// Test resourceUsageHandler
func TestResourceUsageHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/resource-usage", nil)
	w := httptest.NewRecorder()

	resourceUsageHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %v", w.Code)
	}

	var usage ResourceUsage
	err := json.NewDecoder(w.Body).Decode(&usage)
	if err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	if usage.CPUUsage < 0 || usage.MemoryUsage < 0 || usage.DiskUsage < 0 {
		t.Errorf("Expected positive resource usage values, got %+v", usage)
	}
}

// Test sendAlert functionality (mock example)
func TestSendAlert(t *testing.T) {
	alertEnabled = true
	defaultLimit = SetLimit{
		CPUThreshold:    50,
		MemoryThreshold: 50,
		DiskThreshold:   50,
	}

	tests := []struct {
		resource string
		usage    float64
		alert    bool
	}{
		{"CPU", 60, true},
		{"Memory", 40, false},
		{"Disk", 55, true},
	}

	for _, test := range tests {
		mu.Lock()
		alertEnabled = false // Mocking disabled alert notifications
		mu.Unlock()
		sendAlert(test.resource, test.usage)
	}
}
