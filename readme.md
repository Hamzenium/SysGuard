# Resource Monitor with Alerts
This project is a resource monitoring system built in Go (Golang) that tracks CPU, memory, disk usage, and temperature on your machine. The system periodically checks the resource usage and compares it against predefined thresholds. If any resource exceeds the threshold, a notification is triggered to alert the user. The application features a backend API and a graphical user interface (GUI) that can be used to monitor and manage these resources.

## Features

- **Real-time Resource Monitoring**: Monitors CPU, memory, disk, and system temperature usage and updates in real-time.
- **Alerts**: Sends notifications when resource usage exceeds configured thresholds.
- **Cross-Platform GUI**: A graphical user interface to monitor system resources and toggle alert notifications.
- **Backend API**: Provides an API to fetch resource usage data and control alert settings.

![image](https://github.com/user-attachments/assets/d1914265-8e60-481a-86d7-fc48cd7b6668)

## Design Overview

### Concurrency with Goroutines

Concurrency is a key design decision in this project. Go’s goroutines make it easy to run multiple tasks simultaneously, which is useful for both resource monitoring and alert handling. A dedicated goroutine monitors system resources like CPU, memory, disk, and temperature, checking their usage periodically. If any resource exceeds its configured threshold, an alert is triggered.

### Alert Management

Alerts are handled concurrently using mutexes to prevent race conditions when toggling the alert status. One goroutine runs the resource monitoring loop, while separate goroutines handle alerting and user interactions with the GUI.

### Frontend Updates

The frontend is updated in a separate goroutine that regularly polls the backend for updated resource usage data. This ensures the GUI reflects the latest data in real time, displaying live progress bars for each resource.

### Synchronization with Mutexes

Mutexes are employed to avoid race conditions when accessing shared resources like the `alertEnabled` variable. This ensures that the variable is safely accessed by multiple goroutines.

## Frontend GUI

The frontend GUI is built using Fyne, a cross-platform GUI library for Go. The GUI includes the following features:

- **Real-Time Monitoring**: Displays live progress bars for CPU, memory, disk, and temperature usage.
- **Alert Control**: A checkbox allows users to toggle alert notifications. When enabled, notifications are triggered when any resource exceeds its threshold.
- **Polling**: The frontend regularly polls the backend to update resource usage data.

## Backend API

The backend exposes two primary API endpoints:

- **GET /resource-usage**: Returns the current resource usage for CPU, memory, disk, and temperature in JSON format.
- **POST /toggle-alerts**: Allows the frontend to enable or disable alerts.

### Example Response for `/resource-usage`:
The response includes the current resource usage data in JSON format.

### Example Request for `/toggle-alerts`:
The request allows the frontend to toggle alerts by sending a JSON payload with the `enable_alerts` field.

## Error Handling and Logging

The project includes robust error handling, ensuring that system failures or invalid data are logged properly. Errors are logged using the `log` package, which facilitates tracing issues that may arise during execution.

## Installation & Setup

### Prerequisites

- Go 1.18 or higher
- Fyne (GUI library)
- `osx-cpu-temp` (for macOS users only, used for fetching CPU temperature)

### Steps to Run the Application

1. **Clone the Repository**:
    - Clone the repository and navigate to the project directory.

2. **Install Dependencies**:
    - Install the necessary Go dependencies by running the appropriate commands.

3. **Backend Setup**:
    - Navigate to the backend directory and run the Go server.

4. **Frontend Setup**:
    - Open a separate terminal, navigate to the frontend directory, and run the frontend application.

## Using the Makefile

You can use the Makefile to simplify the build and run process. Available commands include:

- **Install dependencies**: Install the required dependencies.
- **Build the backend**: Build the backend.
- **Build the frontend**: Build the frontend.
- **Run the backend**: Start the backend server.
- **Run the frontend**: Start the frontend GUI.
- **Clean up**: Remove generated binaries.
- **Help**: View available commands.

## Technologies Used

- **Go (Golang)**: Used for both the backend API and frontend GUI.
- **Fyne**: Cross-platform GUI library for Go.
- **beeep**: Used for sending macOS desktop notifications when resource usage exceeds thresholds.
- **gopsutil**: Used for gathering system statistics like CPU, memory, and disk usage.
- **osx-cpu-temp**: Used for fetching CPU temperature on macOS (macOS-specific).

## License

This project is licensed under the MIT License - see the LICENSE file for details.
