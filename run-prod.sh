#!/bin/bash
set -e 

echo -e "\n\n<---------------- [BUILD] ---------------->\n\n"

# Build Vite frontend
echo -e "[FRONTEND]\n"

echo -e "building dist...\n"
build_result=$(cd frontend && npm run build)
echo -e "$build_result\n"
sleep 2

echo -e "dist built successfully\n"

# Build Go file server
echo -e "[SERVER]\n\n"
echo -e "building server...\n" 
go build -o reorg-prod cmd/server/main.go
chmod +x reorg-prod
sleep 2

echo -e "server built successfully\n"


# Run app
echo -e "\n\n<---------------- [RUN] ---------------->\n\n"
echo -e "[APP]\n\n"
echo -e "running reorg on port 8081...\n"
trap 'echo -e "\n\ncleaning reorg-prod and exiting..."; rm -f reorg-prod; exit' SIGINT SIGTERM
./reorg-prod
