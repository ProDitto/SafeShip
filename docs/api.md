# API Documentation

This document provides details on the REST API endpoints for the Secure Software Supply Chain MVP.

## Authentication

All endpoints under the `/v1` path require authentication. For this MVP, authentication is mocked and can be provided in one of two ways:

1.  **API Key**: Include an `X-API-Key` header with the value `mock-api-key`.
2.  **Bearer Token**: Include an `Authorization` header with the value `Bearer mock-jwt-token`.

**Example:**
```sh
curl -H "X-API-Key: mock-api-key" http://localhost:8080/v1/images
```

---

## Endpoints

### Images

#### List Images

- **Endpoint**: `GET /v1/images`
- **Description**: Retrieves a list of all published images.
- **Success Response (200 OK)**:
  ```json
  [
    {
      "id": 1,
      "tenant_namespace": "customer-a",
      "digest": "sha256:c3d3b3c3d3...",
      "tags": ["1.0.0", "latest"],
      "slsa_level": 3,
      "created_at": "2023-10-27T10:00:00Z",
      "updated_at": "2023-10-27T10:00:00Z"
    }
  ]
  ```

#### Get Image by ID

- **Endpoint**: `GET /v1/images/{id}`
- **Description**: Retrieves a single image by its unique ID.
- **Success Response (200 OK)**:
  ```json
  {
    "id": 1,
    "tenant_namespace": "customer-a",
    "digest": "sha256:c3d3b3c3d3...",
    "tags": ["1.0.0", "latest"],
    "slsa_level": 3,
    "created_at": "2023-10-27T10:00:00Z",
    "updated_at": "2023-10-27T10:00:00Z"
  }
  ```
- **Error Response (404 Not Found)**: If the image does not exist.

#### Trigger New Image Build

- **Endpoint**: `POST /v1/images`
- **Description**: Simulates triggering a new image build. This creates a `build_event` record with a `pending` status and starts an asynchronous (mocked) build process.
- **Request Body**:
  ```json
  {
    "tenant_namespace": "customer-a"
  }
  ```
- **Success Response (202 Accepted)**:
  ```json
  {
    "id": 5,
    "tenant_namespace": "customer-a",
    "image_id": null,
    "trigger_type": "api",
    "status": "pending",
    "created_at": "2023-10-27T12:00:00Z",
    "updated_at": "2023-10-27T12:00:00Z"
  }
  ```

#### Get Image SBOMs

- **Endpoint**: `GET /v1/images/{id}/sbom`
- **Description**: Retrieves SBOM metadata for a given image. (Currently mocked).
- **Success Response (200 OK)**:
  ```json
  [
    { "format": "SPDX", "uri": "minio://sboms/..." },
    { "format": "CycloneDX", "uri": "minio://sboms/..." }
  ]
  ```

#### Get Image CVEs

- **Endpoint**: `GET /v1/images/{id}/cves`
- **Description**: Retrieves CVE findings for a given image. (Currently mocked).
- **Success Response (200 OK)**:
  ```json
  [
    { "cve_id": "CVE-2023-1111", "severity": "Critical", "fix_available": true },
    { "cve_id": "CVE-2023-2222", "severity": "High", "fix_available": false }
  ]
  ```

#### Get Image Verification Data

- **Endpoint**: `GET /v1/images/{id}/verify`
- **Description**: Retrieves attestation and verification metadata for a given image. (Currently mocked).
- **Success Response (200 OK)**:
  ```json
  {
    "signature": { "key_id": "cosign-key-1", "rekor_entry_uri": "https://rekor.mock.dev/12345" },
    "attestations": [
      { "type": "provenance", "uri": "minio://attestations/..." },
      { "type": "vuln-scan", "uri": "minio://attestations/..." }
    ]
  }
  ```

### Webhooks

#### Trigger Upstream Build

- **Endpoint**: `POST /v1/webhooks/upstream`
- **Description**: A dedicated endpoint for triggering a build from a simulated upstream event (e.g., a Git commit).
- **Request Body**:
  ```json
  {
    "tenant_namespace": "customer-b"
  }
  ```
- **Success Response (202 Accepted)**: Returns the created `build_event` object, similar to `POST /v1/images`.

### Customers

#### List Customers

- **Endpoint**: `GET /v1/customers`
- **Description**: Retrieves a list of all tenants/customers.
- **Success Response (200 OK)**:
  ```json
  [
    {
      "namespace": "customer-a",
      "name": "Customer A Inc.",
      "contact_info": "contact@customera.com",
      "sla_tier": "premium",
      "created_at": "2023-10-27T09:00:00Z",
      "updated_at": "2023-10-27T09:00:00Z"
    }
  ]
  ```

#### Get Customer by Namespace

- **Endpoint**: `GET /v1/customers/{namespace}`
- **Description**: Retrieves a single customer by their unique namespace.
- **Success Response (200 OK)**: Returns a single customer object.
- **Error Response (404 Not Found)**: If the customer does not exist.

