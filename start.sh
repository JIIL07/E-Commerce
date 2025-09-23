#!/bin/bash

echo "Starting E-Commerce Platform..."

# Check if .env exists
if [ ! -f .env ]; then
    echo "Creating .env file from template..."
    cp env.example .env
    echo ""
    echo "IMPORTANT: Please edit .env file with your configuration!"
    echo "Then run this script again."
    exit 1
fi

# Build and start all services
echo "Building Docker images..."
docker-compose build

echo "Starting all services..."
docker-compose up -d

echo ""
echo "========================================"
echo "E-Commerce Platform is starting!"
echo "========================================"
echo ""
echo "Frontend: http://localhost"
echo "Backend API: http://localhost/api"
echo "Health Check: http://localhost/health"
echo ""
echo "To view logs: docker-compose logs -f"
echo "To stop: docker-compose down"
echo ""
