-- Migration: 002_create_brands_table

CREATE TABLE IF NOT EXISTS brands (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    logo VARCHAR(500),
    description TEXT,
    minimum_investment DECIMAL(15,2) NOT NULL DEFAULT 0,
    profit_sharing_percentage DECIMAL(5,2) NOT NULL DEFAULT 0,
    estimated_roi VARCHAR(100),
    location_requirement TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_brands_is_active ON brands(is_active);
CREATE INDEX idx_brands_deleted_at ON brands(deleted_at);
