# Binary name for backend
BINARY_NAME = sysguard

GO = go
GOFMT = gofmt

BACKEND_DIR = backend
FRONTEND_DIR = MySwiftApp
BACKEND_LOG = backend.log
BACKEND_PORT = 8080

# Install Go dependencies
deps:
	@echo "Installing Go dependencies..."
	$(GO) mod tidy

# Build backend
build-backend:
	@echo "Building the backend application..."
	$(GO) build -o $(BACKEND_DIR)/$(BINARY_NAME) $(BACKEND_DIR)/main.go

# Build frontend (Swift app)
build-frontend:
	@echo "Building the frontend Swift application..."
	cd MySwiftApp && swift build -c release

# Build both
build: build-backend build-frontend

# Run backend
run-backend:
	@echo "Running the backend server..."
	@$(BACKEND_DIR)/$(BINARY_NAME) > $(BACKEND_LOG) 2>&1 &
	@echo "Backend started in the background."

# Wait for backend
wait-backend:
	@echo "Waiting for backend to start..."
	@while ! nc -z localhost $(BACKEND_PORT); do sleep 1; done
	@echo "Backend is running!"

# Run frontend (Swift app)
run-frontend:
	@echo "Running the frontend application..."
	cd MySwiftApp && .build/release/MySwiftApp

# Run everything
run-all: deps build run-backend wait-backend run-frontend

# Run using pre-built
run-built: run-backend wait-backend run-frontend

# Clean
clean:
	@echo "Cleaning up..."
	rm -f $(BACKEND_DIR)/$(BINARY_NAME) $(BACKEND_LOG)
	@echo "Cleaning frontend build..."
	cd $(FRONTEND_DIR) && xcodebuild clean

# Help
help:
	@echo "Makefile for SysGuard project"
	@echo ""
	@echo "Available targets:"
	@echo "  deps         - Install Go dependencies"
	@echo "  build        - Build both backend and frontend applications"
	@echo "  build-backend- Build backend Go server"
	@echo "  build-frontend - Build Swift frontend"
	@echo "  run-backend  - Run backend server"
	@echo "  wait-backend - Wait for backend to become ready"
	@echo "  run-frontend - Run frontend application"
	@echo "  run-all      - Install deps, build, run backend and frontend"
	@echo "  run-built    - Run backend and frontend using existing build"
	@echo "  clean        - Clean build artifacts and logs"
	@echo "  help         - Show this help message"
