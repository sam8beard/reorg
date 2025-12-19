#!/bin/bash
set -e 

echo -e "\n\n<---------------- [BUILD] ---------------->\n\n"

# Build Go file server
echo -e "[SERVER]\n\n"
echo -e "building server...\n" 
go build -o reorg-dev cmd/server/main.go
chmod +x reorg-dev

echo -e "server built successfully\n"


# Run app
echo -e "\n\n<---------------- [RUN] ---------------->\n\n"
echo -e "[APP]\n\n"
echo -e "running backend on port 8080...\n"
echo -e "running frontend on port 5173...\n"
EXE_PATH="$PWD/reorg-dev"
./reorg-dev & 
GO_PID=$!
trap 'kill $GO_PID; echo -e "\n\ncleaning reorg-dev and exiting..."; rm -f $EXE_PATH; exit' SIGINT SIGTERM
cd frontend && BROWSER=firefox npm run dev

