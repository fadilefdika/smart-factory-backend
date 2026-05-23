-- Connect to db_smartfactory_production before running this script
-- \c db_smartfactory_production

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS devices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL, -- 'sensor', 'actuator', 'gateway'
    status VARCHAR(100) NOT NULL DEFAULT 'inactive', -- 'active', 'inactive', 'offline', 'maintenance'
    metadata JSONB,
    last_seen TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Buat indeks untuk pencarian berdasarkan tipe dan status
CREATE INDEX idx_devices_type ON devices(type);
CREATE INDEX idx_devices_status ON devices(status);
