BINARY_NAME = monitor

GO = go
GOFMT = gofmt

BACKEND_DIR = backend
FRONTEND_DIR = frontend

deps:
	@echo "Installing Go dependencies..."
	@go mod tidy

build-backend:
	@echo "Building the backend application..."
	$(GO) build -o $(BACKEND_DIR)/$(BINARY_NAME) $(BACKEND_DIR)/server.go

build-frontend:
	@echo "Building the frontend application..."
	$(GO) build -o $(FRONTEND_DIR)/$(BINARY_NAME) $(FRONTEND_DIR)/widget.go

run-backend: build-backend
	@echo "Running the backend server..."
	$(BACKEND_DIR)/$(BINARY_NAME)

run-frontend: build-frontend
	@echo "Running the frontend application..."
	$(FRONTEND_DIR)/$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	rm -f $(BACKEND_DIR)/$(BINARY_NAME) $(FRONTEND_DIR)/$(BINARY_NAME)

help:
	@echo "Makefile for Go application"
	@echo ""
	@echo "Available targets:"
	@echo "  build           - Build both backend and frontend applications"
	@echo "  run-backend     - Build and run the backend server"
	@echo "  run-frontend    - Build and run the frontend application"
	@echo "  clean           - Remove the generated binaries"
	@echo "  deps            - Install Go dependencies and tidy up"
	@echo "  help            - Show this help message"
