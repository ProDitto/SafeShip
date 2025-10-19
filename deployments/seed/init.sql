-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create customers table to represent tenants
CREATE TABLE IF NOT EXISTS customers (
    namespace VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    contact_info VARCHAR(255),
    sla_tier VARCHAR(50) NOT NULL DEFAULT 'standard',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create images table
CREATE TABLE IF NOT EXISTS images (
    id SERIAL PRIMARY KEY,
    tenant_namespace VARCHAR(255) NOT NULL REFERENCES customers(namespace) ON DELETE CASCADE,
    digest VARCHAR(255) UNIQUE NOT NULL,
    tags TEXT[],
    slsa_level INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_images_tenant_namespace ON images(tenant_namespace);

-- Table for SBOM records
CREATE TABLE IF NOT EXISTS sbom_records (
    id SERIAL PRIMARY KEY,
    image_id INT REFERENCES images(id) ON DELETE CASCADE,
    format VARCHAR(50) NOT NULL, -- e.g., 'SPDX', 'CycloneDX'
    uri VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_sbom_records_image_id ON sbom_records(image_id);

-- Table for CVE findings
CREATE TABLE IF NOT EXISTS cve_findings (
    id SERIAL PRIMARY KEY,
    image_id INT REFERENCES images(id) ON DELETE CASCADE,
    cve_id VARCHAR(50) NOT NULL,
    severity VARCHAR(50) NOT NULL, -- e.g., 'Critical', 'High', 'Medium', 'Low'
    description TEXT,
    fix_available BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_cve_findings_image_id ON cve_findings(image_id);

-- Table for attestations
CREATE TABLE IF NOT EXISTS attestations (
    id SERIAL PRIMARY KEY,
    image_id INT REFERENCES images(id) ON DELETE CASCADE,
    type VARCHAR(100) NOT NULL, -- e.g., 'provenance', 'slsa'
    uri VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_attestations_image_id ON attestations(image_id);


-- Seed data
TRUNCATE customers, images, sbom_records, cve_findings, attestations RESTART IDENTITY CASCADE;

INSERT INTO customers (namespace, name, contact_info, sla_tier) VALUES
('acme-corp', 'ACME Corporation', 'security@acme.corp', 'enterprise'),
('startup-inc', 'Startup Inc.', 'devops@startup.inc', 'standard');

INSERT INTO images (tenant_namespace, digest, tags, slsa_level) VALUES
('acme-corp', 'sha256:c3d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3', '{"latest", "1.2.3"}', 3),
('acme-corp', 'sha256:abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890', '{"1.2.2"}', 3),
('startup-inc', 'sha256:fedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321', '{"latest"}', 2),
('startup-inc', 'sha256:11223344556677889900aabbccddeeff11223344556677889900aabbccddeeff', '{"stable"}', 2);

-- Seed SBOMs for image 1
INSERT INTO sbom_records (image_id, format, uri) VALUES
(1, 'SPDX', 'minio://sboms/acme-corp/image-1.spdx.json'),
(1, 'CycloneDX', 'minio://sboms/acme-corp/image-1.cdx.json');

-- Seed CVEs for image 1
INSERT INTO cve_findings (image_id, cve_id, severity, description, fix_available) VALUES
(1, 'CVE-2023-4567', 'Critical', 'Remote code execution vulnerability in lib-xyz.', true),
(1, 'CVE-2023-8910', 'Medium', 'Denial of service in logging component.', false);

-- Seed Attestations for image 1
INSERT INTO attestations (image_id, type, uri) VALUES
(1, 'provenance', 'minio://attestations/acme-corp/image-1-provenance.json'),
(1, 'slsa-v1.0', 'minio://attestations/acme-corp/image-1-slsa.json');

-- Seed SBOMs for image 2
INSERT INTO sbom_records (image_id, format, uri) VALUES
(2, 'SPDX', 'minio://sboms/acme-corp/image-2.spdx.json');

-- Seed CVEs for image 3
INSERT INTO cve_findings (image_id, cve_id, severity, description, fix_available) VALUES
(3, 'CVE-2024-0001', 'High', 'SQL injection vulnerability in base image.', true);

-- Seed Attestations for image 3
INSERT INTO attestations (image_id, type, uri) VALUES
(3, 'provenance', 'minio://attestations/startup-inc/image-3-provenance.json');
