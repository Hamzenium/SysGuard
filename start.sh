#!/bin/sh

# Start backend in background
./backend-monitor &

sleep 1
./frontend-monitor
