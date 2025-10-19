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

-- Build Events Table
CREATE TABLE IF NOT EXISTS build_events (
    id SERIAL PRIMARY KEY,
    tenant_namespace VARCHAR(255) NOT NULL REFERENCES customers(namespace) ON DELETE CASCADE,
    image_id INT REFERENCES images(id) ON DELETE SET NULL, -- Allow image_id to be null if image is deleted
    trigger_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_build_events_tenant_namespace ON build_events(tenant_namespace);
CREATE INDEX IF NOT EXISTS idx_build_events_status ON build_events(status);

-- Table for SBOM records
CREATE TABLE IF NOT EXISTS sbom_records (
    id SERIAL PRIMARY KEY,
    image_id INT NOT NULL REFERENCES images(id) ON DELETE CASCADE,
    format VARCHAR(50) NOT NULL, -- e.g., 'SPDX', 'CycloneDX'
    uri VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_sbom_records_image_id ON sbom_records(image_id);

-- Table for CVE findings
CREATE TABLE IF NOT EXISTS cve_findings (
    id SERIAL PRIMARY KEY,
    image_id INT NOT NULL REFERENCES images(id) ON DELETE CASCADE,
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
    image_id INT NOT NULL REFERENCES images(id) ON DELETE CASCADE,
    type VARCHAR(100) NOT NULL, -- e.g., 'provenance', 'slsa'
    uri VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_attestations_image_id ON attestations(image_id);


-- Seed data
TRUNCATE customers, images, build_events, sbom_records, cve_findings, attestations RESTART IDENTITY CASCADE;

INSERT INTO customers (namespace, name, contact_info, sla_tier) VALUES
('acme-corp', 'ACME Corporation', 'contact@acme.com', 'enterprise'),
('startup-x', 'Startup X', 'devops@startupx.io', 'standard');

INSERT INTO images (tenant_namespace, digest, tags, slsa_level) VALUES
('acme-corp', 'sha256:c3d3e4f5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2c3', ARRAY['latest', '1.2.3'], 3),
('startup-x', 'sha256:d4e4f5g6h7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2c4d5', ARRAY['latest'], 1);

-- Seed SBOMs
INSERT INTO sbom_records (image_id, format, uri) VALUES
(1, 'SPDX', 'minio://sboms/acme-corp/c3d3e4f5.../sbom.spdx.json'),
(1, 'CycloneDX', 'minio://sboms/acme-corp/c3d3e4f5.../sbom.cyclonedx.json'),
(2, 'SPDX', 'minio://sboms/startup-x/d4e4f5g6.../sbom.spdx.json');

-- Seed CVEs
INSERT INTO cve_findings (image_id, cve_id, severity, description, fix_available) VALUES
(1, 'CVE-2023-12345', 'High', 'Remote code execution vulnerability in libfoo', true),
(1, 'CVE-2023-67890', 'Medium', 'Denial of service in bar-utils', false),
(2, 'CVE-2023-54321', 'Critical', 'SQL injection in database driver', true);

-- Seed Attestations
INSERT INTO attestations (image_id, type, uri) VALUES
(1, 'provenance', 'minio://attestations/acme-corp/c3d3e4f5.../provenance.json'),
(1, 'vuln-scan', 'minio://attestations/acme-corp/c3d3e4f5.../scan-report.json');

-- Seed Build Events
INSERT INTO build_events (tenant_namespace, image_id, trigger_type, status) VALUES
('acme-corp', 1, 'manual', 'completed'),
('startup-x', 2, 'webhook', 'completed');
