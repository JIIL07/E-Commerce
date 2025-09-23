# Project Management
.PHONY: help setup build up down logs clean restart migrate-up migrate-down

help:
	@echo "Available commands:"
	@echo "  setup     - Initial project setup"
	@echo "  build     - Build all Docker images"
	@echo "  up        - Start all services"
	@echo "  down      - Stop all services"
	@echo "  logs      - View logs from all services"
	@echo "  clean     - Clean up containers and volumes"
	@echo "  restart   - Restart all services"
	@echo "  migrate-up   - Run database migrations"
	@echo "  migrate-down - Rollback database migrations"

setup:
	@echo "Setting up E-Commerce project..."
	@if [ ! -f .env ]; then \
		cp env.example .env; \
		echo "Created .env file from template. Please edit it with your configuration."; \
	fi
	@echo "Setup complete! Edit .env file and run 'make up' to start."

build:
	@echo "Building all Docker images..."
	docker-compose build

up:
	@echo "Starting all services..."
	docker-compose up -d
	@echo "Services started! Access the application at http://localhost"

down:
	@echo "Stopping all services..."
	docker-compose down

logs:
	docker-compose logs -f

clean:
	@echo "Cleaning up containers and volumes..."
	docker-compose down -v --remove-orphans
	docker system prune -f

restart:
	@echo "Restarting all services..."
	docker-compose restart

migrate-up:
	@echo "Running database migrations..."
	docker-compose exec backend go run ./cmd/server migrate-up

migrate-down:
	@echo "Rolling back database migrations..."
	docker-compose exec backend go run ./cmd/server migrate-down

# Development commands
dev-backend:
	@echo "Starting backend in development mode..."
	cd backend-go && make dev

dev-frontend:
	@echo "Starting frontend in development mode..."
	cd frontend && npm run dev

# Production commands
prod-build:
	@echo "Building for production..."
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml build

prod-up:
	@echo "Starting production services..."
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Database commands
db-shell:
	docker-compose exec postgres psql -U postgres -d ecommerce

db-backup:
	@echo "Creating database backup..."
	docker-compose exec postgres pg_dump -U postgres ecommerce > backup_$(shell date +%Y%m%d_%H%M%S).sql

# Monitoring
status:
	@echo "Service status:"
	docker-compose ps

health:
	@echo "Health check:"
	@curl -s http://localhost/api/health | jq . || echo "Health check failed"
