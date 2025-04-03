BINARY_NAME = image

GO = go
GOFMT = gofmt

BACKEND_DIR = backend
FRONTEND_DIR = frontend
BACKEND_LOG = backend.log
BACKEND_PORT = 8080

deps:
	@echo "Installing Go dependencies..."
	@go mod tidy

build-backend:
	@echo "Building the backend application..."
	$(GO) build -o $(BACKEND_DIR)/$(BINARY_NAME) $(BACKEND_DIR)/main.go

build-frontend:
	@echo "Building the frontend application..."
	$(GO) build -o $(FRONTEND_DIR)/$(BINARY_NAME) $(FRONTEND_DIR)/widget.go

build: build-backend build-frontend

run-backend:
	@echo "Running the backend server..."
	@$(BACKEND_DIR)/$(BINARY_NAME) > $(BACKEND_LOG) 2>&1 &  # Run in background
	@echo "Backend started in the background."

wait-backend:
	@echo "Waiting for backend to start..."
	@while ! nc -z localhost $(BACKEND_PORT); do sleep 1; done
	@echo "Backend is running!"

run-frontend:
	@echo "Running the frontend application..."
	@$(FRONTEND_DIR)/$(BINARY_NAME)

run-all: deps build run-backend wait-backend run-frontend

run-built: run-backend wait-backend run-frontend

docker-run:
	@echo "Building binaries..."
	$(MAKE) build
	@echo "Building Docker image..."
	docker-compose build
	@echo "Starting container..."
	docker-compose up -d

clean:
	@echo "Cleaning up..."
	rm -f $(BACKEND_DIR)/$(BINARY_NAME) $(FRONTEND_DIR)/$(BINARY_NAME) $(BACKEND_LOG)

help:
	@echo "Makefile for Go application"
	@echo ""
	@echo "Available targets:"
	@echo "  build        - Build both backend and frontend applications"
	@echo "  run-backend  - Build and run the backend server"
	@echo "  wait-backend - Wait until the backend is ready using a real check"
	@echo "  run-frontend - Build and run the frontend application"
	@echo "  run-all      - Install deps, build everything, then run backend & frontend"
	@echo "  clean        - Remove generated binaries and logs"
	@echo "  deps         - Install Go dependencies and tidy up"
	@echo "  help         - Show this help message"
