#!/bin/bash

echo "Bring up containers"
docker-compose up -d

echo "Build service_a"
docker exec -ti service_a sh -c "which git || apk add git"
docker exec -ti service_a sh -c "cd /build && go build -o a ./service_a/main.go"

echo "Build service_b"
docker exec -ti service_b sh -c "which git || apk add git"
docker exec -ti service_b sh -c "cd /build && go build -o b ./service_b/main.go"

echo "Start service_a"
docker exec -d service_a sh -c "cd /build && ./a"

echo "Start service_b"
docker exec -d service_b sh -c "cd /build && ./b"
