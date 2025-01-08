.PHONY: build tinywiny down logs

# Build the Docker images
build:
	docker compose build

# Start the application and Redis server in the background
tinywiny:
	docker compose up -d

# Stop and remove all containers
down:
	docker compose down

# Show logs of both services
logs:
	docker compose logs -f
