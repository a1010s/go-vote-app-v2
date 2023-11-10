#!/bin/bash
./myapp &
sleep 30
echo "Stopping the application on the server..."
pkill -f myapp # Kill the process
rm /root/go-vote-app-v2/myapp
rm -rf badger-db
