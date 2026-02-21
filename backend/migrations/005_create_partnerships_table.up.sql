-- Migration: 005_create_partnerships_table

CREATE TABLE IF NOT EXISTS partnerships (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lead_id UUID NOT NULL REFERENCES leads(id) ON DELETE RESTRICT,
    brand_id UUID NOT NULL REFERENCES brands(id) ON DELETE RESTRICT,
    mitra_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    progress_percentage INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    start_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_partnerships_lead_id ON partnerships(lead_id);
CREATE INDEX idx_partnerships_brand_id ON partnerships(brand_id);
CREATE INDEX idx_partnerships_mitra_id ON partnerships(mitra_id);
CREATE INDEX idx_partnerships_status ON partnerships(status);
CREATE INDEX idx_partnerships_deleted_at ON partnerships(deleted_at);
