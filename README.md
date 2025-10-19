# Project: Secure Software Supply Chain MVP

This project is a Minimum Viable Product (MVP) for a simulated secure software supply chain platform. Inspired by tools like Chainguard, Cosign, Sigstore, and SLSA, it provides a demo-capable system that models how organizations can build, verify, distribute, and track container images in a secure, traceable, and auditable manner.

The system is built using containerized services orchestrated via `docker-compose`, and is suitable for demo environments, internal prototyping, or design validation.

## Features

- **Security-First Design**: Images carry associated metadata like signed SBOMs, provenance, and build attestations (simulated).
- **Traceability**: Builds are linked to source events (e.g. upstream changes, CVE discovery, or scheduled rebuilds).
- **Multi-Tenant Support**: Tenants are logically isolated and tracked across the system.
- **SBOM Handling**: SBOMs are generated in both SPDX and CycloneDX formats and stored in an accessible object store (MinIO).
- **CVE Monitoring**: Each image is scanned for known vulnerabilities, triggering SLA timers and alerting logic.
- **Explorability**: A minimalistic web frontend supports exploration of images, metadata, and CVE details.
- **Simulated Systems**: All integrations with external systems (e.g. Rekor, Sigstore, CVE databases) are stubbed or mocked.

## Architecture

The system is composed of several containerized services:

- **Backend (Go)**: A Go application built using Clean Architecture principles, exposing a REST API.
- **Frontend (Vanilla JS)**: A lightweight, static frontend served by Nginx for exploring the data.
- **PostgreSQL**: The primary metadata database for all structured data.
- **MinIO**: An S3-compatible object store for artifacts like SBOMs, logs, and attestations.

All services are orchestrated using `docker-compose`.

For a detailed overview, see the [Architecture Documentation](./docs/architecture.md).

## Getting Started

### Prerequisites

- Docker
- Docker Compose

### Setup

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/your-username/secure-image-service.git
    cd secure-image-service
    ```

2.  **Configure environment variables:**
    Copy the example environment file. The default values are suitable for local development.
    ```sh
    cp .env.example .env
    ```

3.  **Start the services:**
    Use the provided `Makefile` to build and start all services.
    ```sh
    make up
    ```
    This command will start the backend, frontend, database, and MinIO in detached mode.

### Accessing the Services

- **Frontend UI**: [http://localhost](http://localhost)
- **Backend API**: [http://localhost:8080](http://localhost:8080)
- **MinIO Console**: [http://localhost:9001](http://localhost:9001) (Credentials: `minioadmin` / `minioadmin`)

## Development

The `Makefile` provides several commands to simplify the development workflow:

- `make up`: Build and start all services in detached mode.
- `make down`: Stop and remove all containers, networks, and volumes.
- `make logs`: Tail the logs for the `backend` service.
- `make db-reset`: A destructive command that completely stops the stack, wipes the database and MinIO volumes, and restarts the services. This is useful for re-applying the seed data from scratch.

## API

The backend exposes a RESTful API for interacting with the system. For detailed information on endpoints, requests, and responses, please see the [API Documentation](./docs/api.md).

A mock API key (`mock-api-key`) or JWT (`mock-jwt-token`) is required for authenticated endpoints.

**Example cURL:**
```sh
curl -H "X-API-Key: mock-api-key" http://localhost:8080/v1/images
```

## Project Structure

The project is organized as a monorepo:

```
/
├── .github/          # GitHub Actions CI workflows
├── backend/          # Go backend service (Clean Architecture)
├── frontend/         # Vanilla JS frontend (static assets)
├── deployments/      # Docker Compose, seed scripts, etc.
├── docs/             # Architecture and API documentation
├── .env.example      # Sample environment variables
├── Makefile          # Development helper commands
└── README.md         # This file
```

