#!/bin/sh

# Start backend
echo "Starting backend..."
./backend/sysguard &   
backend_pid=$!

sleep 1

echo "Starting frontend..."
open ./frontend/build/Release/SysGuard.app

echo "Backend running with PID $backend_pid. Press Ctrl+C to stop backend."
trap "echo 'Stopping backend...'; kill $backend_pid" EXIT
wait $backend_pid
