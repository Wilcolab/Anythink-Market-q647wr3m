.PHONY: help run build test clean docker-up docker-down seed

# Show help
help:
	@echo "Available commands:"
	@echo "  run          - Run the application"
	@echo "  build        - Build the application"
	@echo "  test         - Run tests (when available)"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-up    - Start PostgreSQL with Docker Compose"
	@echo "  docker-down  - Stop PostgreSQL containers"
	@echo "  seed         - Run database migrations and seeding"

# Run the application
run:
	go run .

# Build the application
build:
	go build -o bin/quiz-api .

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Start PostgreSQL with Docker Compose
docker-up:
	docker-compose up -d postgres
	@echo "PostgreSQL is starting... use 'docker-compose logs -f postgres' to check status"

# Stop Docker containers
docker-down:
	docker-compose down

# Force seed the database (useful for development)
seed:
	@echo "This will connect to the database and run migrations + seeding"
	@echo "Make sure PostgreSQL is running and then run: go run .
