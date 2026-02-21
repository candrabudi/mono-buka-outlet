-- Migration: 006_create_payments_table

CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    partnership_id UUID NOT NULL REFERENCES partnerships(id) ON DELETE RESTRICT,
    brand_id UUID NOT NULL REFERENCES brands(id) ON DELETE RESTRICT,
    type VARCHAR(50) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    proof_url VARCHAR(500),
    verified_status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    verified_by UUID REFERENCES users(id) ON DELETE SET NULL,
    verified_at TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_payments_partnership_id ON payments(partnership_id);
CREATE INDEX idx_payments_brand_id ON payments(brand_id);
CREATE INDEX idx_payments_type ON payments(type);
CREATE INDEX idx_payments_verified_status ON payments(verified_status);
