COMPOSE_FILE := ./deployments/docker-compose.yml

.PHONY: up down logs logs-backend logs-db db-reset help

help:
	@echo "Usage:"
	@echo "  make up           Start all services (build if needed)"
	@echo "  make down         Stop all services"
	@echo "  make logs         Tail backend logs"
	@echo "  make logs-backend Alias for logs"
	@echo "  make logs-db      Tail db service logs"
	@echo "  make db-reset     Reset DB volume & restart db service"

up:
	docker compose -f $(COMPOSE_FILE) up -d --build

down:
	docker compose -f $(COMPOSE_FILE) down

logs: logs-backend

logs-backend:
	docker compose -f $(COMPOSE_FILE) logs -f backend

logs-db:
	docker compose -f $(COMPOSE_FILE) logs -f db

db-reset:
	docker compose -f $(COMPOSE_FILE) down -v
	docker compose -f $(COMPOSE_FILE) up -d db
	@echo "Waiting for database to be ready..."
	@sleep 5
	docker compose -f $(COMPOSE_FILE) up -d backend frontend minio
