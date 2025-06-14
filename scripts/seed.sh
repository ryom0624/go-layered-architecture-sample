#!/bin/bash

# Database seeding script
# This script runs the seed command to populate the database with initial data

set -e

echo "Starting database seeding..."

# Check if .env file exists
if [ ! -f .env ]; then
    echo "Warning: .env file not found. Using default configuration."
    echo "Make sure to copy .env.example to .env and configure it properly."
fi

# Run the seed command
echo "Running seed command..."
go run cmd/seed/main.go

echo "Database seeding completed!"