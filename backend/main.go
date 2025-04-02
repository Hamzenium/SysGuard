package main

import (
	"log"
	"net/http"

	"resource-monitor/backend/handlers"
	"resource-monitor/backend/monitor"
)

func main() {
	monitor.AlertEnabled = true
	log.Println("Default limits:", monitor.DefaultLimit)

	http.HandleFunc("/toggle-alerts", handlers.ToggleAlertHandler)
	http.HandleFunc("/resource-usage", handlers.ResourceUsageHandler)
	http.HandleFunc("/limit-changer", handlers.ToggleLimitHandler)
	http.HandleFunc("/shutdown", handlers.ShutdownHandler)

	server := &http.Server{Addr: ":8080"}

	go monitor.MonitorResources()

	go func() {
		log.Println("Starting backend server on port 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-handlers.ShutdownChan

	log.Println("Shutting down the server...")
	if err := server.Close(); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server stopped.")
}
