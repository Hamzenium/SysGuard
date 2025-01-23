#!/bin/bash

# Get the script's directory
SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)

# Start backend monitor
echo "Starting backend monitor..."
if [ -f "$SCRIPT_DIR/backend/monitor" ]; then
    echo "Found 'monitor' executable in $SCRIPT_DIR/backend."
    "$SCRIPT_DIR/backend/monitor" &
    backend_pid=$!
    echo "Backend monitor started with PID: $backend_pid"
else
    echo "Error: 'monitor' file not found in $SCRIPT_DIR/backend."
    exit 1
fi

# Start frontend SysGuard
echo "Starting frontend SysGuard..."
if [ -f "$SCRIPT_DIR/frontend/monitor" ]; then
    echo "Found 'sysguard' executable in $SCRIPT_DIR/frontend."
    "$SCRIPT_DIR/frontend/monitor" &
    frontend_pid=$!
    echo "Frontend SysGuard started with PID: $frontend_pid"
else
    echo "Error: 'sysguard' file not found in $SCRIPT_DIR/frontend."
    echo "Stopping backend monitor..."
    kill $backend_pid
    exit 1
fi

# Wait for both processes to complete
echo "Both backend and frontend are running. Press Ctrl+C to stop."
trap "echo 'Stopping processes...'; kill $backend_pid $frontend_pid" EXIT
wait
