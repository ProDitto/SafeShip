.PHONY: up down logs db-reset

up:
	docker-compose up -d --build

down:
	docker-compose down

logs:
	docker-compose logs -f backend

db-reset:
	docker-compose down -v
	docker-compose up -d db
	@echo "Waiting for database to be ready..."
	@sleep 5
	docker-compose up -d backend frontend minio

