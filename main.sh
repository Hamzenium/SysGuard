#!/bin/bash

# Get the script's directory
SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)

echo "Starting backend monitor..."
if [ -f "$SCRIPT_DIR/backend/sysguard" ]; then
    echo "Found 'sysguard' backend executable in $SCRIPT_DIR/backend."
    "$SCRIPT_DIR/backend/sysguard" &
    backend_pid=$!
    echo "Backend monitor started with PID: $backend_pid"
else
    echo "Error: 'sysguard' backend executable not found in $SCRIPT_DIR/backend."
    exit 1
fi

echo "Starting frontend SysGuard..."
if [ -d "$SCRIPT_DIR/frontend/build/Release/SysGuard.app" ]; then
    echo "Found 'SysGuard.app' in $SCRIPT_DIR/frontend/build/Release."
    open "$SCRIPT_DIR/frontend/build/Release/SysGuard.app"
    echo "Frontend SysGuard launched (runs in foreground by macOS)."
else
    echo "Error: 'SysGuard.app' not found in $SCRIPT_DIR/frontend/build/Release."
    echo "Stopping backend monitor..."
    kill $backend_pid
    exit 1
fi

echo "Backend monitor is running with PID $backend_pid. Press Ctrl+C to stop backend."
trap "echo 'Stopping backend...'; kill $backend_pid" EXIT
wait $backend_pid
