#!/bin/bash

# Docker-based database seeding script
# This script runs the seed command inside a Docker container

set -e

echo "Starting database seeding with Docker..."

# Check if docker-compose.yml exists
if [ ! -f docker-compose.yml ]; then
    echo "Error: docker-compose.yml not found."
    echo "Please run this script from the project root directory."
    exit 1
fi

# Build and run the seed command in Docker
echo "Building and running seed command in Docker..."
docker compose run --rm app go run cmd/seed/main.go

echo "Docker-based database seeding completed!"