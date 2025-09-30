.PHONY: help dev dev-down build build-fast setup migrate test clean final init seed generate-images

help:
	@echo "Available commands:"
	@echo "  dev        - Start development environment (fast, no monitoring)"
	@echo "  final      - Start final optimized environment (fast with monitoring)"
	@echo "  dev-down   - Stop development environment"
	@echo "  build      - Build Docker images (no cache)"
	@echo "  build-fast - Build Docker images (with cache)"
	@echo "  setup      - Setup project directories"
	@echo "  init       - Initialize database and run migrations"
	@echo "  migrate    - Run database migrations"
	@echo "  seed       - Seed database with all sample data"
	@echo "  seed-categories - Seed only categories"
	@echo "  seed-products   - Seed only products"
	@echo "  seed-users      - Seed only users"
	@echo "  seed-orders     - Seed only orders"
	@echo "  seed-reviews    - Seed only reviews"
	@echo "  generate-images - Generate placeholder images for products"
	@echo "  auto-init   - Full project setup (init + seed + images)"
	@echo "  start-full  - Build, start services and auto-initialize"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean up containers and volumes"

dev:
	@echo "Starting development environment..."
	@echo "Building images in parallel..."
	docker-compose build --parallel
	@echo "Starting services..."
	docker-compose up -d
	@echo "Development environment started!"
	@echo "Frontend: http://localhost:3000"
	@echo "Backend API: http://localhost:5000"
	@echo "API Docs: http://localhost:5000/docs"

final:
	@echo "Starting final optimized environment with multi-stage builds..."
	@echo "Building images with cache optimization..."
	docker-compose build --parallel
	@echo "Starting services..."
	docker-compose up -d
	@echo "Final environment started!"
	@echo "Frontend: http://localhost:3000"
	@echo "Backend API: http://localhost:5000"
	@echo "API Docs: http://localhost:5000/docs"
	@echo "Prometheus: http://localhost:9090"
	@echo "Grafana: http://localhost:3001 (admin/admin)"

dev-down:
	docker-compose down

build:
	@echo "Building optimized Docker images with multi-stage builds..."
	docker-compose build --parallel --no-cache

build-fast:
	@echo "Building Docker images with cache optimization..."
	docker-compose build --parallel

setup:
	@echo "Setting up E-Commerce project..."
ifeq ($(OS),Windows_NT)
	@if not exist .env ( \
		echo Creating .env file from template... && \
		copy env.example .env && \
		echo Please edit .env file with your configuration \
	)
	@if not exist logs mkdir logs
	@if not exist logs\backend mkdir logs\backend
	@if not exist logs\nginx mkdir logs\nginx
	@if not exist logs\frontend mkdir logs\frontend
	@if not exist backend-go\uploads mkdir backend-go\uploads
	@if not exist nginx\ssl mkdir nginx\ssl
else
	@if [ ! -f .env ]; then \
		echo "Creating .env file from template..."; \
		cp env.example .env; \
		echo "Please edit .env file with your configuration"; \
	fi
	@mkdir -p logs/backend logs/nginx logs/frontend
	@mkdir -p backend-go/uploads
	@mkdir -p nginx/ssl
	@chmod +x scripts/*.sh
endif
	@echo "Setup completed!"

init:
	@echo "Initializing database..."
	cd backend-go && go run cmd/main.go -mode=init -wait

migrate:
	docker-compose exec backend ./main -mode=init

seed:
	@echo "Seeding database with all sample data..."
	cd backend-go && go run cmd/main.go -mode=seed

generate-images:
	@echo "Generating placeholder images for products..."
	cd backend-go && go run cmd/main.go -mode=generate-images
	@echo "Placeholder images generated successfully!"
	@echo "Images are available at: http://localhost:5000/api/uploads/"

auto-init:
	@echo "Auto-initializing E-Commerce project..."
	@echo "This will: initialize DB, seed data, and generate placeholder images"
	cd backend-go && go run cmd/main.go -mode=auto-init -wait
	@echo "Project ready! Visit http://localhost:3000"

start-full:
	@echo "Starting E-Commerce with full initialization..."
	@echo "Building and starting services..."
	docker-compose up -d --build
	@echo "Waiting for services to be ready..."
	timeout /t 10 /nobreak >nul 2>&1 || sleep 10
	@echo "Running auto-initialization..."
	cd backend-go && go run cmd/main.go -mode=auto-init -wait
	@echo "E-Commerce is ready!"
	@echo "Frontend: http://localhost:3000"
	@echo "Backend: http://localhost:5000"
	@echo "Admin: http://localhost:5000/admin"

test:
	@echo "Running tests..."
	cd backend-go && go test ./...

clean:
	docker-compose down -v --remove-orphans
	docker system prune -f