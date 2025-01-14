# Name of the binary file
BINARY_NAME = monitor

# Go-related variables
GO = go
GOFMT = gofmt

# Directories for backend and frontend
BACKEND_DIR = backend
FRONTEND_DIR = frontend

# Target for building the Go application (both backend and frontend)
build-backend:
	@echo "Building the backend application..."
	$(GO) build -o $(BACKEND_DIR)/$(BINARY_NAME) $(BACKEND_DIR)/server.go

build-frontend:
	@echo "Building the frontend application..."
	$(GO) build -o $(FRONTEND_DIR)/$(BINARY_NAME) $(FRONTEND_DIR)/widget.go

# Target for building both backend and frontend applications
build: build-backend build-frontend

# Target for running the backend server
run-backend: build-backend
	@echo "Running the backend server..."
	$(BACKEND_DIR)/$(BINARY_NAME)

# Target for running the frontend GUI
run-frontend: build-frontend
	@echo "Running the frontend application..."
	$(FRONTEND_DIR)/$(BINARY_NAME)

# Target for cleaning up generated files
clean:
	@echo "Cleaning up..."
	rm -f $(BACKEND_DIR)/$(BINARY_NAME) $(FRONTEND_DIR)/$(BINARY_NAME)

# Format the Go code
fmt:
	@echo "Formatting Go code..."
	$(GOFMT) -w $(BACKEND_DIR) $(FRONTEND_DIR)

# Install dependencies (if you are using Go modules)
deps:
	@echo "Installing dependencies..."
	$(GO) mod tidy

# Target to run everything: build, run, and clean up after
all: build run-backend run-frontend clean

# Help message for usage instructions
help:
	@echo "Makefile for Go application"
	@echo ""
	@echo "Available targets:"
	@echo "  build           - Build both backend and frontend applications"
	@echo "  run-backend     - Build and run the backend server"
	@echo "  run-frontend    - Build and run the frontend application"
	@echo "  clean           - Remove the generated binaries"
	@echo "  fmt             - Format Go source code"
	@echo "  deps            - Install Go dependencies"
	@echo "  all             - Build, run both backend and frontend, and clean"
	@echo "  help            - Show this help message"
