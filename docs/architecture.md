# Architecture Documentation

## Overview

This document summarizes the finalized Minimum Viable Product (MVP) for a simulated secure software supply chain platform. Inspired by tools like Chainguard, Cosign, Sigstore, and SLSA, the MVP provides a demo-capable system that models how organizations can build, verify, distribute, and track container images in a secure, traceable, and auditable manner.

The system is built using containerized services orchestrated via `docker-compose`, and is suitable for demo environments, internal prototyping, or design validation.

---

## Product Vision

The platform allows organizations (tenants) to:

- Simulate secure image builds and distribution pipelines.
- Track metadata including SBOMs, provenance, attestations, and CVEs.
- Model SLA compliance for CVE remediation.
- Observe image lineage, tenant usage, and build provenance.
- Serve as an educational or decision-support tool to evaluate secure supply chain needs.

---

## Key MVP Goals

- **Security-First Design**: Images carry associated metadata like signed SBOMs, provenance, and build attestations (simulated).
- **Traceability**: Builds are linked to source events (e.g. upstream changes, CVE discovery, or scheduled rebuilds).
- **Multi-Tenant Support**: Tenants are logically isolated and tracked across the system using a `tenant_namespace`.
- **SBOM Handling**: SBOMs are generated in both SPDX and CycloneDX formats and stored in an accessible object store (MinIO).
- **CVE Monitoring**: Each image is scanned for known vulnerabilities, triggering SLA timers and alerting logic.
- **Explorability**: A minimalistic web frontend supports exploration of images, metadata, and CVE details.
- **Simulated Systems**: All integrations with external systems (e.g. Rekor, Sigstore, CVE databases) are stubbed or mocked.

---

## Architecture Overview

### Services & Components

| Component                 | Description                                                                 |
|---------------------------|-----------------------------------------------------------------------------|
| **Backend (Go)**          | Serves the REST API, implements business logic using Clean Architecture.    |
| **Metadata DB (PostgreSQL)**| Stores all structured image, build, CVE, and tenant data.                   |
| **Object Store (MinIO)**  | Stores build logs, SBOM files, attestations, and provenance documents.      |
| **Frontend (Nginx)**      | Serves the static Vanilla JS frontend for the user interface.               |

All services are containerized and orchestrated via `docker-compose`.

---

## Data Model Highlights

### Core Tables

- `images`: Image metadata including digest, tags, associated SBOMs, and SLSA level.
- `build_events`: Track origin (e.g. upstream trigger, CVE), build status, and timestamps.
- `sbom_records`: SBOM format, URI, associated image.
- `cve_findings`: CVEs tied to specific images, with fix availability and risk rating.
- `attestations`: Simulated signed metadata linking image, build, and SBOMs.
- `customers`: Registered tenants; includes metadata like contact info and SLA tier.
- `customer_image_usage`: Maps image usage per tenant, version pinning, and runtime info.
- `sla_violations`: Tracks when CVEs exceed SLA deadlines and trigger escalations.
- `notifications`: Sent alerts related to builds, CVEs, or SLA breaches.
- `audit_logs`: Signed log entries of sensitive actions (e.g. CVE registration, signature generation).

### Design Considerations

- All metadata is tenant-aware via `tenant_namespace`.
- Rich join queries support lineage analysis (e.g. "Which customers use an image affected by CVE-1234?").
- Timestamped metadata supports auditing and SLA monitoring.

---

## Simulated Workflows

### 1. Build Trigger → Image Publication

- **Triggered by**:
  - Upstream code change (mock webhook via `POST /v1/webhooks/upstream`)
  - Manual API call (`POST /v1/images`)
  - Periodic rebuild (simulated cron)
  - CVE discovery in base image
- **Process**:
  - An orchestrator service (mocked) simulates a build.
  - A `build_event` is created in `pending` state.
  - Upon completion (simulated), the orchestrator calls back to the API to finalize the build.
  - Artifacts (logs, SBOMs, attestations) are stored in MinIO.
  - The database is updated with image, build, and artifact metadata.
  - The `build_event` is updated to `completed`.

### 2. CVE Detection → SLA Escalation

- A CVE scanner (mocked) assigns mock CVEs to affected images.
- An SLA timer is calculated based on severity and customer SLA policy.
- If a CVE is not fixed in time, an entry in `sla_violations` is created, triggering a simulated alert.
- The UI or API can list current SLA violations by tenant or severity.

---

## Security Model (Simulated)

| Area              | Approach                                                                 |
|-------------------|--------------------------------------------------------------------------|
| **Auth**          | Mock API keys and JWT for API access.                                    |
| **RBAC**          | Basic tenant/user scoping via `tenant_namespace`.                        |
| **Signing**       | Mocked Cosign-style signatures for images and attestations.              |
| **SBOM Format**   | SPDX and CycloneDX stored in MinIO, referenced via URIs.                 |
| **Audit Logging** | Log entries in `audit_logs` per sensitive event.                         |
| **Vuln Scan**     | Simulated CVE ingestion with timestamps, fixability, and CVSS scoring.   |

> **Note:** All security flows are simulated for the MVP and are not suitable for production use.

---

## Frontend UI

- Built with plain JavaScript, no framework.
- **Pages**:
  - **Image Explorer**: List and filter available images.
  - **Image Detail View**: Show SBOM links, CVEs, provenance data.
- Static assets are served via Nginx.

---

## Out-of-Scope for MVP

- Real-time CVE feeds (e.g. NVD, OSV).
- Production-ready cryptographic signature validation.
- OCI registry push/pull integrations.
- Role-based UI access.
- Complex frontend features (e.g. search suggestions, charts).
```
