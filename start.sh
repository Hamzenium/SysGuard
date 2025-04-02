#!/bin/sh

./backend-monitor &

sleep 1
./frontend-monitor
