-- Migration: 004_create_locations_table

CREATE TABLE IF NOT EXISTS locations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lead_id UUID NOT NULL REFERENCES leads(id) ON DELETE CASCADE,
    brand_id UUID NOT NULL REFERENCES brands(id) ON DELETE RESTRICT,
    lat DECIMAL(10,8),
    lng DECIMAL(11,8),
    address TEXT,
    photo VARCHAR(500),
    approval_status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    survey_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_locations_lead_id ON locations(lead_id);
CREATE INDEX idx_locations_brand_id ON locations(brand_id);
CREATE INDEX idx_locations_approval_status ON locations(approval_status);
