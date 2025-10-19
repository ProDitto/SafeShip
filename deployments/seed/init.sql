-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Clear existing data and reset sequences
-- TRUNCATE TABLE customers, images, build_events, sbom_records, cve_findings, attestations, customer_image_usage, sla_violations, notifications, audit_logs RESTART IDENTITY CASCADE;

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


-- Seed Customers
INSERT INTO customers (namespace, name, contact_info, sla_tier) VALUES
('customer-a', 'Customer A Inc.', 'contact@customera.com', 'premium'),
('customer-b', 'Customer B Corp.', 'support@customerb.com', 'standard');

-- Seed Images
INSERT INTO images (tenant_namespace, digest, tags, slsa_level) VALUES
('customer-a', 'sha256:c3d3b3c3d3b3c3d3d3b3c3d3d3b3c3d3d3b3c3d3d3b3c3d3d3b3c3d3d3b3c3d3', '{"1.0.0", "latest"}', 3),
('customer-a', 'sha256:a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1', '{"1.0.1"}', 3),
('customer-a', 'sha256:b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2', '{"2.0.0-clean"}', 4),
('customer-b', 'sha256:d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4', '{"v1-stable"}', 2);

-- Seed SBOM Records
INSERT INTO sbom_records (image_id, format, uri) VALUES
(1, 'SPDX', 'minio://sboms/customer-a/image1-1.0.0.spdx.json'),
(1, 'CycloneDX', 'minio://sboms/customer-a/image1-1.0.0.cyclonedx.json'),
(2, 'SPDX', 'minio://sboms/customer-a/image2-1.0.1.spdx.json'),
(4, 'SPDX', 'minio://sboms/customer-b/image4-v1-stable.spdx.json');

-- Seed CVE Findings
INSERT INTO cve_findings (image_id, cve_id, severity, description, fix_available) VALUES
(1, 'CVE-2023-1111', 'Critical', 'Remote code execution vulnerability in lib-a.', true),
(1, 'CVE-2023-2222', 'High', 'SQL injection in component-b.', false),
(2, 'CVE-2023-3333', 'Medium', 'Cross-site scripting in web framework.', true),
(4, 'CVE-2024-5555', 'High', 'Denial of service in network stack.', true);

-- Seed Attestations
INSERT INTO attestations (image_id, type, uri) VALUES
(1, 'provenance', 'minio://attestations/customer-a/image1-provenance.json'),
(1, 'vuln-scan', 'minio://attestations/customer-a/image1-vuln-scan.json'),
(2, 'provenance', 'minio://attestations/customer-a/image2-provenance.json'),
(3, 'provenance', 'minio://attestations/customer-a/image3-provenance.json'),
(4, 'provenance', 'minio://attestations/customer-b/image4-provenance.json');

-- Seed Build Events
INSERT INTO build_events (tenant_namespace, image_id, trigger_type, status) VALUES
('customer-a', 1, 'webhook', 'completed'),
('customer-a', 2, 'manual', 'completed'),
('customer-a', 3, 'scheduled', 'completed'),
('customer-b', 4, 'api', 'completed'),
('customer-b', NULL, 'webhook', 'pending'); -- A build that hasn't completed yet

-- Seed Customer Image Usage
INSERT INTO customer_image_usage (tenant_namespace, image_id, version_pinned, runtime_info) VALUES
('customer-a', 1, true, 'prod-cluster-1'),
('customer-a', 2, false, 'staging-cluster'),
('customer-b', 4, true, 'prod-main');

-- Seed SLA Violations
-- Assuming CVE-2023-1111 (finding_id=1) for customer-a has breached the 'premium' SLA
INSERT INTO sla_violations (tenant_namespace, cve_finding_id, status) VALUES
('customer-a', 1, 'active');
-- Assuming CVE-2024-5555 (finding_id=4) for customer-b has breached the 'standard' SLA
INSERT INTO sla_violations (tenant_namespace, cve_finding_id, status) VALUES
('customer-b', 4, 'active');


-- Seed Notifications
INSERT INTO notifications (tenant_namespace, type, payload, status) VALUES
('customer-a', 'SLA_VIOLATION', '{"cve_id": "CVE-2023-1111", "severity": "Critical", "image_digest": "sha256:c3d3..."}', 'sent'),
('customer-a', 'BUILD_COMPLETE', '{"image_id": 1, "tags": ["1.0.0", "latest"]}', 'sent');

-- Seed Audit Logs
INSERT INTO audit_logs (tenant_namespace, action, actor, details) VALUES
('customer-a', 'build_triggered', 'webhook-service', '{"trigger": "upstream_commit", "commit_sha": "abcdef123"}'),
('customer-a', 'image_published', 'system', '{"image_id": 1, "digest": "sha256:c3d3..."}'),
('system', 'cve_registered', 'cve-scanner', '{"cve_id": "CVE-2023-1111", "image_id": 1}'),
('customer-b', 'build_triggered', 'api-key-user', '{"trigger": "manual_api_call"}');
