-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Clear existing data and reset sequences
TRUNCATE
    customers,
    images,
    build_events,
    sbom_records,
    cve_findings,
    attestations
RESTART IDENTITY CASCADE;

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
    uri VARCHAR(1024) NOT NULL,
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
    uri VARCHAR(1024) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_attestations_image_id ON attestations(image_id);

-- New Tables for this step
CREATE TABLE IF NOT EXISTS customer_image_usage (
    id SERIAL PRIMARY KEY,
    tenant_namespace VARCHAR(255) REFERENCES customers(namespace) ON DELETE CASCADE,
    image_id INT REFERENCES images(id) ON DELETE CASCADE,
    version_pinned BOOLEAN DEFAULT FALSE,
    runtime_info VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sla_violations (
    id SERIAL PRIMARY KEY,
    tenant_namespace VARCHAR(255) REFERENCES customers(namespace) ON DELETE CASCADE,
    cve_finding_id INT REFERENCES cve_findings(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    resolved_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    tenant_namespace VARCHAR(255) REFERENCES customers(namespace) ON DELETE CASCADE,
    type VARCHAR(100),
    payload JSONB,
    sent_at TIMESTAMPTZ,
    status VARCHAR(50) DEFAULT 'pending'
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id SERIAL PRIMARY KEY,
    tenant_namespace VARCHAR(255),
    action VARCHAR(100),
    actor VARCHAR(255),
    details JSONB,
    timestamp TIMESTAMPTZ DEFAULT NOW()
);


-- Seed data
INSERT INTO customers (namespace, name, contact_info, sla_tier) VALUES
('acme-corp', 'ACME Corporation', 'contact@acme.corp', 'premium'),
('startup-x', 'Startup X Inc.', 'devops@startup-x.com', 'standard');

INSERT INTO images (tenant_namespace, digest, tags, slsa_level) VALUES
('acme-corp', 'sha256:abcdef123456', ARRAY['latest', '1.2.3'], 3),
('startup-x', 'sha256:fedcba654321', ARRAY['latest'], 1);

-- Seed SBOMs
INSERT INTO sbom_records (image_id, format, uri) VALUES
(1, 'SPDX', 'minio://sboms/acme-corp/image1.spdx.json'),
(1, 'CycloneDX', 'minio://sboms/acme-corp/image1.cdx.json'),
(2, 'SPDX', 'minio://sboms/startup-x/image2.spdx.json');

-- Seed CVEs
INSERT INTO cve_findings (image_id, cve_id, severity, description, fix_available) VALUES
(1, 'CVE-2023-0001', 'Critical', 'Remote code execution vulnerability', true),
(1, 'CVE-2023-0002', 'Medium', 'Cross-site scripting', true),
(2, 'CVE-2023-0003', 'High', 'Denial of service in library foo', false);

-- Seed Attestations
INSERT INTO attestations (image_id, type, uri) VALUES
(1, 'provenance', 'minio://attestations/acme-corp/image1.provenance.json'),
(1, 'vuln-scan', 'minio://attestations/acme-corp/image1.vuln.json');

-- Seed Build Events
INSERT INTO build_events (tenant_namespace, image_id, trigger_type, status) VALUES
('acme-corp', 1, 'manual', 'completed'),
('startup-x', 2, 'webhook', 'completed');

-- Seed data for new tables
INSERT INTO customer_image_usage (tenant_namespace, image_id, runtime_info, version_pinned) VALUES
('acme-corp', 1, 'prod-cluster-1', true),
('startup-x', 2, 'staging-cluster', false);

INSERT INTO sla_violations (tenant_namespace, cve_finding_id, status, created_at) VALUES
('acme-corp', 1, 'active', NOW() - INTERVAL '10 days'); -- A pre-existing critical violation for ACME

INSERT INTO audit_logs (tenant_namespace, action, actor, details) VALUES
('acme-corp', 'image_published', 'ci-pipeline-1', '{"image_id": 1, "digest": "sha256:abcdef123456"}'),
('startup-x', 'image_published', 'ci-pipeline-2', '{"image_id": 2, "digest": "sha256:fedcba654321"}');

