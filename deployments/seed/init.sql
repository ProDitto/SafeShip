-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create customers table to represent tenants
CREATE TABLE customers (
    namespace VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    contact_info VARCHAR(255),
    sla_tier VARCHAR(50) NOT NULL DEFAULT 'standard',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create images table
CREATE TABLE images (
    id SERIAL PRIMARY KEY,
    tenant_namespace VARCHAR(255) NOT NULL REFERENCES customers(namespace),
    digest VARCHAR(255) NOT NULL UNIQUE,
    tags TEXT[],
    slsa_level INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed data for customers
INSERT INTO customers (namespace, name, contact_info, sla_tier) VALUES
('acme-corp', 'ACME Corporation', 'security@acme.corp', 'premium'),
('startup-inc', 'Startup Inc.', 'devops@startup.inc', 'standard');

-- Seed data for images
INSERT INTO images (tenant_namespace, digest, tags, slsa_level) VALUES
('acme-corp', 'sha256:c3d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3', ARRAY['latest', '1.2.3'], 3),
('acme-corp', 'sha256:abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890', ARRAY['1.2.2'], 3),
('startup-inc', 'sha256:fedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321', ARRAY['latest'], 1),
('startup-inc', 'sha256:11223344556677889900aabbccddeeff11223344556677889900aabbccddeeff', ARRAY['v2-beta'], 1);

-- Create indexes for performance
CREATE INDEX idx_images_tenant_namespace ON images(tenant_namespace);

